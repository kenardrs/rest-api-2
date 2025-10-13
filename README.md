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
- ✅ Validação de dados de entrada e tratamento de erros
- ✅ Cliente HTTP funcional com tratamento de erros
- ✅ Respostas padronizadas de sucesso e erro

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

## 📦 Dependências

- [github.com/google/uuid](https://github.com/google/uuid) - Geração de UUIDs

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

Atualmente, a aplicação utiliza configurações hardcoded. Futuras versões incluirão:

- Variáveis de ambiente
- Arquivo de configuração
- Configuração de banco de dados
- Configuração de logging

## 🚧 Roadmap

### ✅ **Implementado**
- [x] Endpoints HTTP básicos (GET/POST /users)
- [x] Validação de entrada e tratamento de erros
- [x] Cliente HTTP de exemplo funcional
- [x] Respostas padronizadas (sucesso e erro)
- [x] Validação de email único
- [x] Estrutura de projeto organizada

### 🔄 **Em Desenvolvimento**
- [ ] Middleware de logging avançado
- [ ] Validação mais robusta de dados de entrada
- [ ] Testes unitários abrangentes

### 📋 **Planejado**
- [ ] Integração com banco de dados (PostgreSQL/MySQL)
- [ ] Documentação OpenAPI/Swagger
- [ ] Autenticação e autorização
- [ ] Paginação para listagem
- [ ] Dockerização
- [ ] CI/CD pipeline
- [ ] Métricas e observabilidade
- [ ] Rate limiting
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