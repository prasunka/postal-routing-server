package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

type Payload struct {
	Forwardfrom string
	Forwardto   string
	Signature   string
}

// Taken from eth_sign_verify.go (https://gist.github.com/dcb9/385631846097e1f59e3cba3b1d42f3ed#file-eth_sign_verify-go)
func (app *application) verifySig(from, sigHex string, msg []byte) bool {
	sig, err := hexutil.Decode(sigHex)
	if err != nil {
		app.errorLog.Println("Invalid signature!")
		return false
	}

	msg = accounts.TextHash(msg)
	sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1

	recovered, err := crypto.SigToPub(msg, sig)
	if err != nil {
		return false
	}

	recoveredAddr := crypto.PubkeyToAddress(*recovered)

	return from == recoveredAddr.Hex()
}

func (app *application) createRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload Payload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("%v\n", payload)

	msg := []byte("I allow mails to be forwarded to " + payload.Forwardto)

	if app.verifySig(payload.Forwardfrom, payload.Signature, msg) {
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)

		id, err := app.endpoints.Insert(payload.Forwardto)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		_, err = app.routes.Insert(id, payload.Forwardfrom)

		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusCreated)
		resp["message"] = "Status Created"

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			app.errorLog.Printf("Error happened in JSON marshal. Err: %s", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)

		}

		w.Write(jsonResp)

	} else {
		http.Error(w, "Signature verification failed. Begone intruder!", http.StatusUnauthorized)
	}
}
