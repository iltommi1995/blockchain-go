package wallet_server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/iltommi1995/blockchain-go/pkg/blockchain/amount_response"
	blockchain_transaction_request "github.com/iltommi1995/blockchain-go/pkg/blockchain/transaction_request"
	"github.com/iltommi1995/blockchain-go/pkg/utils"
	wallet_transaction "github.com/iltommi1995/blockchain-go/pkg/wallet/transaction"
	wallet_transaction_request "github.com/iltommi1995/blockchain-go/pkg/wallet/transaction_request"
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
		var t wallet_transaction_request.TransactionRequest
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
		/*
			fmt.Println(len(*t.SenderPublicKey))
			fmt.Println(len(*t.SenderBloackchainAddress))
			fmt.Println(len(*t.SenderPrivateKey))
			fmt.Println(*t.RecipientBlockchainAddress)
			fmt.Println(*t.Value)
		*/

		// Converto la Public Key da stringa ecdsa.PublicKey
		publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
		// Converto la Private Key da stringa ecdsa.PrivateKey
		privateKey := utils.PrivateKeyFromString(*t.SenderPrivateKey, publicKey)
		// Converto il value in un float
		value, err := strconv.ParseFloat(*t.Value, 32)
		// In caso di errore dico che c'è stato un errore di parsing
		if err != nil {
			log.Println("ERROR: parse error")
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		//
		value32 := float32(value)
		/*
			fmt.Println(publicKey)
			fmt.Println(privateKey)
			fmt.Printf("%.1f\n", value32)
		*/

		w.Header().Add("Content-Type", "application/json")

		// Creo una transazione lato wallet, passando i dati necessari
		transaction := wallet_transaction.NewTransaction(
			privateKey,
			publicKey,
			*t.SenderBloackchainAddress,
			*t.RecipientBlockchainAddress,
			value32)
		// Creo la signature della transaction
		signature := transaction.GenerateSignature()
		// Versione string della signature
		signatureStr := signature.String()

		// Creo transaction request lato server
		bt := &blockchain_transaction_request.TransactionRequest{
			SenderBlockchainAddress:    t.SenderBloackchainAddress,
			RecipientBlockchainAddress: t.RecipientBlockchainAddress,
			SenderPublicKey:            t.SenderPublicKey,
			Value:                      &value32,
			Signature:                  &signatureStr,
		}
		// Converto in json la transaction request
		m, _ := json.Marshal(bt)
		// Creo un buffer della transaction request
		buf := bytes.NewBuffer(m)
		// Faccio una post request all'endpoint del blockchain server, inviando i dati
		// della transaction request
		resp, _ := http.Post(ws.Gateway()+"/transactions", "application/json", buf)
		// Se lo status code è 201, restituisco success
		if resp.StatusCode == 201 {
			io.WriteString(w, string(utils.JsonStatus("success")))
			return
		}
		// Altrimenti restituisco fail
		io.WriteString(w, string(utils.JsonStatus("fail")))
	default:
		// Per qualsiasi altro HTTP Method, dò errore
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

// Resolver dell'endpoint "/wallet/amount"
func (ws *WalletServer) WalletAmount(w http.ResponseWriter, req *http.Request) {
	// Controllo il metodo
	switch req.Method {
	// Se è GET
	case http.MethodGet:
		// Recupero l'indirizzo dal query param
		blockchainAddress := req.URL.Query().Get("blockchain_address")
		// Preparo l'endpoint corretto
		endpoint := fmt.Sprintf("%s/amount", ws.Gateway())

		// Preparo il client
		client := &http.Client{}
		// Creo GET request all'endpoint
		bcsReq, _ := http.NewRequest("GET", endpoint, nil)
		// Aggiungo i query param
		q := bcsReq.URL.Query()
		q.Add("blockchain_address", blockchainAddress)
		// Faccio l'encode dell'url
		bcsReq.URL.RawQuery = q.Encode()

		// Creo la response, inviando la request
		bcsResp, err := client.Do(bcsReq)

		// Se ci sono errori, restituisco fail
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		w.Header().Add("Content-Type", "application/json")
		// Se lo status code della response è 200
		if bcsResp.StatusCode == 200 {
			// Faccio decode del body della response
			decoder := json.NewDecoder(bcsResp.Body)
			var bar amount_response.AmountResponse
			// Decode in amount_response
			err := decoder.Decode(&bar)

			// Se c'è un errore nel decoding restituisco l'errore
			if err != nil {
				log.Printf("ERROR: %v", err)
				io.WriteString(w, string(utils.JsonStatus("fail")))
				return
			}

			// Creo il json da dare in risposta
			m, _ := json.Marshal(struct {
				Message string  `json:"message"`
				Amount  float32 `json:"amount"`
			}{
				Message: "success",
				Amount:  bar.Amount,
			})
			// Risposta
			io.WriteString(w, string(m[:]))
			// Se lo stato non è 200 restituisco fail
		} else {
			io.WriteString(w, string(utils.JsonStatus("fail")))
		}
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
	http.HandleFunc("/wallet/amount", ws.WalletAmount)
	http.HandleFunc("/transaction", ws.CreateTransaction)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa((int(ws.port))), nil))
}
