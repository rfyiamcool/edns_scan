package main

import (
	"fmt"

	edns "github.com/rfyiamcool/edns_scan"
)

func main() {
	dns_server := "8.8.8.8"
	domain := "xiaorui.cc"
	client_ip := "210.12.138.91"

	ipList, rtt, err := edns.ResolvAddrTypeA(dns_server, domain, client_ip)
	fmt.Println(ipList, rtt, err)
}
