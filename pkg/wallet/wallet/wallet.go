package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

// Il wallet ha:
// - private key
// - public key
// - blockchain address
type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockchainAddress string
}

// Funzione per creare nuovo wallet
func NewWallet() *Wallet {
	w := new(Wallet)
	// Ottengo la private key randomica con la libreria ecdsa
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	w.privateKey = privateKey

	// Passaggi per generare l'address:

	// 1. Dalla private key ottengo la public key
	w.publicKey = &w.privateKey.PublicKey

	// Creare address
	// 2. SHA-256 della chiave pubblica (32 bytes)
	h2 := sha256.New()
	h2.Write(w.publicKey.X.Bytes())
	h2.Write(w.publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)

	// 3. Hash function RIPEMD-160 sul risultato di SHA-256 (20 byres)
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)

	// 4. Aggiungere version byte davanti all'hash RIPEMD-160 (0x00 per mainnet)
	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:])

	// 5 - Si fa SHA-256 del risultato ottenuto nel punto 4
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)

	// 6. Si fa SHA-256 del risultato del punto 5
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)

	// 7. Prendere i primi 4 bytes del risultato del punto 5
	// come checksum
	checksum := digest6[:4]

	// 8. Agiungere i 4 bytes di checksum del punto 7 alla fine
	// del output del punto 4 (), si arriva quindi a 25 byes
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[21:], checksum[:])

	// 9. Convertire il risultato da byte string a base58
	address := base58.Encode(dc8)

	w.blockchainAddress = address
	return w
}

// è una specie di getter, alternativa a scrivere in
// maiuscolo la prima lettera delle proprietà di
// uno struct
func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

// Getter di private key come string
func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

// Getter di public key
func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

// Getter di public key come string
func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%064x%064x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

// Getter di blockchain address
func (w *Wallet) BlockchainAddress() string {
	return w.blockchainAddress
}

// Json del wallet
func (w *Wallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PrivateKey        string `json:"private_key"`
		PublicKey         string `json:"public_key"`
		BlockchainAddress string `json:"blockchain_address"`
	}{
		PrivateKey:        w.PrivateKeyStr(),
		PublicKey:         w.PublicKeyStr(),
		BlockchainAddress: w.BlockchainAddress(),
	})
}
