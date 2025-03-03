package util

func IsVisiableUSASCII(b byte) bool {
	/*
		VCHAR   =  %x21-7E  ; visible (printing) characters
		RFC5234 B.1
	*/
	return b >= 0x21 && b <= 0x7E
}

func IsUSASCII(b byte) bool {
	// USASCII ranges from 0 to 127 (dec) or 0 to 7F (hex)
	return b <= 0x7F
}

func IsCTLCharASCII(b byte) bool {
	/*
		CTL =  %x00-1F / %x7F    ; controls
		RFC5234 B.1
	*/
	return b <= 0x1F || b == 0x7F
}
