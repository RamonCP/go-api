# go-api — Lab de Estudo em Go

> Este repositório é um **laboratório de aprendizado** — um ambiente de experimentação para praticar desenvolvimento de APIs em Go com boas práticas de arquitetura, banco de dados relacional e containerização. Não é um produto em produção.

## Sobre o Projeto

API REST de gerenciamento de produtos construída em Go, com foco em demonstrar a aplicação de **Arquitetura Hexagonal (Ports & Adapters)** — um padrão arquitetural que isola a lógica de negócio das dependências externas (frameworks, banco de dados, etc.).

O domínio é simples de propósito: um CRUD de produtos, para que a atenção fique inteiramente na arquitetura e nas tecnologias, sem distrações de regras de negócio complexas.

## Tecnologias

| Tecnologia         | Papel no projeto                             |
| ------------------ | -------------------------------------------- |
| **Go 1.25**        | Linguagem principal                          |
| **Gin**            | Framework HTTP (adapter de entrada)          |
| **PostgreSQL 12**  | Banco de dados relacional (adapter de saída) |
| **lib/pq**         | Driver Go para PostgreSQL                    |
| **golang-migrate** | Migrations de banco de dados versionadas     |
| **Docker**         | Containerização da aplicação                 |
| **Docker Compose** | Orquestração local dos serviços              |
| **go-cmp**         | Comparação de structs em testes              |
| **GitHub Actions** | CI: lint, testes e build automatizados       |

## Arquitetura

O projeto aplica **Arquitetura Hexagonal**, organizando o código em três camadas bem definidas:

```
go-api/
├── cmd/go-api/          # Ponto de entrada — wiring de dependências (DI manual)
└── internal/
    ├── core/            # Núcleo isolado de qualquer framework
    │   ├── domain/      # Modelos de domínio (Product)
    │   ├── ports/       # Interfaces: driving (entrada) e driven (saída)
    │   └── services/    # Regras de negócio (use cases)
    └── adapters/
        ├── http/        # Adapter primário: handlers Gin + DTOs
        └── postgres/    # Adapter secundário: repositório SQL + migrations
```

**Fluxo de uma requisição:**

```
HTTP Request → Gin Handler → ProductService (interface) → ProductRepository (interface) → PostgreSQL
```

O `core` nunca importa `adapters`. Os adapters implementam interfaces definidas em `ports`, tornando toda a camada central testável sem banco de dados ou HTTP real.

## Endpoints

Base URL: `http://localhost:8000`

| Método   | Rota           | Descrição                     |
| -------- | -------------- | ----------------------------- |
| `GET`    | `/health`      | Readiness: API + banco (`200`/`503`) |
| `GET`    | `/products`    | Lista todos os produtos       |
| `GET`    | `/product/:id` | Busca produto por ID          |
| `POST`   | `/product`     | Cria um novo produto          |
| `PUT`    | `/product/:id` | Atualiza um produto existente |
| `DELETE` | `/product/:id` | Remove um produto             |

**Payload (JSON):**

```json
{
  "name": "Nintendo Switch",
  "price": 1999.9
}
```

## Como Rodar

### Pré-requisitos

- [Docker](https://www.docker.com/) e [Docker Compose](https://docs.docker.com/compose/) instalados

### 1. Subir com Docker Compose

```bash
docker compose up --build
```

Isso irá:

1. Subir o PostgreSQL na porta `5432`
2. Compilar e subir a API Go na porta `8000`
3. Executar as migrations automaticamente na inicialização
4. Seed com dados iniciais de exemplo

### 2. Testar os endpoints

```bash
# Listar produtos
curl http://localhost:8000/products

# Buscar por ID
curl http://localhost:8000/product/1

# Criar produto
curl -X POST http://localhost:8000/product \
  -H "Content-Type: application/json" \
  -d '{"name": "Produto Teste", "price": 49.90}'

# Atualizar produto
curl -X PUT http://localhost:8000/product/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Produto Atualizado", "price": 59.90}'

# Deletar produto
curl -X DELETE http://localhost:8000/product/1
```

### 3. Rodar localmente (sem Docker)

Requer Go 1.25+ e um PostgreSQL local rodando na porta 5432 com usuário `postgres` e senha `1234`.

```bash
# Instalar dependências
go mod download

# Subir apenas o banco
docker compose up go_db -d

# Rodar a aplicação
go run cmd/go-api/main.go
```

### 4. Derrubar o ambiente

```bash
docker compose down -v   # -v remove o volume do banco
```

## Testes

A estratégia de testes segue o padrão da comunidade Go (Google/Uber Style Guides): **biblioteca padrão `testing` + [`go-cmp`](https://github.com/google/go-cmp)**, sem frameworks de asserção (como `testify`). Isso mantém os testes desacoplados de APIs de terceiros e com mensagens de falha auto-explicativas.

### Tipos de teste e a arquitetura hexagonal

A pirâmide de testes mapeia diretamente nas camadas do projeto:

| Camada                | Tipo                  | Como testar                                                          |
| --------------------- | --------------------- | ------------------------------------------------------------------- |
| `core/services`       | Unitário              | Fakes das interfaces de `ports` — sem banco nem HTTP, em ms          |
| `core/domain`         | Unitário              | Direto, sobre valores puros                                         |
| `adapters/http`       | Unitário / integração | Mapeamento de DTOs; handlers com `httptest`                          |
| `adapters/postgres`   | Integração            | Banco real (ex.: testcontainers), atrás de uma build tag `integration` |

O maior benefício do hexagonal aparece aqui: como o `core` depende de **interfaces**, os testes unitários do `ProductService` rodam com um repositório _fake_ em memória — em milissegundos, sem subir o Postgres.

### Convenções

- **Table-driven tests** com subtests (`t.Run`): um caso por linha da tabela.
- Mensagens auto-diagnosticáveis: `Func(input) = got, want X` (sempre "got" antes de "want").
- `cmp.Diff` para comparar structs (aponta exatamente o campo que divergiu).
- _Test doubles_ escritos à mão (fakes), não mocks gerados.
- Black-box (`package xxx_test`) por padrão; white-box (`package xxx`) apenas para exercitar código não-exportado.

### Testes existentes

| Arquivo                                          | O que cobre                                                |
| ------------------------------------------------ | ---------------------------------------------------------- |
| `internal/core/services/product_service_test.go` | `GetProductById` com fake do repositório (black-box)       |
| `internal/adapters/http/product_dto_test.go`     | `toProductResponse` — mapeamento domínio → DTO (white-box) |

### Como rodar

```bash
# Todos os testes
go test ./...

# Verboso (mostra cada teste e subteste)
go test ./... -v

# Um teste específico (regex no nome)
go test ./... -run TestProductService_GetProductById

# Cobertura por pacote
go test ./... -cover

# Relatório visual de cobertura (abre no navegador)
go test ./... -coverprofile=cover.out && go tool cover -html=cover.out

# Detector de data races
go test ./... -race
```

> Os testes também rodam automaticamente no CI (GitHub Actions) a cada push e pull request.

## Próximos Passos

Este lab está em evolução. Os experimentos planejados para as próximas iterações são:

- [x] **Testes unitários** — cobrir o `core/services` com mocks das interfaces de repositório, demonstrando um dos maiores benefícios da arquitetura hexagonal
- [ ] **Testes de integração** — testar os adapters (HTTP handlers e repositório Postgres) com banco real
- [ ] **Variáveis de ambiente** — externalizar credenciais do banco (atualmente hardcoded) via `.env` ou flags de configuração
- [x] **Dockerfile multi-stage** — reduzir o tamanho da imagem final separando build e runtime
- [ ] **Middleware de logging** — adicionar logs estruturados com `slog` (stdlib) ou `zerolog`
- [x] **Health check endpoint** — `GET /health` para verificação de disponibilidade da API e do banco
- [ ] **Segundo adapter de saída** — implementar um repositório em memória para comparar com o Postgres sem mudar nada no `core`
- [ ] **OpenAPI/Swagger** — documentar os endpoints com `swaggo/swag`
- [x] **CI com GitHub Actions** — pipeline de build, lint e testes automatizados
