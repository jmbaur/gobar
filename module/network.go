package module

import (
	"fmt"
	"time"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

type Network struct {
	Interface string
}

func (n Network) Interval() time.Duration {
	return 10 * time.Second
}

func (n Network) String() string {
	var ipv4, ipv6 string

	link, err := netlink.LinkByName(n.Interface)
	if err != nil {
		return fmt.Sprintf("%s: %s", n.Interface, err)
	}
	v4addrs, err := netlink.AddrList(link, unix.AF_INET)
	if err != nil {
		return fmt.Sprintf("%s: %s", n.Interface, err)
	}
	for _, a := range v4addrs {
		if a.Flags&unix.IFA_F_MANAGETEMPADDR > 0 {
			continue
		}
		if a.IP.IsGlobalUnicast() && a.IP.IsPrivate() {
			ipv4 = a.IPNet.String()
		}
	}
	v6addrs, err := netlink.AddrList(link, unix.AF_INET6)
	if err != nil {
		return fmt.Sprintf("%s: %s", n.Interface, err)
	}
	for _, a := range v6addrs {
		if a.Flags&unix.IFA_F_MANAGETEMPADDR > 0 {
			continue
		}
		if a.Flags&unix.IFA_F_TEMPORARY > 0 || a.Flags&unix.IFA_F_PERMANENT > 0 {
			if a.IP.IsGlobalUnicast() || a.IP.IsPrivate() {
				ipv6 = a.IPNet.String()
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
