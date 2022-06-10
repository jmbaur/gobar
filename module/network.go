package module

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"

	col "github.com/jmbaur/gobar/color"
	"github.com/jmbaur/gobar/i3"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

var (
	ErrInvalidPattern = errors.New("invalid interface pattern string")
	ErrNoMatch        = errors.New("no matching interface")
)

type iface struct {
	link netlink.Link
	ipv4 net.IP
	ipv6 net.IP
}

type Network struct {
	// The name of the network interface
	Interface *string
	Pattern   *string
	patternRe *regexp.Regexp

	ifaces []iface
}

func (n *Network) sendError(c chan Update, err error, position int) {
	c <- Update{
		Block: i3.Block{
			FullText: fmt.Sprintf("network: %s", err),
			Color:    col.Red,
		},
		Position: position,
	}
}

func (n *Network) valid() bool {
	return (n.Pattern != nil && n.Interface == nil) ||
		(n.Pattern == nil && n.Interface != nil)
}

func (n *Network) init() error {
	if n.Pattern != nil {
		var err error
		n.patternRe, err = regexp.Compile(*n.Pattern)
		if err != nil {
			return err
		}

		links, err := netlink.LinkList()
		if err != nil {
			return err
		}
		matchedNone := true
		for _, link := range links {
			if matched := n.patternRe.Match([]byte(link.Attrs().Name)); matched {
				matchedNone = false
				n.ifaces = append(n.ifaces, iface{
					link: link,
				})
			}
		}
		if matchedNone {
			return ErrNoMatch
		}
	} else if n.Interface != nil && n.ifaces == nil {
		link, err := netlink.LinkByName(*n.Interface)
		if err != nil {
			return err
		}
		n.ifaces = append(n.ifaces, iface{
			link: link,
		})
	}

	for i, iface := range n.ifaces {
		v4addrs, err := netlink.AddrList(iface.link, unix.AF_INET)
		if err != nil {
			return err
		}
		for _, addr := range v4addrs {
			if chooseIPv4(addr.IP) {
				n.ifaces[i].ipv4 = addr.IP
			}
		}

		v6addrs, err := netlink.AddrList(iface.link, unix.AF_INET6)
		if err != nil {
			return err
		}
		for _, addr := range v6addrs {
			if chooseIPv6(addr.IP, addr.Flags) {
				n.ifaces[i].ipv6 = addr.IP
			}
		}
	}

	return nil
}

func (n *Network) print(c chan Update, position int) {
	var (
		color            = col.Normal
		fullTextUnjoined = []string{}
	)

	for _, iface := range n.ifaces {
		name := iface.link.Attrs().Name
		if iface.ipv4 != nil || iface.ipv6 != nil {
			color = col.Green
		}

		switch true {
		case iface.ipv4 != nil && iface.ipv6 != nil:
			fullTextUnjoined = append(fullTextUnjoined, fmt.Sprintf("%s: %s %s", name, iface.ipv4, iface.ipv6))
		case iface.ipv4 != nil && iface.ipv6 == nil:
			fullTextUnjoined = append(fullTextUnjoined, fmt.Sprintf("%s: %s", name, iface.ipv4))
		case iface.ipv4 == nil && iface.ipv6 != nil:
			fullTextUnjoined = append(fullTextUnjoined, fmt.Sprintf("%s: %s", name, iface.ipv6))
		default:
			if n.patternRe != nil {
				continue
			} else {
				fullTextUnjoined = append(fullTextUnjoined, fmt.Sprintf("%s: n/a", name))
				color = col.Red
			}
		}
	}

	c <- Update{
		Block: i3.Block{
			FullText: strings.Join(fullTextUnjoined, "; "),
			Color:    color,
		},
		Position: position,
	}
}

func (n *Network) Run(tx chan Update, rx chan i3.ClickEvent, position int) {
	if valid := n.valid(); !valid {
		n.sendError(tx, ErrInvalidPattern, position)
	}

	if err := n.init(); err != nil {
		n.sendError(tx, err, position)
		return
	}
	n.print(tx, position)

	updates := make(chan netlink.AddrUpdate)
	done := make(chan struct{}, 1)
	defer func() {
		close(updates)
		close(done)
	}()

	if err := netlink.AddrSubscribe(updates, done); err != nil {
		n.sendError(tx, err, position)
		return
	}

	for update := range updates {
		for i, iface := range n.ifaces {
			if update.LinkIndex == iface.link.Attrs().Index {
				if update.NewAddr {
					if len(update.LinkAddress.IP) == net.IPv4len && chooseIPv4(update.LinkAddress.IP) {
						n.ifaces[i].ipv4 = update.LinkAddress.IP
					} else if chooseIPv6(update.LinkAddress.IP, update.Flags) {
						n.ifaces[i].ipv6 = update.LinkAddress.IP
					} else {
						continue
					}
				} else {
					if iface.ipv4.Equal(update.LinkAddress.IP) {
						n.ifaces[i].ipv4 = nil
					} else if iface.ipv6.Equal(update.LinkAddress.IP) {
						n.ifaces[i].ipv6 = nil
					} else {
						continue
					}
				}
				n.print(tx, position)
			}
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
