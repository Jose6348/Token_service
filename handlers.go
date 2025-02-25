package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

var db *sql.DB

func CriarConta(w http.ResponseWriter, r *http.Request) {
	var tipoConta string
	tipoConta = r.URL.Query().Get("tipo")

	switch tipoConta {
	case "fisica":
		criarContaFisica(w, r)
	case "juridica":
		criarContaJuridica(w, r)
	default:
		http.Error(w, "Tipo de conta inválido", http.StatusBadRequest)
	}
}

func criarContaFisica(w http.ResponseWriter, r *http.Request) {
	var pf PessoaFisica
	if err := json.NewDecoder(r.Body).Decode(&pf); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `
    INSERT INTO pessoa_fisica (renda_mensal, idade, nome_completo, celular, email, categoria, saldo)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id`
	id := 0
	err := db.QueryRow(sqlStatement, pf.RendaMensal, pf.Idade, pf.NomeCompleto, pf.Celular, pf.Email, pf.Categoria, pf.Saldo).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pf.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pf)
}

func criarContaJuridica(w http.ResponseWriter, r *http.Request) {
	var pj PessoaJuridica
	if err := json.NewDecoder(r.Body).Decode(&pj); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `
    INSERT INTO pessoa_juridica (faturamento, idade, nome_fantasia, celular, email_corporativo, categoria, saldo)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id`
	id := 0
	err := db.QueryRow(sqlStatement, pj.Faturamento, pj.Idade, pj.NomeFantasia, pj.Celular, pj.EmailCorporativo, pj.Categoria, pj.Saldo).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pj.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pj)
}

func ConsultarSaldo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var saldo float64
	err := db.QueryRow("SELECT saldo FROM pessoa_fisica WHERE id = $1 UNION SELECT saldo FROM pessoa_juridica WHERE id = $1", id).Scan(&saldo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"saldo": saldo})
}

func DepositarDinheiro(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var depósito struct {
		Valor float64 `json:"valor"`
	}
	if err := json.NewDecoder(r.Body).Decode(&depósito); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `
    UPDATE pessoa_fisica SET saldo = saldo + $1 WHERE id = $2
    UNION
    UPDATE pessoa_juridica SET saldo = saldo + $1 WHERE id = $2`
	_, err := db.Exec(sqlStatement, depósito.Valor, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func SacarDinheiro(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var saque struct {
		Valor float64 `json:"valor"`
	}
	if err := json.NewDecoder(r.Body).Decode(&saque); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `
    UPDATE pessoa_fisica SET saldo = saldo - $1 WHERE id = $2 AND saldo >= $1
    UNION
    UPDATE pessoa_juridica SET saldo = saldo - $1 WHERE id = $2 AND saldo >= $1`
	res, err := db.Exec(sqlStatement, saque.Valor, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, "Saldo insuficiente ou conta não encontrada", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func TransferirDinheiro(w http.ResponseWriter, r *http.Request) {
	var transferencia struct {
		De    int     `json:"de"`
		Para  int     `json:"para"`
		Valor float64 `json:"valor"`
	}
	if err := json.NewDecoder(r.Body).Decode(&transferencia); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer tx.Rollback()

	sqlSaque := `
    UPDATE pessoa_fisica SET saldo = saldo - $1 WHERE id = $2 AND saldo >= $1
    UNION
    UPDATE pessoa_juridica SET saldo = saldo - $1 WHERE id = $2 AND saldo >= $1`
	res, err := tx.Exec(sqlSaque, transferencia.Valor, transferencia.De)
	if err != nil || res.RowsAffected() == 0 {
		http.Error(w, "Saldo insuficiente ou conta de origem não encontrada", http.StatusBadRequest)
		return
	}

	sqlDeposito := `
    UPDATE pessoa_fisica SET saldo = saldo + $1 WHERE id = $2
    UNION
    UPDATE pessoa_juridica SET saldo = saldo + $1 WHERE id = $2`
	res, err = tx.Exec(sqlDeposito, transferencia.Valor, transferencia.Para)
	if err != nil || res.RowsAffected() == 0 {
		http.Error(w, "Conta de destino não encontrada", http.StatusBadRequest)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func FecharConta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	sqlStatement := `
    DELETE FROM pessoa_fisica WHERE id = $1
    UNION
    DELETE FROM pessoa_juridica WHERE id = $1`
	res, err := db.Exec(sqlStatement, id)
	if err != nil || res.RowsAffected() == 0 {
		http.Error(w, "Conta não encontrada", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
