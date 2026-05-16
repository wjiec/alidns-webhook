package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_resolveZone(t *testing.T) {
	domainName := os.Getenv("WEBHOOK_DOMAIN_NAME")
	if len(domainName) == 0 {
		t.Skipf("no domain name set")
	}

	zone, err := resolveZone(t.Context(), "_acme-challenge.foobar."+domainName+".", alidnsServers, true)
	if assert.NoError(t, err) {
		assert.Equal(t, domainName+".", zone)
	}
}
