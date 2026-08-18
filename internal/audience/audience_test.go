package audience

import "testing"

func TestGlob(t *testing.T) {
	if !MatchGlob("api.*", "api.v1") {
		t.Fatal("api.*")
	}
	if MatchGlob("api.*", "web.v1") {
		t.Fatal("should not match")
	}
	if !MatchGlob("a?c", "abc") {
		t.Fatal("?")
	}
}

func TestMatchAudience(t *testing.T) {
	if !MatchAudience([]string{"api"}, []string{"api", "web"}, false) {
		t.Fatal("intersect")
	}
	if MatchAudience([]string{"api"}, []string{"web"}, false) {
		t.Fatal("mismatch")
	}
	if !MatchAudience(nil, []string{"x"}, false) {
		t.Fatal("empty want")
	}
}

func TestConstraint(t *testing.T) {
	c := Constraint{Issuer: "tg", Audience: []string{"api"}}
	iss, sub, aud := c.Check("tg", "u", []string{"api"})
	if !iss || !sub || !aud {
		t.Fatal(iss, sub, aud)
	}
	iss, _, aud = c.Check("other", "u", []string{"web"})
	if iss || aud {
		t.Fatal("should fail")
	}
}
