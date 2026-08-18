package codec

const (
	Magic0 byte = 'T'
	Magic1 byte = 'G'
	Magic2 byte = 'C'
	Magic3 byte = '1'
	VersionV1 byte = 0x01
)

var Magic = [4]byte{Magic0, Magic1, Magic2, Magic3}

func MagicOK(b []byte) bool {
	return len(b) >= 4 && b[0] == Magic0 && b[1] == Magic1 && b[2] == Magic2 && b[3] == Magic3
}

func WriteMagic(dst []byte) []byte {
	return append(dst, Magic0, Magic1, Magic2, Magic3)
}

func EnvelopeOverhead(macLen int) int {
	// magic(4) + ver(1) + hdrLen(2) + payLen(2) + mac
	return 9 + macLen
}

func MinEnvelopeSize(macLen int) int {
	return EnvelopeOverhead(macLen)
}
