package blockchain_server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/iltommi1995/blockchain-go/pkg/blockchain/blockchain"
	"github.com/iltommi1995/blockchain-go/pkg/wallet/wallet"
)

// File nella cache in cui viene salvata la blockchain
var cache map[string]*blockchain.Blockchain = make(map[string]*blockchain.Blockchain)

// Blockchain server, ha solo la porta su cui runna
// così posso startarne di più su porte diverse
type BlockchainServer struct {
	port uint16
}

// Funzione per creare un nuovo server
func NewBlockchainServer(port uint16) *BlockchainServer {
	return &BlockchainServer{port}
}

// Getter della porta
func (bcs *BlockchainServer) Port() uint16 {
	return bcs.port
}

// Metodo per avere la blockchain
func (bcs *BlockchainServer) GetBloackchain() *blockchain.Blockchain {
	// Recupero la blockchain dalla cache
	bc, ok := cache["blockchain"]

	// Se non c'è
	if !ok {
		// Creo wallet del miner
		minerWallet := wallet.NewWallet()
		// Passiamo anche la porta perché poi servirà alla blockchain per cercare
		// altri nodi
		// Creo nuova blockchain
		bc := blockchain.NewBlockchain(minerWallet.BlockchainAddress(), bcs.Port())
		// Salvo nella cache
		cache["blockchain"] = bc

		log.Printf("private_key %v", minerWallet.PrivateKeyStr())
		log.Printf("public_key %v", minerWallet.PublicKeyStr())
		log.Printf("blockchain_address %v", minerWallet.BlockchainAddress())
	}

	// Ritorno la blockchain
	return bc
}

func HelloWorld(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, World!")
}

// Resolver per endpoint "/"
func (bcs *BlockchainServer) GetChain(w http.ResponseWriter, req *http.Request) {
	// Controlla il metodo della request
	switch req.Method {
	// Se è GET
	case http.MethodGet:
		// Setto headers
		w.Header().Add("Content-Type", "application/json")
		// Blockchain
		bc := bcs.GetBloackchain()
		// RIGA DA CONTROLLARE (json della blockchain, fatto in questo modo alla prima chiamata dà null)
		m, _ := json.Marshal(bc)

		// Restituisco la blockchain
		io.WriteString(w, string(m[:]))
	default:
		// Se è altro HTTP Method do errore
		log.Printf("ERROR: invalid HTTP Method")
	}
}

// Metodo per avviare il server
func (bcs *BlockchainServer) Run() {
	// Crea endpoint e associa resolver
	http.HandleFunc("/", bcs.GetChain)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(int(bcs.Port())), nil))
}
