package codec

import "testing"

func TestTLVRoundtrip(t *testing.T) {
	m := NewMap()
	m.AddString(FIss, "iss")
	m.AddString(FSub, "sub")
	m.AddString(FAud, "a")
	m.AddString(FAud, "b")
	b, err := EncodeMap(m, MaxPayloadBytes)
	if err != nil {
		t.Fatal(err)
	}
	got, err := DecodeMap(b, false, true)
	if err != nil {
		t.Fatal(err)
	}
	if got.FirstString(FIss) != "iss" || got.FirstString(FSub) != "sub" {
		t.Fatal(got.FirstString(FIss), got.FirstString(FSub))
	}
	aud := got.AllString(FAud)
	if len(aud) != 2 || aud[0] != "a" || aud[1] != "b" {
		t.Fatalf("%v", aud)
	}
}

func TestFieldOrderRejected(t *testing.T) {
	w := NewWriter(1024)
	w.TLVString(FSub, "x")
	w.TLVString(FIss, "y")
	b := w.Bytes()
	if _, err := DecodeMap(b, false, true); err == nil {
		t.Fatal("expected order error")
	}
}

func TestAttrsCanonical(t *testing.T) {
	b, err := EncodeAttrs(map[string]string{"z": "1", "a": "2"})
	if err != nil {
		t.Fatal(err)
	}
	m, err := DecodeAttrs(b)
	if err != nil {
		t.Fatal(err)
	}
	if m["z"] != "1" || m["a"] != "2" {
		t.Fatalf("%v", m)
	}
}

func TestWriterLimit(t *testing.T) {
	w := NewWriter(4)
	w.Raw([]byte("hello"))
	if w.Err() == nil {
		t.Fatal("expected too large")
	}
}

func TestMagic(t *testing.T) {
	var buf []byte
	buf = WriteMagic(buf)
	if !MagicOK(buf) {
		t.Fatal("magic")
	}
	buf[0] ^= 1
	if MagicOK(buf) {
		t.Fatal("tampered magic")
	}
}
