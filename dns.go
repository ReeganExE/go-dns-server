package main

import (
	"fmt"

	"github.com/miekg/dns"
)

var noop = struct{}{}

type DnsServer struct {
	defaultIp string
	server    *dns.Server
	domains   map[string]struct{}
}

func (d *DnsServer) parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		if q.Qtype == dns.TypeA {
			fmt.Println("Query for", q.Name)
			if _, ok := d.domains[q.Name]; ok {
				if rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, d.defaultIp)); err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		}
	}
}

func (d *DnsServer) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		d.parseQuery(m)
	}

	w.WriteMsg(m)
}

func (d *DnsServer) ListenAndServe() error {
	return d.server.ListenAndServe()
}

func (d *DnsServer) Shutdown() error {
	return d.server.Shutdown()
}

func NewDNSServer(dnsPort int, defaultIP string, domains []string) *DnsServer {
	domainsMap := make(map[string]struct{}, len(domains))
	for _, s := range domains {
		domainsMap[s+"."] = noop
	}

	dnsServer := &DnsServer{
		defaultIp: defaultIP,
		domains:   domainsMap,
	}

	dnsServer.server = &dns.Server{Addr: fmt.Sprintf(":%d", dnsPort), Net: "udp", Handler: dnsServer}

	return dnsServer
}
