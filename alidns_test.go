// SPDX-License-Identifier: MIT

package main

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"
	"testing"

	alidns "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/stretchr/testify/assert"
)

func newTestAliDNS(t *testing.T) *AliDNS {
	accessKeyId := os.Getenv("WEBHOOK_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("WEBHOOK_ACCESS_KEY_SECRET")
	if len(accessKeyId) == 0 || len(accessKeySecret) == 0 {
		t.Skipf("no accessKeyId or accessKeySecret set")
	}

	cli, err := alidns.NewClient(&openapi.Config{
		Endpoint:        tea.String("dns.aliyuncs.com"),
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
	})
	if assert.NoError(t, err) {
		return &AliDNS{cli: cli}
	}

	return nil
}

func randomDNSRecord(t *testing.T) (string, string) {
	domainName := os.Getenv("WEBHOOK_DOMAIN_NAME")
	if len(domainName) == 0 {
		t.Skipf("no domain name set")
	}

	subDomain := make([]byte, 4)
	if _, err := io.ReadFull(rand.Reader, subDomain); assert.NoError(t, err) {
		return hex.EncodeToString(subDomain) + "." + domainName + ".", domainName + "."
	}
	return "", ""
}

func TestAliDNS_AddRecord(t *testing.T) {
	if c := newTestAliDNS(t); c != nil {
		t.Run("not exists", func(t *testing.T) {
			if fqdn, zone := randomDNSRecord(t); len(fqdn) != 0 {
				err := c.AddRecord(fqdn, zone, "foobar")
				if assert.NoError(t, err) {
					_ = c.DeleteRecord(fqdn, zone)
				}
			}
		})

		t.Run("same value", func(t *testing.T) {
			if fqdn, zone := randomDNSRecord(t); len(fqdn) != 0 {
				err := c.AddRecord(fqdn, zone, "foobar")
				if assert.NoError(t, err) {
					// add record again
					err = c.AddRecord(fqdn, zone, "foobar")
					if assert.NoError(t, err) {
						_ = c.DeleteRecord(fqdn, zone)
					}
				}
			}
		})

		t.Run("updated", func(t *testing.T) {
			if fqdn, zone := randomDNSRecord(t); len(fqdn) != 0 {
				err := c.AddRecord(fqdn, zone, "foobar")
				if assert.NoError(t, err) {
					// add record again
					err = c.AddRecord(fqdn, zone, "updated")
					if assert.NoError(t, err) {
						_ = c.DeleteRecord(fqdn, zone)
					}
				}
			}
		})
	}
}

func TestAliDNS_DeleteRecord(t *testing.T) {
	if c := newTestAliDNS(t); c != nil {
		t.Run("not exists", func(t *testing.T) {
			if fqdn, zone := randomDNSRecord(t); len(fqdn) != 0 {
				err := c.DeleteRecord(fqdn, zone)
				assert.NoError(t, err)
			}
		})

		t.Run("exists", func(t *testing.T) {
			if fqdn, zone := randomDNSRecord(t); len(fqdn) != 0 {
				err := c.AddRecord(fqdn, zone, "foobar")
				if assert.NoError(t, err) {
					_ = c.DeleteRecord(fqdn, zone)
				}
			}
		})
	}
}
