// SPDX-License-Identifier: MIT

package main

import (
	"os"
	"testing"
	"time"

	"github.com/cert-manager/cert-manager/pkg/issuer/acme/dns/util"
	acmetest "github.com/cert-manager/cert-manager/test/acme"
)

func TestRunsSuite(t *testing.T) {
	acmetest.NewFixture(new(AliSolver),
		acmetest.SetResolvedZone(os.Getenv("TEST_ZONE_NAME")),
		acmetest.SetDNSName(util.UnFqdn(os.Getenv("TEST_ZONE_NAME"))),
		acmetest.SetAllowAmbientCredentials(false),
		acmetest.SetManifestPath("testdata/alidns"),
		acmetest.SetDNSServer("223.5.5.5:53"),
		acmetest.SetPropagationLimit(10*time.Minute),
	).RunConformance(t)
}
