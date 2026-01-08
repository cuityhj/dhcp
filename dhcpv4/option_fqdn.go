package dhcpv4

import (
	"github.com/u-root/uio/uio"

	"github.com/cuityhj/dhcp/rfc1035label"
)

// The Flags Field is described by RFC 4702, Section 2.1.
//
// The "S" bit indicates whether the server SHOULD or SHOULD NOT perform
// the A RR (FQDN-to-address) DNS updates.  A client sets the bit to 0
// to indicate the server SHOULD NOT perform the updates and 1 to
// indicate the server SHOULD perform the updates. The state of the bit
// in the reply from the server indicates the action to be taken by the
// server; if 1, the server has taken responsibility for A RR updates
// for the FQDN.
//
// The "O" bit indicates whether the server has overridden the client's
// preference for the "S" bit.  A client MUST set this bit to 0.  A
// server MUST set this bit to 1 if the "S" bit in its reply to the
// client does not match the "S" bit received from the client.
//
//The "N" bit indicates whether the server SHOULD NOT perform any DNS
// updates.  A client sets this bit to 0 to request that the server
// SHOULD perform updates (the PTR RR and possibly the A RR based on the
// "S" bit) or to 1 to request that the server SHOULD NOT perform any
// DNS updates.  A server sets the "N" bit to indicate whether the
// server SHALL (0) or SHALL NOT (1) perform DNS updates.  If the "N"
// bit is 1, the "S" bit MUST be 0.
//
// The "E" bit indicates the encoding of the Domain Name field.
// 1 indicates canonical wire format, without compression.
// 0 indicates a now-deprecated ASCII encoding.
const (
	FLAG_S    = 0x01
	FLAG_O    = 0x02
	FLAG_E    = 0x04
	FLAG_N    = 0x08
	FLAG_MASK = 0xF
)

type FQDN struct {
	Flags       uint8
	Rcode1      uint8
	Rcode2      uint8
	DomainName  string
	domainBytes []byte
}

// FromBytes parses a a single byte into FQDN.
func (fqdn *FQDN) FromBytes(data []byte) error {
	buf := uio.NewBigEndianBuffer(data)
	fqdn.Flags = buf.Read8()
	fqdn.Rcode1 = buf.Read8()
	fqdn.Rcode2 = buf.Read8()
	fqdn.domainBytes = buf.ReadAll()
	if fqdn.Flags&FLAG_E == FLAG_E {
		if domainName, err := rfc1035label.FromBytes(fqdn.domainBytes); err != nil {
			return err
		} else {
			fqdn.DomainName = domainName.String()
		}
	} else {
		fqdn.DomainName = string(fqdn.domainBytes)
	}

	return buf.FinError()
}

// ToBytes returns a serialized stream of bytes for this option.
func (fqdn *FQDN) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write8(fqdn.Flags)
	buf.Write8(fqdn.Rcode1)
	buf.Write8(fqdn.Rcode2)
	buf.WriteBytes(fqdn.domainBytes)
	return buf.Data()
}

// String returns a human-readable string for this option.
func (fqdn *FQDN) String() string {
	return fqdn.DomainName
}

// OptFQDN returns a new FQDN option.
//
// The FQDN option is described by RFC 4702.
func OptFQDN(fqdn *FQDN) Option {
	return Option{
		Code:  OptionFQDN,
		Value: fqdn,
	}
}
