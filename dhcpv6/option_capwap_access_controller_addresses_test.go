package dhcpv6

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/u-root/uio/uio"
)

func TestCAPWAPAccessControllerAddresses(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want []net.IP
	}{
		{
			buf: []byte{
				0, 52, // CAPWAP Access Controller Addresses
				0, 16, // length
				0x2a, 0x03, 0x28, 0x80, 0xff, 0xfe, 0x00, 0x0c, 0xfa, 0xce, 0xb0, 0x0c, 0x00, 0x00, 0x00, 0x35,
			},
			want: []net.IP{
				net.IP{0x2a, 0x03, 0x28, 0x80, 0xff, 0xfe, 0x00, 0x0c, 0xfa, 0xce, 0xb0, 0x0c, 0x00, 0x00, 0x00, 0x35},
			},
		},
		{
			buf: []byte{
				0, 52, // CAPWAP Access Controller Addresses
				0, 32, // length
				0x2a, 0x03, 0x28, 0x80, 0xff, 0xfe, 0x00, 0x0c, 0xfa, 0xce, 0xb0, 0x0c, 0x00, 0x00, 0x00, 0x35,
				0x2a, 0x03, 0x28, 0x80, 0xff, 0xfe, 0x00, 0x0c, 0xfa, 0xce, 0xb0, 0x0c, 0x00, 0x00, 0x00, 0x35,
			},
			want: []net.IP{
				net.IP{0x2a, 0x03, 0x28, 0x80, 0xff, 0xfe, 0x00, 0x0c, 0xfa, 0xce, 0xb0, 0x0c, 0x00, 0x00, 0x00, 0x35},
				net.IP{0x2a, 0x03, 0x28, 0x80, 0xff, 0xfe, 0x00, 0x0c, 0xfa, 0xce, 0xb0, 0x0c, 0x00, 0x00, 0x00, 0x35},
			},
		},
		{
			buf:  nil,
			want: nil,
		},
		{
			buf: []byte{
				0, 52,
				0, 8,
				0x2a, 0x03, 0x28, 0x80, 0xff, 0xfe, 0x00, 0x0c, // invalid IPv6 address
			},
			want: nil,
			err:  uio.ErrUnreadBytes,
		},
		{
			buf:  []byte{0, 52, 0, 1, 0},
			want: nil,
			err:  uio.ErrUnreadBytes,
		},
		{
			buf:  []byte{0, 52, 0},
			want: nil,
			err:  uio.ErrUnreadBytes,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var mo MessageOptions
			if err := mo.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := mo.CAPWAPAccessControllerAddresses(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CAPWAPAccessControllerAddresses = %v, want %v", got, tt.want)
			}

			if tt.want != nil {
				var m MessageOptions
				m.Add(OptCAPWAPAccessControllerAddresses(tt.want...))
				got := m.ToBytes()
				if diff := cmp.Diff(tt.buf, got); diff != "" {
					t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
				}
			}
		})
	}
}
