package dnschange

import (
	"fmt"
	"time"

	"github.com/fdurand/gonetsh/netsh"
	"github.com/jackpal/gateway"
)

func (d DNSStruct) Change(dns string) {
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
			OriginalDNSServer = v.StaticDNSServers
			NetInterface.SetDNSServer(v.Name, dns)
			InterfaceName = v.Name
		}
	}

	time.Sleep(1 * time.Minute)

	d.NetInterface = *NetInterface

	d.RestoreDNS(NetInterface, OriginalDNSServer, InterfaceName)
}

func (d DNSStruct) RestoreDNS(NetInterface netsh.Interface, dns string, iface string) {
	if dns == "" {
		d.NetInterface.ResetDNSServer(iface)
	} else {
		d.NetInterface.SetDNSServer(iface, dns)
	}
}
