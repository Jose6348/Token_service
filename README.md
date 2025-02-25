Token Service 🔐

Um serviço de autenticação baseado em tokens, desenvolvido em Go. Este projeto fornece uma estrutura simples e eficiente para gerenciar autenticação e autorização, permitindo a geração, validação e renovação de tokens de acesso.
🚀 Recursos

    🛠 Autenticação segura baseada em tokens JWT
    📦 Estrutura modular, facilitando a escalabilidade
    📊 Banco de dados integrado via database.go
    🌐 Handlers para requisições HTTP

📂 Estrutura do Projeto

    main.go → Ponto de entrada da aplicação
    handlers.go → Lida com as requisições HTTP
    database.go → Conexão com o banco de dados
    models.go → Definições dos modelos de dados
    migrations/ → Scripts para gerenciamento do banco

📌 Como Usar

    Clone o repositório:

git clone https://github.com/Jose6348/Token_service.git
cd Token_service

Instale as dependências:

go mod tidy

Execute o serviço:

go run main.go
