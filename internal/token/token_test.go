package token

import "testing"

func TestIssueOpen(t *testing.T) {
	secret := make([]byte, 32)
	for i := range secret {
		secret[i] = byte(i + 1)
	}
	p := Payload{
		Issuer:  "iss",
		Subject: "sub",
		Audience: []string{"api"},
		JTI:     "j1",
		Nonce:   "n1",
	}
	raw, err := IssueBytes("k-ab", AlgHS256, secret, p)
	if err != nil {
		t.Fatal(err)
	}
	h, got, err := Open(raw, secret)
	if err != nil {
		t.Fatal(err)
	}
	if h.KID != "k-ab" || h.Alg != AlgHS256 {
		t.Fatalf("%+v", h)
	}
	if !PayloadEqual(p, got) {
		t.Fatalf("%+v vs %+v", p, got)
	}
	raw[len(raw)-1] ^= 1
	if _, _, err := Open(raw, secret); err == nil {
		t.Fatal("expected mac fail")
	}
}

func TestPeekKID(t *testing.T) {
	secret := make([]byte, 32)
	raw, err := IssueBytes("k-zz", AlgHS256, secret, Payload{JTI: "j"})
	if err != nil {
		t.Fatal(err)
	}
	kid, alg, err := PeekKID(raw)
	if err != nil || kid != "k-zz" || alg != AlgHS256 {
		t.Fatal(kid, alg, err)
	}
}

func TestBadMagic(t *testing.T) {
	if _, err := Split([]byte("XXXX")); err == nil {
		t.Fatal("expected error")
	}
}
