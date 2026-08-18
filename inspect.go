package ticketgate

import (
	"github.com/LYH2263/go-ticketgate/internal/token"
)

func InspectKID(raw []byte) (string, byte, error) {
	kid, alg, err := token.PeekKID(raw)
	if err != nil {
		return "", 0, mapErr(err)
	}
	return kid, byte(alg), nil
}

func (g *Gateway) CurrentKID() (string, error) {
	k, err := g.ring.Current()
	if err != nil {
		return "", mapErr(err)
	}
	return k.Material.KID, nil
}
