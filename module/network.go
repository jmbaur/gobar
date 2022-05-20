package module

import (
	"fmt"

	"github.com/jmbaur/gobar/i3"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

type Network struct {
	Interface string
}

func (n Network) Run(c chan Update, position int) error {
	updates := make(chan netlink.AddrUpdate)
	done := make(chan struct{})
	defer func() {
		done <- struct{}{}
	}()

	if err := netlink.AddrSubscribe(updates, done); err != nil {
		return fmt.Errorf("failed to run network module: %v", err)
	}

	link, err := netlink.LinkByName(n.Interface)
	if err != nil {
		return fmt.Errorf("failed to get link for %s: %v", n.Interface, err)
	}

	var fullText, ipv4, ipv6 string

	v4addrs, err := netlink.AddrList(link, unix.AF_INET)
	if err != nil {
		return fmt.Errorf("failed to get address list for %s: %s", n.Interface, err)
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
		return fmt.Errorf("failed to get address list for %s: %v", n.Interface, err)
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

	switch true {
	case ipv4 != "" && ipv6 != "":
		fullText = fmt.Sprintf("%s: %s %s", n.Interface, ipv4, ipv6)
	case ipv4 != "" && ipv6 == "":
		fullText = fmt.Sprintf("%s: %s", n.Interface, ipv4)
	case ipv4 == "" && ipv6 != "":
		fullText = fmt.Sprintf("%s: %s", n.Interface, ipv6)
	default:
		fullText = fmt.Sprintf("%s: n/a", n.Interface)
	}
	c <- Update{
		Block: i3.Block{
			FullText: fullText,
		},
		Position: position,
	}

	for u := range updates {
		if u.LinkIndex == link.Attrs().Index {
			c <- Update{
				Block: i3.Block{
					FullText: fmt.Sprintf("%s: %s", n.Interface, u.LinkAddress.IP),
				},
			}
		}
	}

	return nil
}
