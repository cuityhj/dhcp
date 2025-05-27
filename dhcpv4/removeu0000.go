package dhcpv4

// RemoveU0000 remove \x00 and \u0000
func RemoveU0000(s string) string {
	runes := make([]rune, 0, len(s))
	for _, r := range s {
		if r == 0 {
			continue
		} else {
			runes = append(runes, r)
		}
	}

	return string(runes)
}
