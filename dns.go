package main

import (
	"context"
	"slices"

	"github.com/cert-manager/cert-manager/pkg/issuer/acme/dns/util"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

var (
	// alidnsServers lists the public AliDNS resolvers used to discover the authoritative zone.
	alidnsServers = []string{"223.5.5.5:53", "223.6.6.6:53"}
)

// resolveZone resolves the SOA for an FQDN and returns the primary nameserver of the matching zone.
func resolveZone(ctx context.Context, fqdn string, nameservers []string, recursive bool) (string, error) {
	in, err := dnsQuery(ctx, fqdn, nameservers, recursive)
	if err != nil {
		return "", err
	}

	// Search both the answer and authority sections because SOA records may be
	// returned in either place depending on how the resolver responds.
	for _, ans := range slices.Concat(in.Answer, in.Ns) {
		if soa, ok := ans.(*dns.SOA); ok {
			return soa.Hdr.Name, nil
		}
	}
	return "", errors.New("no SOA records found")
}

// dnsQuery queries the given nameservers for the SOA of an FQDN.
//
// This implementation is adapted from
//   - https://github.com/cert-manager/cert-manager/blob/master/pkg/issuer/acme/dns/util/wait.go
func dnsQuery(ctx context.Context, fqdn string, nameservers []string, recursive bool) (*dns.Msg, error) {
	m := new(dns.Msg)
	m.SetQuestion(fqdn, dns.TypeSOA)
	m.SetEdns0(4096, false)
	m.RecursionDesired = recursive

	udp := &dns.Client{Net: "udp", Timeout: util.DNSTimeout}
	tcp := &dns.Client{Net: "tcp", Timeout: util.DNSTimeout}

	// Will retry the request based on the number of servers (n+1)
	for _, ns := range nameservers {
		in, _, err := udp.ExchangeContext(ctx, m, ns)
		if in != nil && !in.Truncated {
			return in, err
		}

		// Try TCP if UDP fails
		in, _, err = tcp.ExchangeContext(ctx, m, ns)
		if in != nil && !in.Truncated {
			return in, err
		}
	}
	return nil, errors.New("failed to resolve FQDN")
}
