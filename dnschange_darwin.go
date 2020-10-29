package dnschange

import (
	"fmt"
	"net"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fdurand/modifyDNS/scutil"
	"github.com/jackpal/gateway"
)

func Change(dns string) {
	// var OriginalDNSServer string
	// var InterfaceName string
	gatewayIP, _ := gateway.DiscoverGateway()
	spew.Dump(gatewayIP)
	var gatewayInterface string
	Interfaces, _ := net.Interfaces()
	for _, v := range Interfaces {
		eth, _ := net.InterfaceByName(v.Name)
		adresses, _ := eth.Addrs()
		for _, adresse := range adresses {
			_, NetIP, _ := net.ParseCIDR(adresse.String())
			if NetIP.Contains(gatewayIP) {
				gatewayInterface = v.Name
			}
		}
	}
	NetInterface := scutil.New(nil)
	err := NetInterface.GetDNSServers(gatewayInterface)
	if err != nil {
		fmt.Println(err)
	}

	NetInterface.SetDNSServer("127.0.0.69")
	time.Sleep(1 * time.Minute)
	NetInterface.ResetDNSServer()

	// for _, v := range NetInterfaces {
	// 	if gatewayIP.String() == v.DefaultGatewayAddress {
	// 		OriginalDNSServer = v.StaticDNSServers
	// 		NetInterface.SetDNSServer(v.Name, dns)
	// 		InterfaceName = v.Name
	// 	}
	// }

	// time.Sleep(1 * time.Minute)
	// restoreDNS(NetInterface, OriginalDNSServer, InterfaceName)
}

// func restoreDNS(NetInterface netsh.Interface, dns string, iface string) {
// 	if dns == "" {
// 		NetInterface.ResetDNSServer(iface)
// 	} else {
// 		NetInterface.SetDNSServer(iface, dns)
// 	}
// }
