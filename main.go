package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	csrfToken := csrf.Token(r)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<html><body><h1>Servidor Protegido CSRF</h1><form action='/submit' method='POST'><input type='hidden' name='csrf_token' value='%s'><input type='text' name='data'><input type='submit' value='Enviar'></form></body></html>", csrfToken)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	data := r.FormValue("data")
	fmt.Fprintf(w, "Dados recebidos com sucesso: %s", data)
}

func main() {
	// Criando um middleware CSRF
	csrfMiddleware := csrf.Protect([]byte("32-byte-secret-key"), csrf.Secure(false))

	// Criando roteador
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/submit", submitHandler).Methods("POST")

	// Aplicando middleware CSRF
	http.Handle("/", csrfMiddleware(r))

	log.Println("Servidor rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
