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
	// Region can be used to select an access point close to the Webhook cluster node.
	// Generally, it can be left unset.
	//
	// If you need to set it, refer to the table in https://next.api.aliyun.com/product/Alidns and
	// fill in the value from the "Region ID" column on that page.
	Region string `json:"region"`

	// AccessKeyIdRef is a credential for accessing Aliyun OpenAPI, which can be created and managed
	// in the RAM console.
	AccessKeyIdRef cmmeta.SecretKeySelector `json:"accessKeyIdRef"`

	// AccessKeySecretRef is the access credential secret that matches AccessKeyIdRef.
	//
	// This field follows Aliyun's naming style; you can configure either this or SecretAccessKeyRef.
	AccessKeySecretRef cmmeta.SecretKeySelector `json:"accessKeySecretRef"`

	// SecretAccessKeyRef is the access credential secret that matches AccessKeyIdRef.
	// This field follows Amazon's naming style; you can configure either this or AccessKeySecretRef.
	SecretAccessKeyRef cmmeta.SecretKeySelector `json:"secretAccessKeyRef"`
}

// Validate checks if the config of the webhook is valid.
func (cfg *Config) Validate() error {
	if len(cfg.AccessKeyIdRef.Name) == 0 {
		return errors.New("accessKeyIdRef may not be empty")
	}

	if len(cfg.AccessKeySecretRef.Name) == 0 {
		cfg.SecretAccessKeyRef.DeepCopyInto(&cfg.AccessKeySecretRef)
	}
	if len(cfg.AccessKeySecretRef.Name) == 0 {
		return errors.New("accessKeySecretRef may not be empty")
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
