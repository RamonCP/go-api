# syntax=docker/dockerfile:1

# ──────────────────────────────────────────────────────────────
# Stage 1 — build: compila o binário usando o toolchain completo do Go.
# Esta camada é pesada (SDK + dependências), mas NÃO vai para a imagem final.
# ──────────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copia só os manifests primeiro para aproveitar o cache de camadas do Docker:
# enquanto go.mod/go.sum não mudarem, o `go mod download` não roda de novo.
COPY go.mod go.sum ./
RUN go mod download

# Agora o restante do código-fonte (inclui as migrations, embutidas via go:embed).
COPY . .

# Binário estático e enxuto:
#   CGO_ENABLED=0 → sem dependência da libc (drivers são Go puro) → roda em base mínima
#   -ldflags "-s -w" → remove tabela de símbolos e info de debug (binário menor)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/main ./cmd/go-api

# ──────────────────────────────────────────────────────────────
# Stage 2 — runtime: imagem final só com o binário e o mínimo para rodar.
# ──────────────────────────────────────────────────────────────
FROM alpine:3.21

# ca-certificates: confiança em TLS caso a app passe a fazer chamadas HTTPS.
# Usuário não-root: boa prática de segurança em containers.
RUN apk add --no-cache ca-certificates \
    && adduser -D -u 10001 appuser

WORKDIR /app

# Traz APENAS o binário compilado do stage anterior.
COPY --from=builder /app/main .

USER appuser

EXPOSE 8000

# Health check em nível de container: o Docker chama GET /health periodicamente
# e marca o container como healthy/unhealthy. O busybox wget já vem no alpine.
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://localhost:8000/health || exit 1

ENTRYPOINT ["./main"]
