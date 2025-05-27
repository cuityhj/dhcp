package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFQDN(t *testing.T) {
	fqdnbytes := []byte{
		0,                                                          //flags E is 0
		0,                                                          //RCODE1
		0,                                                          //RCODE2
		68, 69, 83, 75, 84, 79, 80, 45, 49, 79, 66, 53, 57, 67, 56, //DESKTOP-1OB59C8
	}
	fqdn := &FQDN{}
	err := fqdn.FromBytes(fqdnbytes)
	require.NoError(t, err)
	require.Equal(t, fqdn.DomainName, "DESKTOP-1OB59C8")

	fqdnbytes = []byte{
		4,                               // flags E is 1
		0,                               // RCODE1
		0,                               // RCODE2
		6, 109, 121, 104, 111, 115, 116, // myhost
		7, 101, 120, 97, 109, 112, 108, 101, // example
		4, 99, 111, 109, 46, // com.
		0,
	}
	fqdn = &FQDN{}
	err = fqdn.FromBytes(fqdnbytes)
	require.NoError(t, err)
	require.Equal(t, fqdn.DomainName, "myhost.example.com.")
}
