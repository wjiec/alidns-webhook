// SPDX-License-Identifier: MIT
//go:build ci

package main

import (
	"encoding/json"
	"testing"

	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	"github.com/stretchr/testify/assert"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

var (
	testAccessKeyIdRef = cmmeta.SecretKeySelector{
		LocalObjectReference: cmmeta.LocalObjectReference{
			Name: "alidns-secret",
		},
		Key: "access-key-id",
	}
	testSecretAccessKeyRef = cmmeta.SecretKeySelector{
		LocalObjectReference: cmmeta.LocalObjectReference{
			Name: "alidns-secret",
		},
		Key: "secret-access-key",
	}
)

func TestConfig_Validate(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		correct := &Config{
			AccessKeyIdRef:     testAccessKeyIdRef,
			SecretAccessKeyRef: testSecretAccessKeyRef,
		}

		loaded, err := loadConfig(&extapi.JSON{Raw: mustMarshal(correct)})
		if assert.NoError(t, err) {
			assert.Equal(t, correct, loaded)
		}
	})

	t.Run("compatible", func(t *testing.T) {
		correct := &Config{
			AccessKeyIdRef:     testAccessKeyIdRef,
			AccessKeySecretRef: testSecretAccessKeyRef,
		}

		loaded, err := loadConfig(&extapi.JSON{Raw: mustMarshal(correct)})
		if assert.NoError(t, err) {
			assert.NotEmpty(t, loaded)
		}
	})

	t.Run("empty json", func(t *testing.T) {
		_, err := loadConfig(&extapi.JSON{Raw: []byte("{}")})
		assert.Error(t, err)
	})

	t.Run("no accessKeyId", func(t *testing.T) {
		bad := &Config{
			SecretAccessKeyRef: testSecretAccessKeyRef,
		}

		_, err := loadConfig(&extapi.JSON{Raw: mustMarshal(bad)})
		assert.Error(t, err)
	})

	t.Run("no secretAccessKey", func(t *testing.T) {
		bad := &Config{
			AccessKeyIdRef: testAccessKeyIdRef,
		}

		_, err := loadConfig(&extapi.JSON{Raw: mustMarshal(bad)})
		assert.Error(t, err)
	})
}

func mustMarshal(v any) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
