package module

import (
	"fmt"

	"github.com/jmbaur/gobar/color"
	"github.com/jmbaur/gobar/i3"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

type Network struct {
	Interface string
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
	done := make(chan struct{})
	defer func() {
		done <- struct{}{}
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

	var fullText, ipv4, ipv6 string
	col := color.Normal

	v4addrs, err := netlink.AddrList(link, unix.AF_INET)
	if err != nil {
		n.sendError(c, err, position)
		return
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
		n.sendError(c, err, position)
		return
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
		col = color.Green
	case ipv4 != "" && ipv6 == "":
		fullText = fmt.Sprintf("%s: %s", n.Interface, ipv4)
		col = color.Green
	case ipv4 == "" && ipv6 != "":
		fullText = fmt.Sprintf("%s: %s", n.Interface, ipv6)
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

	for u := range updates {
		if u.LinkIndex == link.Attrs().Index {
			c <- Update{
				Block: i3.Block{
					FullText: fmt.Sprintf("%s: %s", n.Interface, u.LinkAddress.IP),
				},
			}
		}
	}
}
