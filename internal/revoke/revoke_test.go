package revoke

import (
	"testing"
	"time"

	"github.com/LYH2263/go-ticketgate/internal/clock"
)

func TestExactRevoke(t *testing.T) {
	clk := clock.NewFake(time.Unix(1, 0))
	r := New(clk, ModeExact)
	if r.IsRevoked("j") {
		t.Fatal("fresh")
	}
	if err := r.Revoke("j", time.Time{}); err != nil {
		t.Fatal(err)
	}
	if !r.IsRevoked("j") {
		t.Fatal("revoked")
	}
}

func TestBloomAdd(t *testing.T) {
	b := NewBloom(1024, 4)
	b.Add([]byte("jti-1"))
	if !b.MightContain([]byte("jti-1")) {
		t.Fatal("should contain")
	}
	if b.MightContain([]byte("never-added-zzzz")) {
		t.Fatal("false positive on empty-ish key is unlikely but skip if flaky")
	}
}

func TestRevokeEmpty(t *testing.T) {
	clk := clock.NewFake(time.Unix(1, 0))
	r := New(clk, ModeExact)
	if err := r.Revoke("", time.Time{}); err == nil {
		t.Fatal("expected error")
	}
}
