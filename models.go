package main

type PessoaFisica struct {
	ID           int     `json:"id"`
	RendaMensal  float64 `json:"renda_mensal"`
	Idade        int     `json:"idade"`
	NomeCompleto string  `json:"nome_completo"`
	Celular      string  `json:"celular"`
	Email        string  `json:"email"`
	Categoria    string  `json:"categoria"`
	Saldo        float64 `json:"saldo"`
}

type PessoaJuridica struct {
	ID               int     `json:"id"`
	Faturamento      float64 `json:"faturamento"`
	Idade            int     `json:"idade"`
	NomeFantasia     string  `json:"nome_fantasia"`
	Celular          string  `json:"celular"`
	EmailCorporativo string  `json:"email_corporativo"`
	Categoria        string  `json:"categoria"`
	Saldo            float64 `json:"saldo"`
}
