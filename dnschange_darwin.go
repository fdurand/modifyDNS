package dnschange

import "github.com/fdurand/modifyDNS/scutil"

func Change(dns string) {
	// var OriginalDNSServer string
	// var InterfaceName string
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
