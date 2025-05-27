package dhcpv6

import (
	"fmt"

	"github.com/cuityhj/dhcp/rfc1035label"
	"github.com/u-root/uio/uio"
)

// The Flags Field as defined in RFC4704 Section 4.1.
//
// The "S" bit indicates whether the server SHOULD or SHOULD NOT perform
// the AAAA RR (FQDN-to-address) DNS updates.  A client sets the bit to
// 0 to indicate that the server SHOULD NOT perform the updates and 1 to
// indicate that the server SHOULD perform the updates.  The state of
// the bit in the reply from the server indicates the action to be taken
// by the server; if it is 1, the server has taken responsibility for
// AAAA RR updates for the FQDN.
//
// The "O" bit indicates whether the server has overridden the client's
// preference for the "S" bit.  A client MUST set this bit to 0.  A
// server MUST set this bit to 1 if the "S" bit in its reply to the
// client does not match the "S" bit received from the client.
//
// The "N" bit indicates whether the server SHOULD NOT perform any DNS
// updates.  A client sets this bit to 0 to request that the server
// SHOULD perform updates (the PTR RR and possibly the AAAA RR based on
// the "S" bit) or to 1 to request that the server SHOULD NOT perform
// any DNS updates.  A server sets the "N" bit to indicate whether the
// server SHALL (0) or SHALL NOT (1) perform DNS updates.  If the "N"
// bit is 1, the "S" bit MUST be 0.
const (
	FLAG_S    = 0x01
	FLAG_O    = 0x02
	FLAG_N    = 0x04
	FLAG_MASK = 0x7
)

// OptFQDN implements OptionFQDN option.
//
// https://tools.ietf.org/html/rfc4704
type OptFQDN struct {
	Flags      uint8
	DomainName *rfc1035label.Labels
}

// Code returns the option code.
func (op *OptFQDN) Code() OptionCode {
	return OptionFQDN
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op *OptFQDN) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write8(op.Flags)
	buf.WriteBytes(op.DomainName.ToBytes())
	return buf.Data()
}

func (op *OptFQDN) String() string {
	return fmt.Sprintf("%s: {Flags=%d DomainName=%s}", op.Code(), op.Flags, op.DomainName)
}

// FromBytes deserializes from bytes to build a OptFQDN structure.
func (op *OptFQDN) FromBytes(data []byte) error {
	var err error
	buf := uio.NewBigEndianBuffer(data)
	op.Flags = buf.Read8()
	op.DomainName, err = rfc1035label.FromBytes(buf.ReadAll())
	if err != nil {
		return err
	}
	return buf.FinError()
}
