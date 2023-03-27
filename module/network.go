package module

import (
	"errors"
	"fmt"
	"net"
	"net/netip"
	"regexp"
	"sort"

	col "github.com/jmbaur/gobar/color"
	"github.com/jmbaur/gobar/i3"
	"github.com/vishvananda/netlink"
	"golang.org/x/exp/slices"
	"golang.org/x/sys/unix"
)

var (
	errNetworkInvalidPattern = errors.New("invalid interface pattern string")
	errNetworkNoMatch        = errors.New("no matching interface")
)

var (
	// 2000::/3
	guaPrefix = netip.PrefixFrom(netip.AddrFrom16([16]byte{0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}), 3)
	// fc00::/7
	ulaPrefix = netip.PrefixFrom(netip.AddrFrom16([16]byte{0xfc, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}), 7)
)

type iface struct {
	hideIP   bool
	link     netlink.Link
	ipv4     net.IP
	ipv4Mask net.IPMask
	ipv6     net.IP
	ipv6Mask net.IPMask
}

// Network provides IP address information for chosen network interfaces. The
// interface can be an exact match on the interface name or a match on a name
// regexp. Only works on Linux.
type Network struct {
	Interface *string
	Pattern   *string

	patternRe *regexp.Regexp
	ifaces    []iface
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
			if matched := n.patternRe.MatchString(link.Attrs().Name); matched {
				matchedNone = false
				n.ifaces = append(n.ifaces, iface{link: link, hideIP: true})
			}
		}
		if matchedNone {
			return errNetworkNoMatch
		}
	} else if n.Interface != nil && n.ifaces == nil {
		link, err := netlink.LinkByName(*n.Interface)
		if err != nil {
			return err
		}
		n.ifaces = append(n.ifaces, iface{link: link, hideIP: true})
	}

	for i, iface := range n.ifaces {
		v4addrs, err := netlink.AddrList(iface.link, unix.AF_INET)
		if err != nil {
			return err
		}
		if len(v4addrs) > 0 {
			sort.SliceStable(v4addrs, func(i, j int) bool {
				return prioritizeIPv4(v4addrs[i].IP) >= prioritizeIPv4(v4addrs[j].IP)
			})
			if prioritizeIPv4(v4addrs[0].IP) > 0 {
				n.ifaces[i].ipv4 = v4addrs[0].IP
				n.ifaces[i].ipv4Mask = v4addrs[0].Mask
			}
		}

		v6addrs, err := netlink.AddrList(iface.link, unix.AF_INET6)
		if err != nil {
			return err
		}

		if len(v6addrs) > 0 {
			sort.SliceStable(v6addrs, func(i, j int) bool {
				return prioritizeIPv6(v6addrs[i].IP, v6addrs[i].Flags) >= prioritizeIPv6(v6addrs[j].IP, v6addrs[j].Flags)
			})
			if prioritizeIPv6(v6addrs[0].IP, v6addrs[0].Flags) > 0 {
				n.ifaces[i].ipv6 = v6addrs[0].IP
				n.ifaces[i].ipv6Mask = v6addrs[0].Mask
			}
		}
	}

	return nil
}

func (n *Network) print(tx chan []i3.Block, err error, c col.Color) {
	if err != nil {
		tx <- []i3.Block{{
			Name:      "network",
			Instance:  "network",
			FullText:  fmt.Sprintf("network: %s", err),
			ShortText: "network: error",
			MinWidth:  len("network: error"),
			Color:     c.Red(),
		}}
		return
	}
	if len(n.ifaces) == 0 {
		tx <- []i3.Block{{
			Name:      "network",
			Instance:  "network",
			FullText:  "network: no interfaces",
			ShortText: "network: no interfaces",
			MinWidth:  len("network: no interfaces"),
			Color:     c.Red(),
		}}
		return
	}

	blocks := []i3.Block{}

	disconnectedInterfaces := 0
	for _, iface := range n.ifaces {
		var printColor string

		name := iface.link.Attrs().Name

		switch true {
		case iface.ipv4 != nil && iface.ipv6 != nil:
			printColor = c.Normal()
		case iface.ipv4 != nil && iface.ipv6 == nil:
			printColor = c.Yellow()
		case iface.ipv4 == nil && iface.ipv6 != nil:
			printColor = c.Normal()
		default:
			disconnectedInterfaces++
			if n.patternRe != nil {
				continue
			}

			printColor = c.Red()
		}

		blocks = append(blocks, i3.Block{
			Name:     "network",
			Instance: name,
			FullText: name,
			MinWidth: len(name),
			Color:    printColor,
		})
	}

	if n.patternRe != nil && len(n.ifaces) == disconnectedInterfaces {
		text := "NET: none"
		blocks = append(blocks, i3.Block{
			Name:     "network",
			Instance: "network",
			FullText: text,
			MinWidth: len(text),
			Color:    c.Red(),
		})
	}

	tx <- blocks
}

// Run implements Module.
func (n *Network) Run(tx chan []i3.Block, rx chan i3.ClickEvent, c col.Color) {
	if !n.valid() {
		n.print(tx, errNetworkInvalidPattern, c)
		return
	}

	if err := n.init(); err != nil {
		n.print(tx, err, c)
		return
	}

	// Print initial info for all configured network interfaces.
	n.print(tx, nil, c)

	linkUpdates := make(chan netlink.LinkUpdate)
	addrUpdates := make(chan netlink.AddrUpdate)
	done := make(chan struct{}, 1)
	defer func() {
		close(linkUpdates)
		close(addrUpdates)
		close(done)
	}()

	if err := netlink.LinkSubscribe(linkUpdates, done); err != nil {
		n.print(tx, err, c)
		return
	}

	if err := netlink.AddrSubscribe(addrUpdates, done); err != nil {
		n.print(tx, err, c)
		return
	}

	for {
		select {
		case click := <-rx:
			switch click.Button {
			case i3.MiddleClick:
				for i, iface := range n.ifaces {
					if iface.link.Attrs().Name == click.Instance {
						n.ifaces[i].hideIP = !n.ifaces[i].hideIP
						n.print(tx, nil, c)
					}
				}
			}
		case linkUpdate := <-linkUpdates:
			if n.patternRe == nil || !n.patternRe.MatchString(linkUpdate.Link.Attrs().Name) {
				continue
			}

			idx := slices.IndexFunc(n.ifaces, func(i iface) bool {
				return int(linkUpdate.Index) == i.link.Attrs().Index
			})
			switch linkUpdate.Header.Type {
			case unix.RTM_NEWLINK:
				if idx < 0 {
					n.ifaces = append(n.ifaces, iface{link: linkUpdate.Link, hideIP: true})
				}
			case unix.RTM_DELLINK:
				if idx < 0 {
					continue
				}

				if idx == len(n.ifaces)-1 {
					n.ifaces = n.ifaces[0:idx]
				} else {
					n.ifaces = append(n.ifaces[0:idx], n.ifaces[idx+1:]...)
				}
			}
		case addrUpdate := <-addrUpdates:
			idx := slices.IndexFunc(n.ifaces, func(i iface) bool {
				return i.link.Attrs().Index == addrUpdate.LinkIndex
			})
			if idx < 0 {
				continue
			}

			iface := n.ifaces[idx]

			if addrUpdate.NewAddr {
				if len(addrUpdate.LinkAddress.IP) == net.IPv4len && prioritizeIPv4(addrUpdate.LinkAddress.IP) >= prioritizeIPv4(n.ifaces[idx].ipv4) {
					n.ifaces[idx].ipv4 = addrUpdate.LinkAddress.IP
					n.ifaces[idx].ipv4Mask = addrUpdate.LinkAddress.Mask
				} else if prioritizeIPv6(addrUpdate.LinkAddress.IP, addrUpdate.Flags) >= prioritizeIPv6(n.ifaces[idx].ipv6, 0) {
					n.ifaces[idx].ipv6 = addrUpdate.LinkAddress.IP
					n.ifaces[idx].ipv6Mask = addrUpdate.LinkAddress.Mask
				} else {
					continue
				}
			} else {
				if iface.ipv4.Equal(addrUpdate.LinkAddress.IP) {
					n.ifaces[idx].ipv4 = nil
					n.ifaces[idx].ipv4Mask = nil
				} else if iface.ipv6.Equal(addrUpdate.LinkAddress.IP) {
					n.ifaces[idx].ipv6 = nil
					n.ifaces[idx].ipv6Mask = nil
				} else {
					continue
				}
			}
			n.print(tx, nil, c)
		}
	}
}

func prioritizeIPv4(ip net.IP) int {
	if ip == nil {
		return -1
	}

	score := 0

	if ip.IsGlobalUnicast() {
		score += 100
	}
	if ip.IsPrivate() {
		score += 90
	}
	if ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
		score -= 1000
	}

	return score
}

func prioritizeIPv6(ip net.IP, flags int) int {
	if ip == nil {
		return -1
	}

	score := 0

	if flags&unix.IFA_F_DEPRECATED > 0 {
		score -= 1000
	}

	if flags&unix.IFA_F_TEMPORARY > 0 {
		score += 300
	}

	if flags&unix.IFA_F_PERMANENT > 0 {
		score += 500
	}

	// often not used for routing, often set as primary
	if flags&unix.IFA_F_MANAGETEMPADDR > 0 {
		score -= 2000
	}

	// may have been obtained via DHCPv6 or is a temp addr
	if flags&unix.IFA_F_NOPREFIXROUTE > 0 {
		score += 500
	}

	if parsed, err := netip.ParseAddr(ip.String()); err == nil {
		if guaPrefix.Contains(parsed) {
			score += 500
		}
		if ulaPrefix.Contains(parsed) {
			score -= 10
		}
	}

	// often used for routing, set from temp addr
	if flags&unix.IFA_F_SECONDARY > 0 {
		score += 10
	}

	if !ip.IsPrivate() && ip.IsGlobalUnicast() {
		score += 100
	}

	if ip.IsLinkLocalUnicast() {
		score -= 1000
	}

	return score
}
