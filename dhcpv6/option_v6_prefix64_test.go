package dhcpv6

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/u-root/uio/uio"
)

func TestV6Prefix64(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want *OptV6Prefix64
	}{
		{
			buf: []byte{
				0, 113, // V6Prefix64 option code
				0, 31, // length
				96,                                 // asm_length
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // asm_mprefix64
				96,                                 // ssm_length
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // ssm_mprefix64
				32,                     // unicast_length
				0x20, 0x01, 0x0d, 0xb8, // uprefix64
			},
			want: &OptV6Prefix64{
				AsmLength:     96,
				SsmLength:     96,
				UnicastLength: 32,
				UPrefix64: &net.IPNet{
					IP:   net.IP{0x20, 0x01, 0x0d, 0xb8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					Mask: net.CIDRMask(32, 128),
				},
			},
		},
		{
			buf:  nil,
			want: nil,
		},
		{
			buf:  []byte{0, 113, 0},
			want: nil,
			err:  uio.ErrUnreadBytes,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var mo MessageOptions
			if err := mo.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}

			prefix64 := mo.V6Prefix64()
			if !cmp.Equal(prefix64, tt.want) {
				t.Errorf("Prefixes = %#v, want %#v", prefix64, tt.want)
			}

			if tt.want != nil {
				var b MessageOptions
				b.Add(tt.want)
				got := b.ToBytes()
				if diff := cmp.Diff(tt.buf, got); diff != "" {
					t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
				}
			}
		})
	}
}
