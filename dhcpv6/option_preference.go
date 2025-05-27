package dhcpv6

import (
	"fmt"

	"github.com/u-root/uio/uio"
)

func WithPreference(prefValue uint8) Modifier {
	return func(d DHCPv6) {
		d.AddOption(&OptPreference{PrefValue: prefValue})
	}
}

// Preference option as defined in RFC 8415 Section 21.8.
type OptPreference struct {
	PrefValue uint8
}

// Code returns the Option Code for this option.
func (o *OptPreference) Code() OptionCode {
	return OptionPreference
}

// ToBytes serializes this option.
func (o *OptPreference) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write8(o.PrefValue)
	return buf.Data()
}

// String returns a human-readable representation of the option.
func (o *OptPreference) String() string {
	return fmt.Sprintf("%s: %d", o.Code(), o.PrefValue)
}

// FromBytes builds an OptPreference structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func (o *OptPreference) FromBytes(data []byte) error {
	buf := uio.NewBigEndianBuffer(data)
	o.PrefValue = buf.Read8()
	return buf.FinError()
}
