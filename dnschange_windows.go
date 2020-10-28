package dnschange

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fdurand/gonetsh/netsh"
	"github.com/jackpal/gateway"
)

func Change(dns string) {
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
	restoreDNS(NetInterface, OriginalDNSServer, InterfaceName)
}

func restoreDNS(NetInterface netsh.Interface, dns string, iface string) {
	if dns == "" {
		NetInterface.ResetDNSServer(iface)
	} else {
		NetInterface.SetDNSServer(iface, dns)
	}
}
