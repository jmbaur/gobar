package module

import (
	"fmt"
	"net"

	"github.com/jmbaur/gobar/color"
	"github.com/jmbaur/gobar/i3"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

type Network struct {
	Interface string
	ipv4      net.IP
	ipv6      net.IP
}

func (n Network) sendError(c chan Update, err error, position int) {
	c <- Update{
		Block: i3.Block{
			FullText: fmt.Sprintf("%s: %s", n.Interface, err),
			Color:    color.Red,
		},
		Position: position,
	}
}

func (n Network) Run(c chan Update, position int) {
	updates := make(chan netlink.AddrUpdate)
	done := make(chan struct{}, 1)
	defer func() {
		done <- struct{}{}
		close(updates)
		close(done)
	}()

	if err := netlink.AddrSubscribe(updates, done); err != nil {
		n.sendError(c, err, position)
		return
	}

	link, err := netlink.LinkByName(n.Interface)
	if err != nil {
		n.sendError(c, err, position)
		return
	}

	v4addrs, err := netlink.AddrList(link, unix.AF_INET)
	if err != nil {
		n.sendError(c, err, position)
		return
	}
	for _, addr := range v4addrs {
		if chooseIPv4(addr.IP) {
			n.ipv4 = addr.IP
		}
	}

	v6addrs, err := netlink.AddrList(link, unix.AF_INET6)
	if err != nil {
		n.sendError(c, err, position)
		return
	}
	for _, addr := range v6addrs {
		if chooseIPv6(addr.IP, addr.Flags) {
			n.ipv6 = addr.IP
		}
	}

	hasUpdate := false
	col := color.Green
	for {
		var update netlink.AddrUpdate

		if hasUpdate {
			if update.NewAddr &&
				update.LinkIndex == link.Attrs().Index {
				if len(update.LinkAddress.IP) == net.IPv4len && chooseIPv4(update.LinkAddress.IP) {
					n.ipv4 = update.LinkAddress.IP
				} else if chooseIPv6(update.LinkAddress.IP, update.Flags) {
					n.ipv6 = update.LinkAddress.IP
				}
			} else if !update.NewAddr {
				if n.ipv4.Equal(update.LinkAddress.IP) {
					n.ipv4 = nil
				} else if n.ipv6.Equal(update.LinkAddress.IP) {
					n.ipv6 = nil
				}
			}
		}

		var fullText string
		switch true {
		case !(n.ipv4 == nil) && !(n.ipv6 == nil):
			fullText = fmt.Sprintf("%s: %s %s", n.Interface, n.ipv4, n.ipv6)
			col = color.Green
		case !(n.ipv4 == nil) && n.ipv6 == nil:
			fullText = fmt.Sprintf("%s: %s", n.Interface, n.ipv4)
			col = color.Green
		case n.ipv4 == nil && !(n.ipv6 == nil):
			fullText = fmt.Sprintf("%s: %s", n.Interface, n.ipv6)
			col = color.Green
		default:
			fullText = fmt.Sprintf("%s: n/a", n.Interface)
			col = color.Red
		}
		c <- Update{
			Block: i3.Block{
				FullText: fullText,
				Color:    col,
			},
			Position: position,
		}
		update = <-updates
		hasUpdate = true
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
