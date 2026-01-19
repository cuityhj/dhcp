package dhcpv6

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"github.com/u-root/uio/uio"
)

func TestCaptivePortalParseAndGetter(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want string
	}{
		{
			buf: []byte{
				0, 103, // Captive Portal URL
				0, 17, // length
				'h', 't', 't', 'p', ':', '/', '/', 'p', 'o', 'r', 't', 'a', 'l', '.', 'o', 'r', 'g',
			},
			want: "http://portal.org",
		},
		{
			buf: nil,
		},
		{
			buf: []byte{0, 59, 0},
			err: uio.ErrUnreadBytes,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var mo MessageOptions
			if err := mo.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := mo.CaptivePortal(); got != tt.want {
				t.Errorf("CaptivePortal = %v, want %v", got, tt.want)
			}

			if tt.want != "" {
				var m MessageOptions
				m.Add(OptCaptivePortal(tt.want))
				got := m.ToBytes()
				if diff := cmp.Diff(tt.buf, got); diff != "" {
					t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
				}
			}
		})
	}
}

func TestOptCaptivePortal(t *testing.T) {
	opt := OptCaptivePortal("https://insomniac.slackware.it")
	require.Contains(t, opt.String(), "https://insomniac.slackware.it", "String() should contain the correct BootFileUrl output")
}
