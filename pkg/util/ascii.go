package util

func IsVisiableUSASCII(b byte) bool {
	return b >= 0x21 && b <= 0x7E
}

func IsUSASCII(b byte) bool {
	return b <= 0x7F
}

func IsCTLCharASCII(b byte) bool {
	return b <= 0x1F || b == 0x7F
}
