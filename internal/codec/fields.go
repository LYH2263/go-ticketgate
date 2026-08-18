package codec

// FieldID 是 TGCT 载荷/头的标签。编码必须按升序写出，解码按升序校验。
type FieldID byte

const (
	// Header
	FKid   FieldID = 0x01
	FAlg   FieldID = 0x02
	FFlags FieldID = 0x03

	// Payload — 顺序即兼容性契约，改动会破坏 MAC 输入。
	FIss   FieldID = 0x10
	FSub   FieldID = 0x11
	FAud   FieldID = 0x12
	FIat   FieldID = 0x13
	FNbf   FieldID = 0x14
	FExp   FieldID = 0x15
	FJTI   FieldID = 0x16
	FNonce FieldID = 0x17
	FScope FieldID = 0x18
	FAttrs FieldID = 0x19
	FBind  FieldID = 0x1A
)

type WireType byte

const (
	TBytes WireType = 1
	TU8    WireType = 2
	TU16   WireType = 3
	TU64   WireType = 4
	TNest  WireType = 5
)

type FieldSpec struct {
	ID         FieldID
	Type       WireType
	Repeatable bool
	Header     bool
}

var specs = map[FieldID]FieldSpec{
	FKid:   {FKid, TBytes, false, true},
	FAlg:   {FAlg, TU8, false, true},
	FFlags: {FFlags, TU16, false, true},
	FIss:   {FIss, TBytes, false, false},
	FSub:   {FSub, TBytes, false, false},
	FAud:   {FAud, TBytes, true, false},
	FIat:   {FIat, TU64, false, false},
	FNbf:   {FNbf, TU64, false, false},
	FExp:   {FExp, TU64, false, false},
	FNonce: {FNonce, TBytes, false, false},
	FJTI:   {FJTI, TBytes, false, false},
	FScope: {FScope, TBytes, true, false},
	FAttrs: {FAttrs, TNest, false, false},
	FBind:  {FBind, TBytes, false, false},
}

func Spec(id FieldID) (FieldSpec, bool) {
	s, ok := specs[id]
	return s, ok
}

func IsRepeatable(id FieldID) bool {
	s, ok := specs[id]
	return ok && s.Repeatable
}

func IsHeader(id FieldID) bool {
	s, ok := specs[id]
	return ok && s.Header
}

func Name(id FieldID) string {
	switch id {
	case FKid:
		return "kid"
	case FAlg:
		return "alg"
	case FFlags:
		return "flags"
	case FIss:
		return "iss"
	case FSub:
		return "sub"
	case FAud:
		return "aud"
	case FIat:
		return "iat"
	case FNbf:
		return "nbf"
	case FExp:
		return "exp"
	case FJTI:
		return "jti"
	case FNonce:
		return "nonce"
	case FScope:
		return "scope"
	case FAttrs:
		return "attrs"
	case FBind:
		return "bind"
	default:
		return "unknown"
	}
}

func HeaderOrder() []FieldID {
	return []FieldID{FKid, FAlg, FFlags}
}

func PayloadOrder() []FieldID {
	return []FieldID{FIss, FSub, FAud, FIat, FNbf, FExp, FJTI, FNonce, FScope, FAttrs, FBind}
}
