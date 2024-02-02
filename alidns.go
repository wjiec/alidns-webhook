// SPDX-License-Identifier: MIT

package main

import (
	"context"

	alidns "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	acme "github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	"github.com/cert-manager/cert-manager/pkg/issuer/acme/dns/util"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

const (
	DNSRecordType = "TXT"
	ExactSearch   = "EXACT"
)

// AliSolver implements the provider-specific logic needed to
// 'present' an ACME challenge TXT record for your own DNS provider.
type AliSolver struct {
	ctx  context.Context
	kube *kubernetes.Clientset
}

// Name is used as the name for this DNS solver when referencing it on the ACME
// Issuer resource.
//
// This should be unique **within the group name**, i.e. you can have two
// solvers configured with the same Name() **so long as they do not co-exist
// within a single webhook deployment**.
func (s *AliSolver) Name() string {
	return "alidns"
}

// Present is responsible for actually presenting the DNS record with the
// DNS provider.
//
// This method should tolerate being called multiple times with the same value.
// cert-manager itself will later perform a self check to ensure that the
// solver has correctly configured the DNS provider.
func (s *AliSolver) Present(challenge *acme.ChallengeRequest) error {
	klog.Infof("Presenting TXT record: %v", challenge.ResolvedFQDN)
	dns, err := s.loadAliDNS(challenge)
	if err != nil {
		klog.Errorf("Failed to load alidns cause by %q", err)
		return err
	}

	err = dns.AddRecord(challenge.ResolvedFQDN, challenge.ResolvedZone, challenge.Key)
	if err != nil {
		klog.Errorf("Failed to add TXT record for %q cause by %q",
			challenge.ResolvedFQDN, err.Error())
	}

	return err
}

// CleanUp should delete the relevant TXT record from the DNS provider console.
//
// If multiple TXT records exist with the same record name (e.g.
// _acme-challenge.example.com) then **only** the record with the same `key`
// value provided on the ChallengeRequest should be cleaned up.
//
// This is in order to facilitate multiple DNS validations for the same domain
// concurrently.
func (s *AliSolver) CleanUp(challenge *acme.ChallengeRequest) error {
	klog.Infof("cleaning up TXT record: %v", challenge.ResolvedFQDN)
	dns, err := s.loadAliDNS(challenge)
	if err != nil {
		return err
	}

	err = dns.DeleteRecord(challenge.ResolvedFQDN, challenge.ResolvedZone)
	if err != nil {
		klog.Errorf("Failed to delete TXT record for %q cause by %q",
			challenge.ResolvedFQDN, err.Error())
	}

	return err
}

// Initialize will be called when the webhook first starts.
// This method can be used to instantiate the webhook, i.e. initialising
// connections or warming up caches.
// Typically, the kubeClientConfig parameter is used to build a Kubernetes
// client that can be used to fetch resources from the Kubernetes API, e.g.
// Secret resources containing credentials used to authenticate with DNS
// provider accounts.
// The stopCh can be used to handle early termination of the webhook, in cases
// where a SIGTERM or similar signal is sent to the webhook process.
func (s *AliSolver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) (err error) {
	s.kube, err = kubernetes.NewForConfig(kubeClientConfig)

	ctx, cancel := context.WithCancel(context.Background())
	go func() { <-stopCh; cancel() }()
	s.ctx = ctx

	return
}

// loadAliDNS creates an AliDNS and used to solve a challenge with an ACME server.
func (s *AliSolver) loadAliDNS(challenge *acme.ChallengeRequest) (*AliDNS, error) {
	cfg, err := loadConfig(challenge.Config)
	if err != nil {
		return nil, err
	}

	accessKeyId, err := s.loadSecretData(cfg.AccessKeyIdRef, challenge.ResourceNamespace)
	if err != nil {
		return nil, err
	}

	accessKeySecret, err := s.loadSecretData(cfg.SecretAccessKeyRef, challenge.ResourceNamespace)
	if err != nil {
		return nil, err
	}

	var endpoint = "dns.aliyuncs.com"
	if len(cfg.Region) != 0 {
		endpoint = "alidns." + cfg.Region + ".aliyuncs.com"
	}

	cli, err := alidns.NewClient(&openapi.Config{
		Endpoint:        &endpoint,
		AccessKeyId:     tea.String(string(accessKeyId)),
		AccessKeySecret: tea.String(string(accessKeySecret)),
	})
	if err != nil {
		return nil, err
	}

	return &AliDNS{cli: cli}, nil
}

// loadSecretData loads the specified secret from kubernetes resources.
func (s *AliSolver) loadSecretData(selector cmmeta.SecretKeySelector, ns string) ([]byte, error) {
	secret, err := s.kube.CoreV1().Secrets(ns).Get(s.ctx, selector.Name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "failed reading secret %q", ns+"/"+selector.Name)
	}

	if data, ok := secret.Data[selector.Key]; ok {
		if !s.validSecretData(data) {
			return nil, errors.Wrapf(err, "invalid value for secret %q", ns+"/"+selector.Name)
		}
		return data, nil
	}
	return nil, errors.Errorf("couldn't find key %q in secret %q", selector.Key, ns+"/"+selector.Name)
}

// validSecretData reports whether data contains a control byte.
func (s *AliSolver) validSecretData(data []byte) bool {
	for _, b := range data {
		if b <= ' ' || b == 0x7f || b == '\t' {
			return false
		}
	}
	return true
}

// AliDNS is a client for manipulating Aliyun-DNS
// records through openapi
type AliDNS struct {
	cli *alidns.Client
}

// AddRecord adds the specified dns record via openapi
//
// If the dns record already exists, an attempt is made to update this record.
func (dns *AliDNS) AddRecord(fqdn, zone, value string) error {
	queryReq := new(alidns.DescribeDomainRecordsRequest)
	queryReq.SetDomainName(util.UnFqdn(zone))
	queryReq.SetTypeKeyWord(DNSRecordType)
	queryReq.SetKeyWord(fqdn[:len(fqdn)-len(zone)-1])
	queryReq.SetSearchMode(ExactSearch)
	queryResp, err := dns.cli.DescribeDomainRecords(queryReq)
	if err != nil {
		return err
	}

	// add record when not exists
	if *queryResp.Body.TotalCount == 0 {
		req := new(alidns.AddDomainRecordRequest)
		req.SetType(DNSRecordType)
		req.SetDomainName(util.UnFqdn(zone))
		req.SetRR(fqdn[:len(fqdn)-len(zone)-1])
		req.SetValue(value)

		_, err = dns.cli.AddDomainRecord(req)
		return err
	}

	record := *queryResp.Body.DomainRecords.Record[0] // it's okay
	if *record.Value == value {
		return nil
	}

	// update record when already exists
	req := new(alidns.UpdateDomainRecordRequest)
	req.SetRecordId(*record.RecordId)
	req.SetType(DNSRecordType)
	req.SetRR(fqdn[:len(fqdn)-len(zone)-1])
	req.SetValue(value)

	_, err = dns.cli.UpdateDomainRecord(req)
	return err
}

// DeleteRecord deletes the specified dns record via openapi
//
// No error occurs when the dns record does not exist.
func (dns *AliDNS) DeleteRecord(fqdn, zone string) error {
	req := new(alidns.DeleteSubDomainRecordsRequest)
	req.SetDomainName(util.UnFqdn(zone))
	req.SetRR(fqdn[:len(fqdn)-len(zone)-1])
	req.SetType(DNSRecordType)

	_, err := dns.cli.DeleteSubDomainRecords(req)
	return err
}
