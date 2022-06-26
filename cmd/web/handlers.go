package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/prasunka/postal-routing-server/pkg/models/mysql"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

type Payload struct {
	Forwardfrom string
	Forwardto   string
	Signature   string
}

func createRoute(w http.ResponseWriter, r *http.Request) {
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

	routes := mysql.RouteModel{DB: DB}

	routes.Insert(1, payload.Forwardfrom)

	fmt.Printf("%v\n", payload)

}
