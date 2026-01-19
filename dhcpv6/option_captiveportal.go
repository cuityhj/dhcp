package dhcpv6

import (
	"fmt"
)

// OptCaptivePortal returns a OptionCaptivePortal as defined by RFC 8910. section 2.2.
func OptCaptivePortal(url string) Option {
	return &optCaptivePortal{url}
}

type optCaptivePortal struct {
	url string
}

// Code returns the option code
func (op optCaptivePortal) Code() OptionCode {
	return OptionCaptivePortal
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op optCaptivePortal) ToBytes() []byte {
	return []byte(op.url)
}

func (op optCaptivePortal) String() string {
	return fmt.Sprintf("%s: %s", op.Code(), op.url)
}

// FromBytes builds an optCaptivePortal structure from a sequence
// of bytes. The input data does not include option code and length bytes.
func (op *optCaptivePortal) FromBytes(data []byte) error {
	op.url = string(data)
	return nil
}
