# REST API 2

[![Go Version](https://img.shields.io/badge/Go-1.25.1-blue.svg)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-blue.svg)](https://postgresql.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

Uma REST API robusta para gerenciamento de usuários construída em Go, seguindo os princípios da Clean Architecture com **integração completa ao PostgreSQL**.

> 🚀 **v1.1.0** - Agora com persistência em PostgreSQL, connection pooling e configuração via ambiente!

## 📋 Características

- ✅ Cadastro de usuários com validação de email único
- ✅ Listagem de todos os usuários
- ✅ Clean Architecture (handlers → usecases → repositories → models)
- ✅ **Integração com PostgreSQL** para persistência de dados
- ✅ **Connection pooling** e gerenciamento de recursos
- ✅ **Migrations automáticas** para estrutura do banco
- ✅ Geração automática de UUIDs para usuários
- ✅ Logging estruturado com slog
- ✅ Validação de dados de entrada e tratamento de erros
- ✅ Cliente HTTP funcional com tratamento de erros
- ✅ Respostas padronizadas de sucesso e erro
- ✅ **Configuração via variáveis de ambiente**
- ✅ **Graceful shutdown** com limpeza de recursos

## 🏗️ Arquitetura

O projeto segue os princípios da Clean Architecture com as seguintes camadas:

```
cmd/
├── api/           # Ponto de entrada da aplicação
└── client/        # Cliente de exemplo para testes

internal/
├── config/        # Configurações e variáveis de ambiente
├── database/      # Conexão e pool do PostgreSQL
├── handlers/      # Camada de apresentação (HTTP handlers)
├── usecases/      # Regras de negócio da aplicação
├── repositories/  # Camada de acesso aos dados (PostgreSQL)
│   └── users/     # Repositório específico de usuários
└── models/        # Entidades e DTOs

migrations/        # Scripts SQL para estrutura do banco
.env              # Variáveis de ambiente (desenvolvimento)
docker-compose.yml # Documentação do PostgreSQL
```

### Fluxo de Dados

```
HTTP Request → Handlers → UseCases → Repositories → PostgreSQL
                                                   ↓
                                               Models (Data)
```

## 🏛️ Clean Architecture Detalhada

### 📊 Fluxo de Dependências

```
┌─────────────┐    depends on    ┌──────────────┐
│   Handlers  │ ─────────────→   │   UseCases   │
│(Controllers)│                  │ (Services)   │
└─────────────┘                  └──────────────┘
                                        ↓
                                   depends on
                                        ↓
                                 ┌──────────────┐    depends on    ┌─────────┐
                                 │ Repositories │ ─────────────→   │ Models  │
                                 │   (Data)     │                  │(Entities│
                                 └──────────────┘                  └─────────┘
```

### 🎯 Responsabilidades por Camada

#### 1. 📱 **Handlers** (Camada de Apresentação)
- ✅ **Fazer:** Receber requisições HTTP, validar formato, retornar respostas
- ❌ **Não fazer:** Lógica de negócio, acesso direto aos dados

```go
func (h Handlers) addUsers(w http.ResponseWriter, r *http.Request) {
    // 1. Decodifica JSON com tratamento de erro
    var req models.CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(models.ErrorResponse{Reason: err.Error()})
        return
    }
    
    // 2. Chama UseCase (delega a lógica)
    userID, err := h.useCases.Add(req)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(models.ErrorResponse{Reason: err.Error()})
        return
    }
    
    // 3. Retorna resposta de sucesso
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(models.CreateUserResponse{NewUserID: userID})
}
```

#### 2. 🧠 **UseCases** (Camada de Negócio)
- ✅ **Fazer:** Regras de negócio, orquestração, validações de domínio
- ❌ **Não fazer:** Detalhes HTTP, queries específicas de banco

```go
func (u UseCases) Add(newUser models.User) (uuid.UUID, error) {
    // Regra de negócio: email deve ser único
    if exists := u.repos.User.EmailExists(newUser.Email); exists {
        return uuid.Nil, errors.New("email already exists")
    }
    
    // Lógica de negócio: gerar ID
    newUser.ID = uuid.New()
    u.repos.User.Add(newUser)
    
    return newUser.ID, nil
}
```

#### 3. 💾 **Repositories** (Camada de Dados)
- ✅ **Fazer:** Persistência, queries, cache, APIs externas
- ❌ **Não fazer:** Regras de negócio, validações de domínio

```go
func (u *Users) Add(newUser models.User) {
    u.users = append(u.users, newUser) // Simples persistência
}

func (u *Users) EmailExists(email string) bool {
    for _, user := range u.users {
        if user.Email == email { return true }
    }
    return false
}
```

#### 4. 📦 **Models** (Camada de Entidades)
- ✅ **Fazer:** Estruturas de dados, validações básicas
- ❌ **Não fazer:** Depender de frameworks ou bibliotecas externas

```go
type User struct {
    ID    uuid.UUID `json:"id"`
    Name  string    `json:"name"`
    Email string    `json:"email"`
}
```

### 🎯 **Princípios da Clean Architecture**

#### **Dependency Inversion Principle**
- Camadas externas dependem das internas
- UseCases não sabem sobre HTTP ou banco de dados
- Repositories implementam interfaces definidas pelos UseCases

#### **Separation of Concerns**
- Cada camada tem uma responsabilidade específica
- Mudanças em uma camada não afetam outras
- Fácil de testar isoladamente

#### **Independence**
- Framework independente (pode trocar de HTTP para gRPC)
- Database independente (pode trocar de PostgreSQL para MongoDB)
- UI independente (web, mobile, CLI)

### ✅ **Benefícios desta Arquitetura**

1. **🔄 Testabilidade:** Cada camada pode ser testada isoladamente
2. **🔌 Flexibilidade:** Fácil trocar implementações (banco, framework)
3. **📦 Manutenibilidade:** Código organizado e fácil de entender
4. **👥 Escalabilidade:** Times podem trabalhar em paralelo nas camadas
5. **🛡️ Robustez:** Mudanças em detalhes não afetam regras de negócio

### 🧪 **Exemplo de Teste Isolado**

```go
// Testando UseCase sem HTTP nem banco de dados
func TestAddUser_EmailAlreadyExists(t *testing.T) {
    // Mock repository
    mockRepo := &MockUserRepository{
        emailExists: true, // Simula email já existe
    }
    
    useCase := NewUseCases(&Repositories{User: mockRepo})
    
    // Teste isolado da lógica de negócio
    _, err := useCase.Add(models.User{Email: "test@email.com"})
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "email already exists")
}
```

## 🚀 Início Rápido

### Pré-requisitos

- Go 1.25.1 ou superior
- Docker e Docker Compose
- PostgreSQL 15+ (ou use o container fornecido)
- Git

### Setup do Banco de Dados

#### Opção 1: PostgreSQL já existente
Se você já tem PostgreSQL rodando, ajuste o arquivo `.env`:

```bash
# .env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=seu_database
DB_USER=seu_usuario
DB_PASSWORD=sua_senha
DB_SSLMODE=disable
```

#### Opção 2: Usar Docker
```bash
# PostgreSQL via Docker
docker run -d \
  --name postgres-rest-api \
  -e POSTGRES_DB=app_db \
  -e POSTGRES_USER=app_user \
  -e POSTGRES_PASSWORD=example \
  -p 5432:5432 \
  postgres:16-alpine
```

### Executar Migrations
```bash
# Criar tabelas no banco
docker exec -i <container_name> psql -U <user> -d <database> < migrations/001_create_users_table.sql

# Exemplo:
docker exec -i db psql -U app_user -d app_db < migrations/001_create_users_table.sql
```
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

# Resposta de sucesso (201 Created):
{
  "newUserId": "550e8400-e29b-41d4-a716-446655440000"
}

# Resposta de erro (400 Bad Request - email já existe):
{
  "reason": "email already exists"
}
```

#### Listar usuários
```bash
curl http://localhost:3000/users

# Resposta (200 OK):
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "João Silva",
    "email": "joao@email.com"
  }
]
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
    │   ├── handlers.go          # Core handlers e servidor HTTP
    │   └── users.go             # Endpoints específicos de usuários
    ├── usecases/
    │   └── usecases.go          # Lógica de negócio
    ├── repositories/
    │   ├── repositories.go      # Interfaces dos repositórios
    │   └── users/
    │       └── users.go         # Implementação do repositório de usuários
    └── models/
        ├── users.go             # Modelos relacionados a usuários
        └── errors.go            # Modelos de resposta de erro
```

## 🛠️ Desenvolvimento

### Executar a aplicação

```bash
# 1. Iniciar o servidor da API
go run cmd/api/main.go
# Servidor disponível em http://localhost:3000

# 2. Em outro terminal, testar com o cliente
go run cmd/client/main.go
# Cliente faz requisição POST e testa a API
```

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

## 🔍 Troubleshooting

### **Problemas Comuns**

#### 🔌 **Erro de Conexão com PostgreSQL**
```bash
# Erro: "Failed to connect to database"
# Solução: Verificar se PostgreSQL está rodando
docker ps | grep postgres

# Verificar logs do container
docker logs <container_name>

# Testar conexão manual
docker exec -it <container_name> psql -U <user> -d <database> -c "SELECT 1;"
```

#### 🗄️ **Tabela não existe**
```bash
# Erro: "relation 'users' does not exist"
# Solução: Executar migration
docker exec -i <container_name> psql -U <user> -d <database> < migrations/001_create_users_table.sql
```

#### ⚙️ **Variáveis de ambiente**
```bash
# Verificar se .env está sendo carregado
cat .env

# Verificar variáveis no ambiente
env | grep DB_
```

### **Logs da Aplicação**

A aplicação usa `slog` para logging estruturado:

```bash
# Logs de conexão
INFO Connecting to database host=localhost port=5432 database=app_db
INFO Database connection established successfully

# Logs de operação
INFO User created successfully id=uuid email=user@email.com

# Logs de erro
ERROR Error inserting user error="duplicate key" user={...}
```

### **Verificar Saúde do Sistema**

```bash
# Verificar se API está respondendo
curl http://localhost:3000/users

# Verificar conexões do PostgreSQL
docker exec -it <container_name> psql -U <user> -d <database> -c "SELECT count(*) FROM pg_stat_activity;"

# Verificar tabelas criadas
docker exec -it <container_name> psql -U <user> -d <database> -c "\dt"
```

## 📦 Dependências

### **Core Dependencies**
- [github.com/google/uuid](https://github.com/google/uuid) - Geração de UUIDs
- [github.com/lib/pq](https://github.com/lib/pq) - Driver PostgreSQL para Go
- [github.com/joho/godotenv](https://github.com/joho/godotenv) - Carregamento de variáveis de ambiente

### **Banco de Dados**
- **PostgreSQL 15+** com suporte a:
  - UUIDs nativos (`gen_random_uuid()`)
  - Triggers automáticos (`updated_at`)
  - Índices otimizados
  - Connection pooling

## 🏛️ Padrões de Design Utilizados

### **Clean Architecture (Arquitetura Limpa)**
- **Principle**: Dependency Inversion - camadas externas dependem das internas
- **Benefit**: Independência de frameworks, banco de dados e UI
- **Structure**: Entities → Use Cases → Interface Adapters → Frameworks

### **Repository Pattern**
- **Purpose**: Abstração da camada de acesso aos dados
- **Benefit**: Facilita troca de implementação (in-memory ↔ PostgreSQL ↔ MongoDB)
- **Interface**: Define contrato independente da implementação

### **Dependency Injection**
- **Method**: Constructor injection através do padrão `New()`
- **Benefit**: Facilita testes com mocks e reduz acoplamento
- **Example**: `handlers.New(useCases)` injeta dependências

### **Use Case Pattern**
- **Encapsulation**: Cada caso de uso representa uma operação de negócio
- **Separation**: Isola regras de negócio de detalhes técnicos
- **Testability**: Permite testes unitários da lógica de negócio

### **SOLID Principles Aplicados**
- **S**ingle Responsibility: Cada camada tem uma responsabilidade
- **O**pen/Closed: Extensível via interfaces, fechado para modificação
- **L**iskov Substitution: Repositories são intercambiáveis
- **I**nterface Segregation: Interfaces específicas e coesas
- **D**ependency Inversion: Dependência de abstrações, não concreções

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

### CreateUserResponse
```go
type CreateUserResponse struct {
    NewUserID uuid.UUID `json:"newUserId"`
}
```

### ErrorResponse
```go
type ErrorResponse struct {
    Reason string `json:"reason"`
}
```

## 🔧 Configuração

A aplicação utiliza configuração baseada em variáveis de ambiente:

### **Arquivo .env (Desenvolvimento)**
```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=app_db
DB_USER=app_user
DB_PASSWORD=example
DB_SSLMODE=disable

# Application Configuration
APP_PORT=3000
APP_ENV=development

# Database Pool Configuration
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=300s
```

### **Variáveis de Ambiente (Produção)**
- `DB_HOST` - Host do PostgreSQL
- `DB_PORT` - Porta do PostgreSQL (padrão: 5432)
- `DB_NAME` - Nome do banco de dados
- `DB_USER` - Usuário do banco
- `DB_PASSWORD` - Senha do banco
- `DB_SSLMODE` - Modo SSL (disable/require/verify-full)
- `APP_PORT` - Porta da aplicação (padrão: 3000)
- `DB_MAX_OPEN_CONNS` - Máximo de conexões abertas (padrão: 25)
- `DB_MAX_IDLE_CONNS` - Máximo de conexões ociosas (padrão: 10)
- `DB_CONN_MAX_LIFETIME` - Tempo de vida das conexões (padrão: 5m)

## 🚧 Roadmap

### ✅ **Implementado**
- [x] Endpoints HTTP básicos (GET/POST /users)
- [x] **Integração completa com PostgreSQL**
- [x] **Connection pooling e gerenciamento de recursos**
- [x] **Configuração via variáveis de ambiente (.env)**
- [x] **Migrations automáticas do banco de dados**
- [x] **Graceful shutdown com cleanup**
- [x] Validação de entrada e tratamento de erros
- [x] Cliente HTTP de exemplo funcional
- [x] Respostas padronizadas (sucesso e erro)
- [x] Validação de email único (via banco)
- [x] Estrutura de projeto organizada (Go standards)
- [x] Logging estruturado com slog

### 🔄 **Em Desenvolvimento**
- [ ] Middleware de logging avançado para requests HTTP
- [ ] Validação mais robusta de dados de entrada (struct tags)
- [ ] Testes unitários abrangentes com mocks
- [ ] Health check endpoint para monitoramento

### 📋 **Planejado**
- [ ] Documentação OpenAPI/Swagger
- [ ] Autenticação e autorização (JWT)
- [ ] Paginação para listagem de usuários
- [ ] Containerização da aplicação (Dockerfile)
- [ ] CI/CD pipeline (GitHub Actions)
- [ ] Métricas e observabilidade (Prometheus)
- [ ] Rate limiting
- [ ] Backup automático do banco
- [ ] Ambiente de staging
- [ ] Cache distribuído

## 🤝 Contribuindo

1. Faça um fork do projeto
2. Crie sua feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para detalhes.

## 👥 Autores

- Kenard R Silva - [@kenardrs](https://github.com/kenardrs)

## 🙏 Agradecimentos

- Go community
- Clean Architecture principles by Robert C. Martin
- Todas as bibliotecas open source utilizadas
- Léo Miranda DEV: https://www.youtube.com/watch?v=EXwqzrcVXKg&t=1812s

---

**Nota**: Este é um projeto de estudo/exemplo. Para uso em produção, considere implementar recursos adicionais de segurança, logging, monitoramento e testes mais robustos.