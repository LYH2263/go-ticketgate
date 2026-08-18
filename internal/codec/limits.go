package codec

const (
	MaxTokenBytes   = 64 * 1024
	MaxHeaderBytes  = 1024
	MaxPayloadBytes = 32 * 1024
	MaxStringBytes  = 4096
	MaxAttrKey      = 64
	MaxAttrVal      = 512
	MaxAttrs        = 32
	MaxAudiences    = 8
	MaxScopes       = 16
	MaxKidBytes     = 64
	MaxJTIBytes     = 64
	MaxNonceBytes   = 64
	MaxBindBytes    = 64
	MaxAudLen       = 128
	MaxScopeLen     = 64
	MaxIssuerLen    = 128
	MaxSubjectLen   = 256
)

func CheckString(s string, max int) bool {
	return len(s) <= max
}

func CheckBytes(b []byte, max int) bool {
	return len(b) <= max
}

func OverflowU16(n int) bool {
	return n < 0 || n > 0xFFFF
}
