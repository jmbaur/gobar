package module

import (
	"net"
	"testing"

	"golang.org/x/sys/unix"
)

func TestPrioritizeIpv6(t *testing.T) {
	tt := []struct {
		name  string
		ip    net.IP
		flags int
		want  int
	}{
		{
			name:  "localhost",
			ip:    net.IPv6loopback,
			flags: 0,
			want:  0,
		},
		{
			name:  "ula management address",
			ip:    net.IP{0xfc, 0x00, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01},
			flags: unix.IFA_F_NOPREFIXROUTE | unix.IFA_F_MANAGETEMPADDR,
			want:  -1510,
		},
		{
			name:  "ula temporary address",
			ip:    net.IP{0xfc, 0x00, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01},
			flags: unix.IFA_F_SECONDARY | unix.IFA_F_TEMPORARY,
			want:  300,
		},
		{
			name:  "ula dhcpv6 address",
			ip:    net.IP{0xfc, 0x00, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01},
			flags: unix.IFA_F_NOPREFIXROUTE,
			want:  490,
		},
		{
			name:  "gua management address",
			ip:    net.IP{0x20, 0x00, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01},
			flags: unix.IFA_F_NOPREFIXROUTE | unix.IFA_F_MANAGETEMPADDR,
			want:  -900,
		},
		{
			name:  "gua temporary address",
			ip:    net.IP{0x20, 0x00, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01},
			flags: unix.IFA_F_SECONDARY | unix.IFA_F_TEMPORARY,
			want:  910,
		},
	}

	for _, tc := range tt {
		got := prioritizeIPv6(tc.ip, tc.flags)
		if got != tc.want {
			t.Fatalf("%s: got %d, wanted %d\n", tc.name, got, tc.want)
		}
	}
}
