package blockchain_server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/iltommi1995/blockchain-go/pkg/blockchain/amount_response"
	"github.com/iltommi1995/blockchain-go/pkg/blockchain/blockchain"
	blockchain_transaction "github.com/iltommi1995/blockchain-go/pkg/blockchain/transaction"
	blockchain_transaction_request "github.com/iltommi1995/blockchain-go/pkg/blockchain/transaction_request"
	"github.com/iltommi1995/blockchain-go/pkg/utils"
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

		// Controllare catena
		/*if len(bc.Chain()) == 1 && len(bc.Neighbors()) > 0 {
			fmt.Println("Risolti i conflitti")
			bc.ResolveConflicts()
		}*/

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
		fmt.Println("Stampo la blockchain prima di inviarla")
		bc.Print()
		// RIGA DA CONTROLLARE (json della blockchain, fatto in questo modo alla prima chiamata dà null)
		m, _ := json.Marshal(bc)

		// Restituisco la blockchain
		io.WriteString(w, string(m[:]))
	default:
		// Se è altro HTTP Method do errore
		log.Printf("ERROR: invalid HTTP Method")
	}
}

// Resolver dell'endpoint "/transactions"
func (bcs *BlockchainServer) Transactions(w http.ResponseWriter, req *http.Request) {
	// Controllo il metodo
	switch req.Method {
	// Se è GET
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcs.GetBloackchain()
		transactions := bc.TransactionPool()
		// Restituisco il json contenente il transaction pool della blockchain
		// e il numero delle transazioni contenute nel transaction pool
		m, _ := json.Marshal(struct {
			Transactions []*blockchain_transaction.Transaction `json:"transactions"`
			Length       int                                   `json:"length"`
		}{
			Transactions: transactions,
			Length:       len(transactions),
		})
		io.WriteString(w, string(m[:]))
	case http.MethodPost:
		// Se è POST
		// Decodifico il json del body della request
		decoder := json.NewDecoder(req.Body)
		// Creo una transaction request lato server
		var t blockchain_transaction_request.TransactionRequest
		// faccio il decode di t
		err := decoder.Decode(&t)
		// Se c'è un errore
		if err != nil {
			// Dico che è fallita la costruzione della transazione
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		// Valido la Transaction Request, vedendo se ci sono tutti i dati
		// necessari, se no do errore
		if !t.Validate() {
			log.Println("ERROR: missing field(s)")
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		// Prendo la publicKey e la trasformo da String a *ecdsa.PublicKey
		publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
		// Prendo la signature e la trasformo da String a utils.Signature
		signature := utils.SignatureFromString(*t.Signature)

		bc := bcs.GetBloackchain()

		// Creo la transazione lato server
		isCreated := bc.CreateTransaction(
			*t.SenderBlockchainAddress,
			*t.RecipientBlockchainAddress,
			*t.Value,
			publicKey,
			signature,
		)

		w.Header().Add("Content-Type", "application/json")
		// Creo una variabile m per la response
		var m []byte
		// Controllo se la transazione è stata creata correttamente
		if !isCreated {
			// Se no restituisco il messaggio di errore
			w.WriteHeader(http.StatusBadRequest)
			m = utils.JsonStatus("fail")
		} else {
			// Se sì restituisco il messaggio di successo
			w.WriteHeader(http.StatusCreated)
			m = utils.JsonStatus("success")
		}
		io.WriteString(w, string(m))
	// Se è metodo PUT
	case http.MethodPut:
		decoder := json.NewDecoder(req.Body)
		var t blockchain_transaction_request.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		if !t.Validate() {
			log.Println("ERROR: missing field(s)")
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
		signature := utils.SignatureFromString(*t.Signature)

		bc := bcs.GetBloackchain()
		// Qui le modifiche fatte
		isUpdated := bc.AddTransaction(
			*t.SenderBlockchainAddress,
			*t.RecipientBlockchainAddress,
			*t.Value,
			publicKey,
			signature,
		)

		w.Header().Add("Content-Type", "application/json")
		var m []byte
		if !isUpdated {
			w.WriteHeader(http.StatusBadRequest)
			m = utils.JsonStatus("fail")
		} else {
			m = utils.JsonStatus("success")
		}
		io.WriteString(w, string(m))
	case http.MethodDelete:
		bc := bcs.GetBloackchain()
		bc.ClearTransactionPool()
		io.WriteString(w, string(utils.JsonStatus("success")))

	default:
		// Se è un altro metodo
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

// Resolver dell'endpoint "/mine"
func (bcs *BlockchainServer) Mine(w http.ResponseWriter, req *http.Request) {
	// Controllo il metodo HTTP
	switch req.Method {
	// Se è GET
	case http.MethodGet:
		bc := bcs.GetBloackchain()
		// Effettuo il imining
		isMined := bc.Mining()

		// Controllo se il mining è avvenuto correttamente e creo la response
		var m []byte
		if !isMined {
			w.WriteHeader(http.StatusBadRequest)
			m = utils.JsonStatus("fail")
		} else {
			m = utils.JsonStatus("success")
		}
		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(m))
	// Se è un altro metodo
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

// Resolver dell'endpoint "/mine/start"
// Serve ad automatizzare il processo di mining
func (bcs *BlockchainServer) StartMine(w http.ResponseWriter, req *http.Request) {
	// Controllo il metodo HTTP
	switch req.Method {
	// Se è GET
	case http.MethodGet:
		bc := bcs.GetBloackchain()
		// Avvio il Mining
		bc.StartMining()
		// Costruisco la response
		m := utils.JsonStatus("success")
		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(m))
	// Se è un altro metodo
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

// Resolver dell'endpoint "/amount"
func (bcs *BlockchainServer) Amount(w http.ResponseWriter, req *http.Request) {
	// Controllo il metodo HTTP
	switch req.Method {
	// Se è GET
	case http.MethodGet:
		// Recupero il query param con il blockchain address
		blockchainAddress := req.URL.Query().Get("blockchain_address")
		// Recupero il bilancio
		amount := bcs.GetBloackchain().CalculateTotalAmount(blockchainAddress)
		// Preparo la risposta
		ar := &amount_response.AmountResponse{Amount: amount}
		m, _ := ar.MarshalJSON()
		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(m[:]))
	// Se è un altro metodo
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) Consensus(w http.ResponseWriter, req *http.Request) {
	// Controllo il metodo HTTP
	switch req.Method {
	// Se è PUT
	case http.MethodPut:
		bc := bcs.GetBloackchain()
		replaced := bc.ResolveConflicts()

		w.Header().Add("Content-Type", "application/json")
		if replaced {
			io.WriteString(w, string(utils.JsonStatus("success")))
		} else {
			io.WriteString(w, string(utils.JsonStatus("fail")))
		}
	// Se è un altro metodo
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

// Metodo per avviare il server
func (bcs *BlockchainServer) Run() {
	bcs.GetBloackchain()
	log.Println("Running blockchain...")
	bcs.GetBloackchain().Run()

	// Crea endpoint e associa resolver
	http.HandleFunc("/", bcs.GetChain)
	http.HandleFunc("/transactions", bcs.Transactions)
	http.HandleFunc("/mine", bcs.Mine)
	http.HandleFunc("/mine/start", bcs.StartMine)
	http.HandleFunc("/amount", bcs.Amount)
	http.HandleFunc("/consensus", bcs.Consensus)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(int(bcs.Port())), nil))
}
