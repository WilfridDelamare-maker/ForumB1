FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -ldflags="-extldflags=-static" -o forum .


FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/forum ./forum
COPY templates/ ./templates/
COPY static/ ./static/

# Crée les dossiers nécessaires au démarrage
# static/upload/ : dossier des images uploadées (non copié car exclu du .dockerignore)
# data/          : dossier de la BDD SQLite (monté par le container db via volume)
RUN mkdir -p /app/static/upload /app/data

EXPOSE 8085

CMD ["./forum"]
