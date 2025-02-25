CREATE TABLE pessoa_fisica (
    id SERIAL PRIMARY KEY,
    renda_mensal DECIMAL,
    idade INT,
    nome_completo VARCHAR(255),
    celular VARCHAR(20),
    email VARCHAR(255),
    categoria VARCHAR(50),
    saldo DECIMAL
);

CREATE TABLE pessoa_juridica (
    id SERIAL PRIMARY KEY,
    faturamento DECIMAL,
    idade INT,
    nome_fantasia VARCHAR(255),
    celular VARCHAR(20),
    email_corporativo VARCHAR(255),
    categoria VARCHAR(50),
    saldo DECIMAL
);