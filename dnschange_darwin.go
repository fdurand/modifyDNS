package dnschange

import (
	"net"

	"github.com/davecgh/go-spew/spew"
	"github.com/fdurand/modifyDNS/scutil"
	"github.com/jackpal/gateway"
)

func Change(dns string) {
	// var OriginalDNSServer string
	// var InterfaceName string
	gatewayIP, _ := gateway.DiscoverGateway()
	spew.Dump(gatewayIP)
	Interfaces := net.Interfaces()
	for _, v := range Interfaces() {
		spew.Dump(v)
		eth, _ := net.InterfaceByName(v)
		adresses, _ := eth.Addrs()
		for _, adresse := range adresses {
			IP, NetIP, _ = net.ParseCIDR(adresse.String())
			if NetIP.Contains(gatewayIP) {
				spew.Dump(v)
			}
		}
	}
	NetInterface := scutil.New(nil)
	NetInterface.GetDNSServers()
	// if err != nil {
	// 	fmt.Println(err)
	// }
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
