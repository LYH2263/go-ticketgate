package keyring

import (
	"testing"
	"time"

	"github.com/LYH2263/go-ticketgate/internal/clock"
	"github.com/LYH2263/go-ticketgate/internal/token"
)

func TestRotateGrace(t *testing.T) {
	clk := clock.NewFake(time.Unix(1000, 0))
	r, err := New(clk, 5*time.Second, token.AlgHS256)
	if err != nil {
		t.Fatal(err)
	}
	cur, err := r.Current()
	if err != nil {
		t.Fatal(err)
	}
	old := cur.Material.KID
	if _, err := r.Rotate(); err != nil {
		t.Fatal(err)
	}
	if !r.InGrace(old) {
		t.Fatal("expected grace")
	}
	if _, err := r.LookupVerify(old); err != nil {
		t.Fatal(err)
	}
	clk.Advance(5 * time.Second)
	if !r.Retired(old) {
		t.Fatal("expected retired")
	}
	if _, err := r.LookupVerify(old); err == nil {
		t.Fatal("retired key should not verify")
	}
}

func TestGracePeriod(t *testing.T) {
	clk := clock.NewFake(time.Unix(1, 0))
	r, err := New(clk, time.Minute, token.AlgHS256)
	if err != nil {
		t.Fatal(err)
	}
	if r.GracePeriod() != time.Minute {
		t.Fatal(r.GracePeriod())
	}
	r.SetGrace(2 * time.Second)
	if r.GracePeriod() != 2*time.Second {
		t.Fatal(r.GracePeriod())
	}
}
