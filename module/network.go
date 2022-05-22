package module

import (
	"fmt"
	"net"

	col "github.com/jmbaur/gobar/color"
	"github.com/jmbaur/gobar/i3"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

type Network struct {
	// The name of the network interface
	Interface string

	link netlink.Link
	ipv4 net.IP
	ipv6 net.IP
}

func (n *Network) sendError(c chan Update, err error, position int) {
	c <- Update{
		Block: i3.Block{
			FullText: fmt.Sprintf("%s: %s", n.Interface, err),
			Color:    col.Red,
		},
		Position: position,
	}
}

func (n *Network) init() error {
	if n.link == nil {
		var err error
		if n.link, err = netlink.LinkByName(n.Interface); err != nil {
			return err
		}
	}

	v4addrs, err := netlink.AddrList(n.link, unix.AF_INET)
	if err != nil {
		return err
	}
	for _, addr := range v4addrs {
		if chooseIPv4(addr.IP) {
			n.ipv4 = addr.IP
		}
	}

	v6addrs, err := netlink.AddrList(n.link, unix.AF_INET6)
	if err != nil {
		return err
	}
	for _, addr := range v6addrs {
		if chooseIPv6(addr.IP, addr.Flags) {
			n.ipv6 = addr.IP
		}
	}

	return nil
}

func (n *Network) print(c chan Update, position int) {
	var (
		color    = col.Normal
		fullText string
	)

	if n.ipv4 != nil || n.ipv6 != nil {
		color = col.Green
	}

	switch true {
	case n.ipv4 != nil && n.ipv6 != nil:
		fullText = fmt.Sprintf("%s: %s %s", n.Interface, n.ipv4, n.ipv6)
	case n.ipv4 != nil && n.ipv6 == nil:
		fullText = fmt.Sprintf("%s: %s", n.Interface, n.ipv4)
	case n.ipv4 == nil && n.ipv6 != nil:
		fullText = fmt.Sprintf("%s: %s", n.Interface, n.ipv6)
	default:
		fullText = fmt.Sprintf("%s: n/a", n.Interface)
		color = col.Red
	}
	c <- Update{
		Block: i3.Block{
			FullText: fullText,
			Color:    color,
		},
		Position: position,
	}
}

func (n *Network) Run(c chan Update, position int) {
	if err := n.init(); err != nil {
		n.sendError(c, err, position)
		return
	}
	n.print(c, position)

	updates := make(chan netlink.AddrUpdate)
	done := make(chan struct{}, 1)
	defer func() {
		close(updates)
		close(done)
	}()

	if err := netlink.AddrSubscribe(updates, done); err != nil {
		n.sendError(c, err, position)
		return
	}

	for update := range updates {
		if update.LinkIndex == n.link.Attrs().Index {
			if update.NewAddr {
				if len(update.LinkAddress.IP) == net.IPv4len && chooseIPv4(update.LinkAddress.IP) {
					n.ipv4 = update.LinkAddress.IP
				} else if chooseIPv6(update.LinkAddress.IP, update.Flags) {
					n.ipv6 = update.LinkAddress.IP
				} else {
					continue
				}
			} else {
				if n.ipv4.Equal(update.LinkAddress.IP) {
					n.ipv4 = nil
				} else if n.ipv6.Equal(update.LinkAddress.IP) {
					n.ipv6 = nil
				} else {
					continue
				}
			}
			n.print(c, position)
		}
	}
}

func chooseIPv4(ip net.IP) bool {
	return ip.IsPrivate()
}

func chooseIPv6(ip net.IP, flags int) bool {
	if flags&unix.IFA_F_MANAGETEMPADDR > 0 {
		return false
	}
	return !ip.IsPrivate() && ip.IsGlobalUnicast()
}
