package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRemoveU0000(t *testing.T) {
	input := "\x00hello\x00 \x00world\x00"
	expect := "hello world"
	output := RemoveU0000(input)
	require.Equal(t, output, expect)

	input = "\u0000hello\u0000 \u0000world\u0000"
	output = RemoveU0000(input)
	require.Equal(t, output, expect)
}
