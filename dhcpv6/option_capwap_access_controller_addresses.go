package dhcpv6

import (
	"fmt"
	"net"

	"github.com/u-root/uio/uio"
)

// OptCAPWAPAccessControllerAddresses returns a CAPWAP AC DHCPv6 Option
// This Option as defined in RFC 5417 Section 3.
func OptCAPWAPAccessControllerAddresses(ips ...net.IP) Option {
	return &optCAPWAPAccessControllerAddresses{Addresses: ips}
}

type optCAPWAPAccessControllerAddresses struct {
	Addresses []net.IP
}

// Code returns the Option Code for this option.
func (o *optCAPWAPAccessControllerAddresses) Code() OptionCode {
	return OptionCAPWAPAccessControllerAddresses
}

// ToBytes serializes this option.
func (o *optCAPWAPAccessControllerAddresses) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	for _, ns := range o.Addresses {
		buf.WriteBytes(ns.To16())
	}

	return buf.Data()
}

// String returns a human-readable representation of the option.
func (o *optCAPWAPAccessControllerAddresses) String() string {
	return fmt.Sprintf("%s: %v", o.Code(), o.Addresses)
}

// FromBytes builds an optCAPWAPAccessControllerAddresses structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func (o *optCAPWAPAccessControllerAddresses) FromBytes(data []byte) error {
	buf := uio.NewBigEndianBuffer(data)
	o.Addresses = make([]net.IP, 0, buf.Len()/net.IPv6len)
	for buf.Has(net.IPv6len) {
		o.Addresses = append(o.Addresses, buf.CopyN(net.IPv6len))
	}

	return buf.FinError()
}
