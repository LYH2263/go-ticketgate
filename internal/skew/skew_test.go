package skew

import (
	"testing"
	"time"
)

func TestBoundary(t *testing.T) {
	now := time.Unix(100, 0)
	exp := time.Unix(100, 0)
	if Expired(now, exp, 0) {
		t.Fatal("now==exp should be valid")
	}
	if !Expired(now.Add(time.Second), exp, 0) {
		t.Fatal("now>exp should expire")
	}
	nbf := time.Unix(100, 0)
	if TooEarly(now, nbf, 0) {
		t.Fatal("now==nbf should be valid")
	}
	if !TooEarly(now.Add(-time.Second), nbf, 0) {
		t.Fatal("now<nbf should be early")
	}
	if Expired(now, exp, 2*time.Second) {
		t.Fatal("skew should keep valid")
	}
}
