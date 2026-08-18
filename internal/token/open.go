package token

func PeekKID(raw []byte) (string, Alg, error) {
	env, err := Split(raw)
	if err != nil {
		return "", 0, err
	}
	h, err := DecodeHeader(env.Header)
	if err != nil {
		return "", 0, err
	}
	return h.KID, h.Alg, nil
}

func Open(raw []byte, secret []byte) (Header, Payload, error) {
	env, err := Split(raw)
	if err != nil {
		return Header{}, Payload{}, err
	}
	h, err := DecodeHeader(env.Header)
	if err != nil {
		return Header{}, Payload{}, err
	}
	if h.Alg.MACSize() != len(env.MAC) {
		return Header{}, Payload{}, errAlg
	}
	msg, err := MACMessage(raw, len(env.MAC))
	if err != nil {
		return Header{}, Payload{}, err
	}
	if !VerifyMAC(h.Alg, secret, msg, env.MAC) {
		return Header{}, Payload{}, errString("token: bad mac")
	}
	p, err := DecodePayload(env.Payload)
	if err != nil {
		return Header{}, Payload{}, err
	}
	return h, p, nil
}
