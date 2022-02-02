package module

import (
	"fmt"
	"log"
	"net"
	"time"
)

type Network struct {
	Interface string
}

func (n Network) Interval() time.Duration {
	return 10 * time.Second
}

func (n Network) String() string {
	defer log.Println("Updated network module")
	iface, err := net.InterfaceByName(n.Interface)
	if err != nil {
		log.Println(err)
		return fmt.Sprintf("%s: %s", n.Interface, err)
	}
	addrs, err := iface.Addrs()
	if err != nil {
		log.Println(err)
		return fmt.Sprintf("%s: %s", n.Interface, err)
	}
	ipv4 := ""
	ipv6 := ""
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		log.Println(addr, ip)
		if ip.IsPrivate() || ip.IsGlobalUnicast() {
			switch len(ip) {
			case net.IPv4len:
				ipv4 = ip.String()
			case net.IPv6len:
				ipv6 = ip.String()
			}
		}
	}
	if ipv4 != "" && ipv6 != "" {
		return fmt.Sprintf("%s: %s %s", n.Interface, ipv4, ipv6)
	}
	if ipv4 != "" && ipv6 == "" {
		return fmt.Sprintf("%s: %s", n.Interface, ipv4)
	}
	if ipv4 == "" && ipv6 != "" {
		return fmt.Sprintf("%s: %s", n.Interface, ipv6)
	}
	return fmt.Sprintf("%s: n/a", n.Interface)
}
