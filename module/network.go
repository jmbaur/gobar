package module

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"sort"

	col "github.com/jmbaur/gobar/color"
	"github.com/jmbaur/gobar/i3"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

var (
	ErrInvalidPattern = errors.New("invalid interface pattern string")
	ErrNoMatch        = errors.New("no matching interface")
)

var allInterfaces = ""

type iface struct {
	hideIP   bool
	link     netlink.Link
	ipv4     net.IP
	ipv4Mask net.IPMask
	ipv6     net.IP
	ipv6Mask net.IPMask
}

type Network struct {
	// The name of the network interface
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
				n.ifaces = append(n.ifaces, iface{link: link})
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
		n.ifaces = append(n.ifaces, iface{link: link})
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
			n.ifaces[i].ipv4 = v4addrs[0].IP
			n.ifaces[i].ipv4Mask = v4addrs[0].Mask
		}

		v6addrs, err := netlink.AddrList(iface.link, unix.AF_INET6)
		if err != nil {
			return err
		}

		if len(v6addrs) > 0 {
			sort.SliceStable(v6addrs, func(i, j int) bool {
				return prioritizeIPv6(v6addrs[i].IP, v6addrs[i].Flags) >= prioritizeIPv6(v6addrs[j].IP, v6addrs[j].Flags)
			})
			n.ifaces[i].ipv6 = v6addrs[0].IP
			n.ifaces[i].ipv6Mask = v6addrs[0].Mask
		}
	}

	return nil
}

func (n *Network) print(c chan i3.Block, ifaceName string, err error) {
	if err != nil {
		c <- i3.Block{
			Name:      "network",
			Instance:  "network",
			FullText:  fmt.Sprintf("network: %s", err),
			ShortText: "network: error",
			MinWidth:  len("network: error"),
			Color:     col.Red,
		}
		return
	}
	if len(n.ifaces) == 0 {
		c <- i3.Block{
			Name:      "network",
			Instance:  "network",
			FullText:  "network: no interfaces",
			ShortText: "network: no interfaces",
			MinWidth:  len("network: no interfaces"),
			Color:     col.Red,
		}
		return
	}

	for _, iface := range n.ifaces {
		// Don't refresh the block if the interface has no new data.
		if ifaceName != allInterfaces && iface.link.Attrs().Name != ifaceName {
			continue
		}

		var (
			color     = col.Normal
			fullText  string
			shortText string
		)

		name := iface.link.Attrs().Name

		switch true {
		case iface.ipv4 != nil && iface.ipv6 != nil:
			color = col.Normal
			v4Size, _ := iface.ipv4Mask.Size()
			v6Size, _ := iface.ipv6Mask.Size()
			shortText = fmt.Sprintf("%s: %s/%d %s/%d", name, iface.ipv4.Mask(iface.ipv4Mask), v4Size, iface.ipv6.Mask(iface.ipv6Mask), v6Size)
			fullText = fmt.Sprintf("%s: %s %s", name, iface.ipv4, iface.ipv6)
		case iface.ipv4 != nil && iface.ipv6 == nil:
			color = col.Yellow
			v4Size, _ := iface.ipv4Mask.Size()
			shortText = fmt.Sprintf("%s: %s/%d", name, iface.ipv4.Mask(iface.ipv4Mask), v4Size)
			fullText = fmt.Sprintf("%s: %s", name, iface.ipv4)
		case iface.ipv4 == nil && iface.ipv6 != nil:
			color = col.Normal
			v6Size, _ := iface.ipv6Mask.Size()
			shortText = fmt.Sprintf("%s: %s/%d", name, iface.ipv6.Mask(iface.ipv6Mask), v6Size)
			fullText = fmt.Sprintf("%s: %s", name, iface.ipv6)
		default:
			if n.patternRe != nil {
				continue
			} else {
				color = col.Red
				fullText = fmt.Sprintf("%s: not connected", name)
				shortText = fullText
			}
		}

		if iface.hideIP {
			fullText = shortText
		}

		c <- i3.Block{
			Name:      "network",
			Instance:  name,
			FullText:  fullText,
			ShortText: shortText,
			MinWidth:  len(shortText),
			Color:     color,
		}
	}
}

func (n *Network) Run(tx chan i3.Block, rx chan i3.ClickEvent) {
	if !n.valid() {
		n.print(tx, "", ErrInvalidPattern)
		return
	}

	if err := n.init(); err != nil {
		n.print(tx, "", err)
		return
	}
	// Print initial info for all configured network interfaces.
	n.print(tx, allInterfaces, nil)

	linkUpdates := make(chan netlink.LinkUpdate)
	addrUpdates := make(chan netlink.AddrUpdate)
	done := make(chan struct{}, 1)
	defer func() {
		close(linkUpdates)
		close(addrUpdates)
		close(done)
	}()

	if err := netlink.LinkSubscribe(linkUpdates, done); err != nil {
		n.print(tx, "", err)
		return
	}

	if err := netlink.AddrSubscribe(addrUpdates, done); err != nil {
		n.print(tx, "", err)
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
						n.print(tx, click.Instance, nil)
					}
				}
			}
		case linkUpdate := <-linkUpdates:
			if n.patternRe == nil {
				continue
			}
			if !n.patternRe.MatchString(linkUpdate.Link.Attrs().Name) {
				continue
			}
			switch linkUpdate.Header.Type {
			case unix.RTM_NEWLINK:
				var existing bool
				for _, iface := range n.ifaces {
					if iface.link.Attrs().Index == linkUpdate.Link.Attrs().Index {
						existing = true
					}
				}
				if !existing {
					n.ifaces = append(n.ifaces, iface{link: linkUpdate.Link})
				}
			case unix.RTM_DELLINK:
				for i, iface := range n.ifaces {
					if iface.link.Attrs().Index == linkUpdate.Link.Attrs().Index {
						if i == len(n.ifaces)-1 {
							n.ifaces = append(n.ifaces[0:i], n.ifaces[i+1:]...)
						} else {
							n.ifaces = n.ifaces[0:i]
						}
						break
					}
				}
			}
		case addrUpdate := <-addrUpdates:
			for i, iface := range n.ifaces {
				if addrUpdate.LinkIndex == iface.link.Attrs().Index {
					if addrUpdate.NewAddr {
						if len(addrUpdate.LinkAddress.IP) == net.IPv4len && prioritizeIPv4(addrUpdate.LinkAddress.IP) >= prioritizeIPv4(n.ifaces[i].ipv4) {
							n.ifaces[i].ipv4 = addrUpdate.LinkAddress.IP
							n.ifaces[i].ipv4Mask = addrUpdate.LinkAddress.Mask
						} else if prioritizeIPv6(addrUpdate.LinkAddress.IP, addrUpdate.Flags) >= prioritizeIPv6(n.ifaces[i].ipv6, 0) {
							n.ifaces[i].ipv6 = addrUpdate.LinkAddress.IP
							n.ifaces[i].ipv6Mask = addrUpdate.LinkAddress.Mask
						} else {
							continue
						}
					} else {
						if iface.ipv4.Equal(addrUpdate.LinkAddress.IP) {
							n.ifaces[i].ipv4 = nil
							n.ifaces[i].ipv4Mask = nil
						} else if iface.ipv6.Equal(addrUpdate.LinkAddress.IP) {
							n.ifaces[i].ipv6 = nil
							n.ifaces[i].ipv6Mask = nil
						} else {
							continue
						}
					}
					n.print(tx, iface.link.Attrs().Name, nil)
				}
			}
		}
	}
}

func prioritizeIPv4(ip net.IP) int {
	score := 0
	switch true {
	case ip.IsGlobalUnicast():
		score += 100
	case ip.IsPrivate():
		score += 90
	}
	return score
}

func prioritizeIPv6(ip net.IP, flags int) int {
	score := 0

	if flags&unix.IFA_F_DEPRECATED > 0 {
		score -= 1000
	}
	if flags&unix.IFA_F_TEMPORARY > 0 {
		score += 300
	}
	if flags&unix.IFA_F_MANAGETEMPADDR > 0 {
		score += 10
	}
	if !ip.IsPrivate() && ip.IsGlobalUnicast() {
		score += 100
	}
	if !ip.IsPrivate() && ip.IsGlobalUnicast() {
		score += 90
	}
	if ip.IsPrivate() {
		score += 10
	}

	return score
}
