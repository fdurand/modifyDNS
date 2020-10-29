package dnschange

import (
	"fmt"

	"github.com/fdurand/gonetsh/netsh"
	"github.com/jackpal/gateway"
)

func (d *DNSStruct) Change(dns string) {
	var OriginalDNSServer string
	var InterfaceName string
	gatewayIP, _ := gateway.DiscoverGateway()
	NetInterface := netsh.New(nil)
	NetInterfaces, err := NetInterface.GetInterfaces()
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range NetInterfaces {
		if gatewayIP.String() == v.DefaultGatewayAddress {
			d.NetInterface = NetInterface
			d.SetDNSServer(dns)
		}
	}
}

func (d *DNSStruct) RestoreDNS(dns string) {
	d.NetInterface.(netsh.Interface).ResetDNSServer()
}
