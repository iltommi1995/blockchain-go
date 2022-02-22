package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
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
	return fmt.Sprintf("%064x%064x", s.R, s.S)
}

// Funzione per passare da String a 2 big.Int
func String2BigIntTuple(s string) (big.Int, big.Int) {
	bx, _ := hex.DecodeString(s[:64])
	by, _ := hex.DecodeString(s[64:])

	var bix big.Int
	var biy big.Int

	_ = bix.SetBytes(bx)
	_ = biy.SetBytes(by)

	return bix, biy
}

// Funzione che restituisce la signature a partire dalla versione string della signature
func SignatureFromString(s string) *Signature {
	x, y := String2BigIntTuple(s)
	return &Signature{&x, &y}
}

// Funzione che restituisce *ecdsa.PublicKey a partire dalla versione string della Public Key
func PublicKeyFromString(s string) *ecdsa.PublicKey {
	x, y := String2BigIntTuple(s)
	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     &x,
		Y:     &y,
	}
}

// Metodo che restituisce *ecdsa.PrivateKey a partire dalla versione string della Private Key
func PrivateKeyFromString(s string, publicKey *ecdsa.PublicKey) *ecdsa.PrivateKey {
	b, _ := hex.DecodeString(s[:])
	var bi big.Int
	_ = bi.SetBytes(b)
	return &ecdsa.PrivateKey{
		PublicKey: *publicKey,
		D:         &bi,
	}
}
