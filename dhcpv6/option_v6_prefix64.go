package dhcpv6

import (
	"fmt"
	"net"

	"github.com/u-root/uio/uio"
)

const (
	MULTICAST_PREFIX64_LENGTH = 12
	MULTICAST_PREFIX64_BITS   = 96
)

// unicast-length:  the prefix length for the IPv6 unicast prefix to be
//      used to synthesize the IPv4-embedded IPv6 addresses of the
//      multicast sources, as an 8-bit unsigned integer.  As specified in
//      [RFC6052], the unicast-length MUST be one of 32, 40, 48, 56, 64,
//      or 96.  This field represents the number of valid leading bits in
//      the prefix.
var ValidUnicastLength = map[uint8]bool{
	32: true,
	40: true,
	48: true,
	56: true,
	64: true,
	96: true,
}

// OptV6Prefix64 returns a V6_Prefix64 DHCPv6 Option with unicast-length
//		and uPrefix64, This Option as defined in RFC 8115 Section 3.
//
// asm-length:  the prefix length for the ASM IPv4-embedded prefix, as
//      an 8-bit unsigned integer.  This field represents the number of
//      valid leading bits in the prefix.  This field MUST be set to 96.
//
// ssm-length:  the prefix length for the SSM IPv4-embedded prefix, as
//      an 8-bit unsigned integer.  This field represents the number of
//      valid leading bits in the prefix.  This field MUST be set to 96.
func OptV6UnicastPrefix64(unicastLength uint8, uprefix64 *net.IPNet) Option {
	return &OptV6Prefix64{
		AsmLength:     MULTICAST_PREFIX64_BITS,
		SsmLength:     MULTICAST_PREFIX64_BITS,
		UnicastLength: unicastLength,
		UPrefix64:     uprefix64,
	}
}

type OptV6Prefix64 struct {
	AsmLength     uint8
	AsmMPrefix64  *net.IPNet
	SsmLength     uint8
	SsmMPrefix64  *net.IPNet
	UnicastLength uint8
	UPrefix64     *net.IPNet
}

// Code returns the Option Code for this option.
func (o *OptV6Prefix64) Code() OptionCode {
	return OptionV6Prefix64
}

// ToBytes serializes this option.
func (o *OptV6Prefix64) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write8(o.AsmLength)
	if o.AsmMPrefix64 != nil {
		buf.WriteBytes(o.AsmMPrefix64.IP.To16()[:MULTICAST_PREFIX64_LENGTH])
	} else {
		write96Zero(buf)
	}

	buf.Write8(o.SsmLength)
	if o.SsmMPrefix64 != nil {
		buf.WriteBytes(o.SsmMPrefix64.IP.To16()[:MULTICAST_PREFIX64_LENGTH])
	} else {
		write96Zero(buf)
	}

	buf.Write8(o.UnicastLength)
	if o.UPrefix64 != nil {
		buf.WriteBytes(o.UPrefix64.IP.To16()[:o.UnicastLength/8])
	} else {
		writeUnicastZeroPrefix64(buf, o.UnicastLength)
	}

	return buf.Data()
}

func writeUnicastZeroPrefix64(buf *uio.Lexer, unicastLength uint8) {
	switch unicastLength {
	case 32:
		buf.Write32(0)
	case 40:
		write40Zero(buf)
	case 48:
		write48Zero(buf)
	case 56:
		write56Zero(buf)
	case 64:
		buf.Write64(0)
	case 96:
		write96Zero(buf)
	}
}

// String returns a human-readable representation of the option.
func (o *OptV6Prefix64) String() string {
	if o.UPrefix64 != nil {
		return fmt.Sprintf("V6 PREFIX64 UPREFIX64: %s", o.UPrefix64.String())
	} else if o.SsmMPrefix64 != nil {
		return fmt.Sprintf("V6 PREFIX64 SSM MPREFIX64: %s", o.SsmMPrefix64.String())
	} else if o.AsmMPrefix64 != nil {
		return fmt.Sprintf("V6 PREFIX64 ASM MPREFIX64: %s", o.AsmMPrefix64.String())
	} else {
		return "V6 PREFIX64"
	}
}

// FromBytes builds an OptV6Prefix64 structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func (o *OptV6Prefix64) FromBytes(data []byte) error {
	buf := uio.NewBigEndianBuffer(data)
	o.AsmLength = buf.Read8()
	if o.AsmLength != MULTICAST_PREFIX64_BITS {
		return fmt.Errorf("invalid asm length %d, expected %d", o.AsmLength, MULTICAST_PREFIX64_BITS)
	}

	o.AsmMPrefix64 = ipnetFromIpBytes(buf.CopyN(MULTICAST_PREFIX64_LENGTH), MULTICAST_PREFIX64_BITS)

	o.SsmLength = buf.Read8()
	if o.SsmLength != MULTICAST_PREFIX64_BITS {
		return fmt.Errorf("invalid ssm length %d, expected %d", o.SsmLength, MULTICAST_PREFIX64_BITS)
	}

	o.SsmMPrefix64 = ipnetFromIpBytes(buf.CopyN(MULTICAST_PREFIX64_LENGTH), MULTICAST_PREFIX64_BITS)

	o.UnicastLength = buf.Read8()
	if valid := ValidUnicastLength[o.UnicastLength]; !valid {
		return fmt.Errorf("invalid unicast length %d, not in [32, 40, 48, 56, 64, 96]", o.UnicastLength)
	}

	o.UPrefix64 = ipnetFromIpBytes(buf.CopyN(int(o.UnicastLength/8)), int(o.UnicastLength))

	return buf.FinError()
}

func ipnetFromIpBytes(ipBytes []byte, ones int) *net.IPNet {
	if isAllZero(ipBytes) {
		return nil
	}

	ip := make(net.IP, 16)
	copy(ip, ipBytes)
	return &net.IPNet{
		IP:   ip,
		Mask: net.CIDRMask(ones, 128),
	}
}

func isAllZero(data []byte) bool {
	for _, d := range data {
		if d != 0 {
			return false
		}
	}

	return true
}
