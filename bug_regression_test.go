package ticketgate

import (
	"errors"
	"testing"
	"time"

	"github.com/LYH2263/go-ticketgate/internal/clock"
)

type failSavePersist struct {
	saved []string
	err   error
}

func (p *failSavePersist) Save(jtis []string) error {
	if p.err != nil {
		return p.err
	}
	p.saved = append([]string(nil), jtis...)
	return nil
}

func (p *failSavePersist) Load() ([]string, error) {
	return append([]string(nil), p.saved...), nil
}

func TestBug08_RevokePersistErrorRollsBack(t *testing.T) {
	clk := clock.NewFake(time.Unix(1_700_000_000, 0))
	ps := &failSavePersist{err: errors.New("disk full")}
	g := newGW(t, clk, WithRevokePersist(ps))
	issued, err := g.IssueDetailed(Claims{Subject: "u", ID: "jti-persist-1"})
	if err != nil {
		t.Fatal(err)
	}
	err = g.Revoke(issued.Claims.ID)
	if err == nil {
		t.Fatal("expected persist error")
	}
	if g.IsRevoked(issued.Claims.ID) {
		t.Fatal("memory must roll back when persist fails")
	}
	if _, err := g.Verify(issued.Token); err != nil {
		t.Fatalf("ticket should still verify after failed persist: %v", err)
	}
	g2 := newGW(t, clk, WithRevokePersist(ps))
	if g2.IsRevoked(issued.Claims.ID) {
		t.Fatal("reopen must not show unpersisted revoke")
	}
}
