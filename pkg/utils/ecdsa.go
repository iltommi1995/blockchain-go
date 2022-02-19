package utils

import (
	"fmt"
	"math/big"
)

type Signature struct {
	// R sono le coordinate della nostra chiave pubblica
	R *big.Int
	// S viene calcolato in riferimento a informazioni come transactions hash e
	// la temporary public key per generare la signature
	S *big.Int
}

// to string della signature
func (s *Signature) String() string {
	return fmt.Sprintf("%x%x", s.R, s.S)
}
