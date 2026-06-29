# Roadmap do MVP — go-api → "Goodreads"

Este documento descreve o roadmap para evoluir esta API (Go + arquitetura
hexagonal + Postgres + chi) de um CRUD genérico de `Product` para um MVP
inspirado no Goodreads.

## Visão e loop de valor

O loop central que o MVP precisa entregar:

> **descubro um livro → marco que quero ler / li → dou nota → isso ajuda outros a descobrir**

Toda feature do MVP serve esse loop. O que não serve fica para depois.

## Princípios

- **Fatias verticais**: cada fase é entregável e testável de ponta a ponta
  (domain → ports → services → adapters), ~1 PR coeso.
- **Arquitetura hexagonal preservada**: o domínio não conhece HTTP nem SQL.
- **Refactor de domínio, não de infra**: mantemos o estilo atual (`database/sql`
  puro, DTOs de resposta, testes black-box em services e white-box em DTOs).
- **Migrations forward-only**: não reescrevemos migrations já aplicadas.

## Escopo do MVP

### Entra

| # | Feature | Por quê é essencial |
|---|---------|---------------------|
| 1 | Catálogo de livros (`Book` + `Author`) | Sem acervo não há o que fazer. |
| 2 | Autenticação de usuário | "Minha estante" exige identidade. |
| 3 | Estante com status (Quero ler / Lendo / Lido) | A feature que define o Goodreads. |
| 4 | Rating 1–5 + média por livro | Fecha o loop de valor e alimenta descoberta. |
| 5 | Busca por título/autor | Sem busca, o catálogo é inútil na prática. |

### Fica de fora (pós-MVP)

Resenhas longas, seguir usuários, feed social, reading challenge, listas
(Listopia), recomendações, estantes customizadas. Todas valiosas, nenhuma
essencial pro loop. Recomendação só faz sentido depois de existir massa de
ratings.

> **Quase-MVP**: resenha em texto é só um campo `text` opcional no rating —
> custo quase zero, pode entrar junto da Fase 3 se desejado.

---

## Fase 0 — Refactor de base: `Product` → `Book` + `Author`

**Objetivo:** o CRUD que já existe passa a falar de livros e autores.

**Decisão de modelagem:** Autor como entidade **1:N** (um livro tem um autor;
um autor tem vários livros). Abre caminho para "página do autor" no futuro sem
o custo de uma relação N:N agora.

### Modelo de dados

```
author(id SERIAL PK, name VARCHAR NOT NULL)
book(id SERIAL PK, title VARCHAR NOT NULL, author_id INT FK→author(id),
     isbn VARCHAR, description TEXT, published_year INT, page_count INT)
```

`Book` carrega `AuthorID` na escrita; nas leituras o repositório faz `JOIN author`
e popula o autor aninhado na resposta JSON.

### Mudanças por camada

- **Migrations** (novas):
  - `000003_create_authors_table` — cria `author`
  - `000004_create_books_table` — `DROP TABLE product` + cria `book` com FK
  - `000005_seed_books` — autores + livros de exemplo
- **Domain** (`internal/core/domain/`): `product.go` → `book.go` (`Book` com
  `Author` embutido para leitura); novo `author.go` (`Author{ID, Name}`).
- **Ports**: `ProductService`/`ProductRepository` → `BookService`/`BookRepository`;
  novos `AuthorService`/`AuthorRepository` (Create, GetByID, GetAll).
- **Services**: `book_service.go` + `author_service.go`; adapta o teste para
  `book_service_test.go`.
- **Adapters Postgres**: `book_repository.go` (queries com JOIN) + `author_repository.go`.
- **Adapters HTTP**: `book_handler.go` + `author_handler.go`; `book_dto.go`
  (`BookResponse` com autor aninhado) + DTO de autor; `book_dto_test.go`.
- **Rotas** (`cmd/go-api/main.go`):
  - `/books`, `/book/{id}` (GET/POST/PUT/DELETE)
  - `/authors`, `/author/{id}` (GET), `POST /author`

### Decisões de escopo

- Author: Create + Get + List (sem update/delete ainda).
- Criar livro valida que `author_id` existe.
- `health` fica intacto.

### Critérios de pronto

- `go build ./...` e `go test ./...` verdes.
- Smoke test: criar autor → criar livro → listar livros com autor aninhado.

---

## Fase 1 — Usuários & Autenticação

**Objetivo:** ter identidade e requisições autenticadas — pré-requisito de tudo
que é "do usuário".

### Modelo de dados

```
app_user(id SERIAL PK, email VARCHAR UNIQUE NOT NULL,
         password_hash VARCHAR NOT NULL, display_name VARCHAR,
         created_at TIMESTAMPTZ DEFAULT now())
```

### Mudanças por camada

- **Migration** `000006_create_users_table`.
- **Domain**: `user.go` (`User`); evitar expor `PasswordHash` em respostas.
- **Ports**: `UserService` (Register, Authenticate, GetByID),
  `UserRepository` (Create, GetByEmail, GetByID).
- **Services**: `user_service.go` — hash de senha com `bcrypt`, geração de JWT.
- **Adapters HTTP**:
  - `POST /auth/register`, `POST /auth/login` (retorna JWT)
  - Middleware de autenticação chi que valida o JWT e injeta o `userID` no
    `context.Context`.
- **Config**: segredo do JWT via env (seguir o padrão de `internal/config`).

### Decisões de escopo

- JWT stateless (sem tabela de sessão) para simplicidade.
- Validação de email/senha mínima; sem verificação de email no MVP.

### Critérios de pronto

- Registrar → logar → acessar rota protegida com o token funciona.
- Senha nunca trafega/retorna em texto; nunca logada.
- Testes de service para hashing e fluxo de autenticação.

---

## Fase 2 — Estantes (o coração)

**Objetivo:** o usuário marca o status de leitura dos livros. A partir daqui o
produto é reconhecível como um Goodreads.

### Modelo de dados

```
user_book(
  user_id INT FK→app_user(id),
  book_id INT FK→book(id),
  status  VARCHAR NOT NULL CHECK (status IN ('want_to_read','reading','read')),
  updated_at TIMESTAMPTZ DEFAULT now(),
  PRIMARY KEY (user_id, book_id)
)
```

### Mudanças por camada

- **Migration** `000007_create_user_book_table`.
- **Domain**: `shelf.go` — `ShelfEntry{UserID, BookID, Status}` e um tipo
  `ReadingStatus` com os valores válidos.
- **Ports**: `ShelfService` (SetStatus, RemoveFromShelf, ListByUser[, filtro por
  status]), `ShelfRepository`.
- **Services**: `shelf_service.go` — valida transições/valores de status.
- **Adapters HTTP** (todas protegidas por auth, usam o `userID` do contexto):
  - `PUT /me/books/{bookId}` — define/atualiza status (upsert)
  - `DELETE /me/books/{bookId}` — remove da estante
  - `GET /me/books?status=reading` — minha estante, com filtro opcional
- O `userID` vem **sempre** do token, nunca do corpo da requisição.

### Decisões de escopo

- Apenas as 3 estantes padrão (sem estantes customizadas).
- Progresso de leitura (página atual / %) fica para pós-MVP.

### Critérios de pronto

- Adicionar livro à estante, mudar status e listar por status funcionam.
- Um usuário não consegue ler/alterar a estante de outro.

---

## Fase 3 — Ratings

**Objetivo:** fechar o loop de valor — nota do usuário + média por livro.

### Modelo de dados

```
rating(
  user_id INT FK→app_user(id),
  book_id INT FK→book(id),
  score   SMALLINT NOT NULL CHECK (score BETWEEN 1 AND 5),
  review_text TEXT,                       -- "quase-MVP": opcional
  created_at TIMESTAMPTZ DEFAULT now(),
  PRIMARY KEY (user_id, book_id)
)
```

### Mudanças por camada

- **Migration** `000008_create_rating_table`.
- **Domain**: `rating.go` — `Rating{UserID, BookID, Score, ReviewText}`.
- **Ports**: `RatingService` (Rate (upsert), GetUserRating, ListByBook),
  `RatingRepository` com agregação (`AVG(score)`, `COUNT(*)`).
- **Services**: `rating_service.go` — valida `score` 1–5.
- **Adapters HTTP**:
  - `PUT /books/{bookId}/rating` (protegida) — cria/atualiza minha nota
  - `GET /books/{bookId}/ratings` — lista notas/resenhas do livro
  - `BookResponse` passa a incluir `average_rating` e `ratings_count`.

### Decisões de escopo

- **Média agregada via query** (`AVG`/`COUNT`), nunca calculada em memória sobre
  todas as linhas.
- Um rating por usuário por livro (upsert pela PK composta).

### Critérios de pronto

- Avaliar livro, atualizar a nota e ver a média/contagem atualizadas.
- `score` fora de 1–5 é rejeitado.

---

## Fase 4 — Busca & descoberta

**Objetivo:** descoberta funciona; MVP completo.

### Mudanças por camada

- **Migration** `000009_add_search_indexes` — índices para busca (começar com
  índices em `book.title` / `author.name`; evoluir para full-text do Postgres
  (`tsvector` + GIN) se necessário).
- **Ports/Services**: estender `BookService` com `Search(query, filtros)`.
- **Adapters Postgres**: query de busca (`ILIKE` no MVP; caminho de evolução
  para `to_tsvector`/`plainto_tsquery`).
- **Adapters HTTP**:
  - `GET /books?q=...&genre=...` — busca por título/autor + filtro por gênero
  - Paginação (`limit`/`offset` ou cursor).

### Decisões de escopo

- Começar simples (`ILIKE`), medir, e só então migrar para full-text.
- Gênero pode entrar como campo simples em `book` ou tabela `genre` dedicada —
  decidir no início da fase conforme a necessidade de filtros.

### Critérios de pronto

- Buscar por trecho de título e por nome de autor retorna resultados relevantes.
- Resultados paginados.

---

## Quadro-resumo das fases

| Fase | Entrega | Domínios novos | Risco |
|------|---------|----------------|-------|
| 0 | `Product` → `Book` + `Author` | Book, Author | Baixo |
| 1 | Usuários & Auth (JWT) | User | Médio |
| 2 | Estantes (status de leitura) | ShelfEntry | Médio |
| 3 | Ratings + média por livro | Rating | Médio |
| 4 | Busca & filtros | — (estende Book) | Médio |

## Temas transversais (aplicar ao longo das fases)

- **Tratamento de erro**: hoje muitos handlers retornam `500` genérico. Conforme
  o domínio cresce, mapear erros de domínio para status corretos
  (`404` não encontrado, `400` validação, `401/403` auth, `409` conflito).
- **Migrations**: sempre par up/down; forward-only.
- **Testes**: manter o padrão atual (services black-box com fakes manuais; DTOs
  white-box). Cobrir validações novas de cada fase.
- **Config**: novos segredos/parametros (JWT, etc.) via `internal/config` + env.
- **Observabilidade**: já há `middleware.Logger`/`Recoverer`; considerar
  request ID e logs estruturados quando houver tráfego autenticado.

## Pós-MVP (backlog priorizado)

1. Resenhas ricas + curtidas/comentários em resenhas
2. Seguir usuários + feed de atividades
3. Reading Challenge (meta anual)
4. Listas (Listopia) e estantes customizadas
5. Recomendações (baseadas em estantes/ratings)
6. Página do autor (já habilitada pela modelagem 1:N da Fase 0)
7. Progresso de leitura (página atual / %)
