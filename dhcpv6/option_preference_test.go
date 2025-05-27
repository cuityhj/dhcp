package dhcpv6

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/u-root/uio/uio"
)

func TestPreference(t *testing.T) {
	for i, tt := range []struct {
		buf  []byte
		err  error
		want uint8
	}{
		{
			buf: []byte{
				0, 7, // Preference
				0, 1, // length
				2,
			},
			want: 2,
		},
		{
			buf: []byte{
				0, 7, // Preference
				0, 1, // length
			},
			err: uio.ErrBufferTooShort,
		},
		{
			buf: []byte{
				0, 7, // Preference
				0, 2, // length
				0, 2,
			},
			err: uio.ErrUnreadBytes,
		},
		{
			buf: []byte{0, 7, 0},
			err: uio.ErrUnreadBytes,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var mo MessageOptions
			if err := mo.FromBytes(tt.buf); !errors.Is(err, tt.err) {
				t.Errorf("FromBytes = %v, want %v", err, tt.err)
			}
			if got := mo.Preference(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Preference= %v, want %v", got, tt.want)
			}

			if tt.err == nil {
				var m MessageOptions
				m.Add(&OptPreference{tt.want})
				got := m.ToBytes()
				if diff := cmp.Diff(tt.buf, got); diff != "" {
					t.Errorf("ToBytes mismatch (-want, +got): %s", diff)
				}
			}
		})
	}
}
