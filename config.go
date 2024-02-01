// SPDX-License-Identifier: MIT

package main

import (
	"encoding/json"
	"fmt"

	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	"github.com/pkg/errors"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

// Config is a structure that is used to decode into when
// solving a DNS01 challenge.
//
// This information is provided by cert-manager, and may be a reference to
// additional configuration that's needed to solve the challenge for this
// particular certificate or issuer.
//
// This typically includes references to Secret resources containing DNS
// provider credentials, in cases where a 'multi-tenant' DNS solver is being
// created.
//
// You should not include sensitive information here. If credentials need to
// be used by your provider here, you should reference a Kubernetes Secret
// resource and fetch these credentials using a Kubernetes clientset.
type Config struct {
	Region         string                   `json:"region"` // optional
	AccessKeyIdRef cmmeta.SecretKeySelector `json:"testAccessKeyIdRef"`
	// AccessKeySecretRef will serve as the alias name for SecretAccessKeyRef
	AccessKeySecretRef cmmeta.SecretKeySelector `json:"accessKeySecretRef"`
	SecretAccessKeyRef cmmeta.SecretKeySelector `json:"secretAccessKeyRef"`
}

// Validate checks if the config of the webhook is valid.
func (cfg *Config) Validate() error {
	if len(cfg.AccessKeyIdRef.Name) == 0 {
		return errors.New("testAccessKeyIdRef may not be empty")
	}

	if len(cfg.SecretAccessKeyRef.Name) == 0 {
		cfg.AccessKeySecretRef.DeepCopyInto(&cfg.SecretAccessKeyRef)
	}
	if len(cfg.SecretAccessKeyRef.Name) == 0 {
		return errors.New("AccessKeySecretRef may not be empty")
	}

	return nil
}

// loadConfig decodes JSON configuration into the Config struct.
func loadConfig(cfgJSON *extapi.JSON) (*Config, error) {
	var cfg Config

	// handle the 'base case' where no configuration has been provided
	if cfgJSON == nil {
		return &cfg, nil
	}

	if err := json.Unmarshal(cfgJSON.Raw, &cfg); err != nil {
		return nil, fmt.Errorf("error decoding solver config: %v", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate solver config: %v", err)
	}

	return &cfg, nil
}
