package dnschange

import (
	"fmt"
	"net"

	"github.com/fdurand/modifyDNS/scutil"
	"github.com/jackpal/gateway"
)

func (d DNSStruct) Change(dns string) {
	gatewayIP, _ := gateway.DiscoverGateway()
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

	NetInterface.SetDNSServer(dns)

	d.NetInterface = NetInterface
}

func (d DNSStruct) RestoreDNS() {
	d.NetInterface.(scutil.Interface).ResetDNSServer()
}
