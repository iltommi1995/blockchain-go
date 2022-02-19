package wallet_server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/iltommi1995/blockchain-go/pkg/utils"
	"github.com/iltommi1995/blockchain-go/pkg/wallet/transaction_request"
	"github.com/iltommi1995/blockchain-go/pkg/wallet/wallet"
)

// Path della directory dei template
const tempDir = "pkg/wallet_server/templates"

// Wallet server ha 2 proprietà:
// - porta su cui sarà in ascolto
// - gateway, che è l'url del blockchain server
type WalletServer struct {
	port    uint16
	gateway string
}

// Funzione per creare il wallet server
func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{port, gateway}
}

// Getter di port
func (ws *WalletServer) Port() uint16 {
	return ws.port
}

// Getter di gateway
func (ws *WalletServer) Gateway() string {
	return ws.gateway
}

// Resolver per l'endpoint "/"
func (ws *WalletServer) Index(w http.ResponseWriter, req *http.Request) {
	// Controllo il metodo
	switch req.Method {
	// Se è GET
	case http.MethodGet:
		// Parso il template
		t, _ := template.ParseFiles(path.Join(tempDir, "index.html"))
		// Restituisco il template in risposta
		t.Execute(w, "")
	default:
		// Per qualsiasi altro HTTP Method, dò errore
		log.Printf("ERROR: invalid HTTP Method")
	}
}

// Resolver per l'endpoint "/wallet"
func (ws *WalletServer) Wallet(w http.ResponseWriter, req *http.Request) {
	// Controllo il metodo
	switch req.Method {
	// Se è POST
	case http.MethodPost:
		// Setto gli headers
		w.Header().Add("Content-Type", "application/json")
		// Creo un nuovo wallet
		myWallet := wallet.NewWallet()
		// Restituisco, in json, i dati del wallet
		m, _ := myWallet.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		// Per qualsiasi altro HTTP Method, dò errore
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

// Resolver per l'endpoint "/transaction"
func (ws *WalletServer) CreateTransaction(w http.ResponseWriter, req *http.Request) {
	// Controllo il metodo
	switch req.Method {
	// Se è POST
	case http.MethodPost:
		// Faccio il decode del body della request,
		// contenente i dati della transazione
		decoder := json.NewDecoder(req.Body)
		// Transaction request
		var t transaction_request.TransactionRequest
		// Faccio il decode
		err := decoder.Decode(&t)
		// Se il decode fallisce, restituisco errore
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		// La transazione non è valida, restituisco errore
		if !t.Validate() {
			log.Println("ERROR: missing field(s)")
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		fmt.Println(*t.SenderPublicKey)
		fmt.Println(*t.SenderBloackchainAddress)
		fmt.Println(*t.SenderPrivateKey)
		fmt.Println(*t.RecipientBlockchainAddress)
		fmt.Println(*t.Value)
	default:
		// Per qualsiasi altro HTTP Method, dò errore
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

// Funzione per avviare il server
func (ws *WalletServer) Run() {
	// Qui si creano gli endpoint e si associano i resolver
	http.HandleFunc("/", ws.Index)
	http.HandleFunc("/wallet", ws.Wallet)
	http.HandleFunc("/transaction", ws.CreateTransaction)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa((int(ws.port))), nil))
}
