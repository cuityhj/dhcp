package dhcpv6

import "github.com/u-root/uio/uio"

func write40Zero(b *uio.Lexer) {
	b.Write32(0)
	b.Write8(0)
}

func write48Zero(b *uio.Lexer) {
	b.Write32(0)
	b.Write16(0)
}

func write56Zero(b *uio.Lexer) {
	b.Write32(0)
	b.Write16(0)
	b.Write8(0)
}

func write96Zero(b *uio.Lexer) {
	b.Write64(0)
	b.Write32(0)
}
