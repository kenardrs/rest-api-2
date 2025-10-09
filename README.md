# REST API 2

[![Go Version](https://img.shields.io/badge/Go-1.25.1-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

Uma REST API simples para gerenciamento de usuários construída em Go, seguindo os princípios da Clean Architecture.

## 📋 Características

- ✅ Cadastro de usuários com validação de email único
- ✅ Listagem de todos os usuários
- ✅ Clean Architecture (handlers → usecases → repositories → models)
- ✅ Repositório em memória para simplicidade
- ✅ Geração automática de UUIDs para usuários
- ✅ Logging estruturado
- ✅ Validação de dados de entrada

## 🏗️ Arquitetura

O projeto segue os princípios da Clean Architecture com as seguintes camadas:

```
cmd/
├── api/           # Ponto de entrada da aplicação
└── client/        # Cliente de exemplo (futuro)

internal/
├── handlers/      # Camada de apresentação (HTTP handlers)
├── usecases/      # Regras de negócio da aplicação
├── repositories/  # Camada de acesso aos dados
│   └── users/     # Repositório específico de usuários
└── models/        # Entidades e DTOs
```

### Fluxo de Dados

```
HTTP Request → Handlers → UseCases → Repositories → Models
```

## 🚀 Início Rápido

### Pré-requisitos

- Go 1.25.1 ou superior
- Git

### Instalação

1. Clone o repositório:
```bash
git clone <repository-url>
cd rest-api-2
```

2. Baixe as dependências:
```bash
go mod tidy
```

3. Execute a aplicação:
```bash
go run cmd/api/main.go
```

A API estará disponível em `http://localhost:3000`

## 📚 API Endpoints

### Usuários

| Método | Endpoint | Descrição | Corpo da Requisição |
|--------|----------|-----------|-------------------|
| GET    | `/users` | Lista todos os usuários | - |
| POST   | `/users` | Cria um novo usuário | `{"name": "string", "email": "string"}` |

### Exemplos de Uso

#### Criar um usuário
```bash
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Silva",
    "email": "joao@email.com"
  }'
```

#### Listar usuários
```bash
curl http://localhost:3000/users
```

## 🧪 Estrutura do Projeto

```
rest-api-2/
├── README.md
├── go.mod
├── go.sum
├── cmd/
│   ├── api/
│   │   └── main.go              # Ponto de entrada da API
│   └── client/
│       └── main.go              # Cliente de exemplo
└── internal/
    ├── handlers/
    │   └── handlers.go          # HTTP handlers
    ├── usecases/
    │   └── usecases.go          # Lógica de negócio
    ├── repositories/
    │   ├── repositories.go      # Interfaces dos repositórios
    │   └── users/
    │       └── users.go         # Implementação do repositório de usuários
    └── models/
        └── models.go            # Modelos de dados
```

## 🛠️ Desenvolvimento

### Executar em modo de desenvolvimento

```bash
# Executar a API
go run cmd/api/main.go

# Ou com hot reload usando air (se instalado)
air
```

### Executar testes

```bash
# Executar todos os testes
go test ./...

# Executar testes com coverage
go test -cover ./...

# Executar testes com coverage detalhado
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Lint e formatação

```bash
# Formatar código
go fmt ./...

# Executar vet
go vet ./...

# Executar golangci-lint (se instalado)
golangci-lint run
```

## 📦 Dependências

- [github.com/google/uuid](https://github.com/google/uuid) - Geração de UUIDs

## 🏛️ Padrões de Design Utilizados

- **Clean Architecture**: Separação clara de responsabilidades entre camadas
- **Repository Pattern**: Abstração da camada de dados
- **Dependency Injection**: Injeção de dependências entre camadas
- **Use Cases**: Encapsulamento das regras de negócio

## 📝 Modelos de Dados

### User
```go
type User struct {
    ID    uuid.UUID `json:"id"`
    Name  string    `json:"name"`
    Email string    `json:"email"`
}
```

### CreateUserRequest
```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

## 🔧 Configuração

Atualmente, a aplicação utiliza configurações hardcoded. Futuras versões incluirão:

- Variáveis de ambiente
- Arquivo de configuração
- Configuração de banco de dados
- Configuração de logging

## 🚧 Roadmap

- [ ] Implementar endpoints HTTP completos
- [ ] Adicionar validação robusta de entrada
- [ ] Implementar middleware de logging
- [ ] Adicionar testes unitários
- [ ] Integração com banco de dados (PostgreSQL/MySQL)
- [ ] Documentação OpenAPI/Swagger
- [ ] Dockerização
- [ ] CI/CD pipeline
- [ ] Métricas e observabilidade

## 🤝 Contribuindo

1. Faça um fork do projeto
2. Crie sua feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para detalhes.

## 👥 Autores

- Seu Nome - [@seu-usuario](https://github.com/seu-usuario)

## 🙏 Agradecimentos

- Go community
- Clean Architecture principles by Robert C. Martin
- Todas as bibliotecas open source utilizadas

---

**Nota**: Este é um projeto de estudo/exemplo. Para uso em produção, considere implementar recursos adicionais de segurança, logging, monitoramento e testes mais robustos.