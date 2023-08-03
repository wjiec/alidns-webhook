package main

import (
	"os"
	"testing"
	"time"

	"github.com/cert-manager/cert-manager/pkg/issuer/acme/dns/util"
	"github.com/cert-manager/cert-manager/test/acme/dns"
)

func TestRunsSuite(t *testing.T) {
	dns.NewFixture(new(AliSolver),
		dns.SetResolvedZone(os.Getenv("TEST_ZONE_NAME")),
		dns.SetDNSName(util.UnFqdn(os.Getenv("TEST_ZONE_NAME"))),
		dns.SetAllowAmbientCredentials(false),
		dns.SetManifestPath("testdata/alidns"),
		dns.SetDNSServer("223.5.5.5:53"),
		dns.SetPropagationLimit(10*time.Minute),
	).RunConformance(t)
}
