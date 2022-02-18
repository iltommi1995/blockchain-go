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

var cache map[string]*blockchain.Blockchain = make(map[string]*blockchain.Blockchain)

type BlockchainServer struct {
	port uint16
}

func NewBlockchainServer(port uint16) *BlockchainServer {
	return &BlockchainServer{port}
}

func (bcs *BlockchainServer) Port() uint16 {
	return bcs.port
}

func (bcs *BlockchainServer) GetBloackchain() *blockchain.Blockchain {
	bc, ok := cache["blockchain"]

	if !ok {
		minerWallet := wallet.NewWallet()
		// Passiamo anche la porta perché poi servirà alla blockchain per cercare
		// altri nodi
		bc := blockchain.NewBlockchain(minerWallet.BlockchainAddress(), bcs.Port())
		cache["blockchain"] = bc

		log.Printf("private_key %v", minerWallet.PrivateKeyStr())
		log.Printf("public_key %v", minerWallet.PublicKeyStr())
		log.Printf("blockchain_address %v", minerWallet.BlockchainAddress())
	}
	return bc
}

func HelloWorld(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, World!")
}

func (bcs *BlockchainServer) GetChain(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcs.GetBloackchain()
		// RIGA DA CONTROLLARE
		m, _ := json.Marshal(bc)

		io.WriteString(w, string(m[:]))
	default:
		log.Printf("ERROR: invalid HTTP Method")
	}
}

func (bcs *BlockchainServer) Run() {
	http.HandleFunc("/", bcs.GetChain)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(int(bcs.Port())), nil))
}
