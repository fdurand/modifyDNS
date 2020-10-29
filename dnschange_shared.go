package dnschange

type DNSStruct struct {
	NetInterface interface{}
}

func NewDNSChange() *DNSStruct {
	d := &DNSStruct{}
	return d
}
