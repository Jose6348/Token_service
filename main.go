package main

import (
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	db := SetupDatabase()
	defer db.Close()

	r := mux.NewRouter()
	csrfMiddleware := csrf.Protect([]byte("32-byte-long-auth-key"))

	r.HandleFunc("/conta", CriarConta).Methods("POST")
	r.HandleFunc("/conta/{id}/saldo", ConsultarSaldo).Methods("GET")
	r.HandleFunc("/conta/{id}/deposito", DepositarDinheiro).Methods("POST")
	r.HandleFunc("/conta/{id}/saque", SacarDinheiro).Methods("POST")
	r.HandleFunc("/conta/transferencia", TransferirDinheiro).Methods("POST")
	r.HandleFunc("/conta/{id}", FecharConta).Methods("DELETE")

	log.Println("Servidor iniciado na porta 8080")
	http.ListenAndServe(":8080", csrfMiddleware(r))
}
