package edns

import (
	"errors"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
)

const (
	DefaultPort = ":53"
)

func ResolvAddrTypeA(server string, domain string, clientIp string) ([]string, time.Duration, error) {
	var addr_list []string

	resp, rtt, err := resolve(server, domain, clientIp)
	if err != nil {
		return addr_list, rtt, err
	}

	TypeOk := false
	for _, answer := range resp {
		if answer.Header().Rrtype == dns.TypeA {
			TypeOk = true
			// fmt.Println(answer.Header())
			addr_list = append(addr_list, answer.(*dns.A).A.String())
		}
	}

	if !TypeOk {
		return addr_list, rtt, errors.New("not type A")
	}

	return addr_list, rtt, nil
}

func resolve(server string, domain string, clientIp string) ([]dns.RR, time.Duration, error) {
	// queryType
	var qtype uint16
	qtype = dns.TypeA

	// dnsServer
	if !strings.HasSuffix(server, DefaultPort) {
		server += DefaultPort
	}

	domain = dns.Fqdn(domain)

	msg := new(dns.Msg)
	msg.SetQuestion(domain, qtype)
	msg.RecursionDesired = true

	if clientIp != "" {
		opt := new(dns.OPT)
		opt.Hdr.Name = "."
		opt.Hdr.Rrtype = dns.TypeOPT
		e := new(dns.EDNS0_SUBNET)
		e.Code = dns.EDNS0SUBNET
		e.Family = 1 // ipv4
		e.SourceNetmask = 32
		e.SourceScope = 0
		e.Address = net.ParseIP(clientIp).To4()
		opt.Option = append(opt.Option, e)
		msg.Extra = []dns.RR{opt}
	}

	client := &dns.Client{
		DialTimeout:  5 * time.Second,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
	resp, rtt, err := client.Exchange(msg, server)

	if isEdnsClientSubnet(resp.IsEdns0()) == nil {
		return resp.Answer, rtt, errors.New("is not edns response")
	}

	if err != nil {
		return resp.Answer, rtt, err
	}

	if resp == nil || resp.Rcode != dns.RcodeSuccess {
		return resp.Answer, rtt, errors.New("Test1: no answer")
	}

	return resp.Answer, rtt, err
}

func isEdnsClientSubnet(o *dns.OPT) *dns.EDNS0_SUBNET {
	for _, s := range o.Option {
		switch e := s.(type) {
		case *dns.EDNS0_SUBNET:
			return e
		}
	}
	return nil
}
