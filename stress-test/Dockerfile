# Imagem base com Go
FROM golang:1.22-alpine AS builder

# Define diretório de trabalho
WORKDIR /app

# Copia os arquivos de dependência
COPY go.mod go.sum ./

# Baixa as dependências
RUN go mod download

# Copia o código fonte
COPY . .

# Compila o aplicativo
RUN CGO_ENABLED=0 GOOS=linux go build -o stress-test

# Imagem final menor
FROM alpine:latest

# Instala certificados CA para fazer requisições HTTPS
RUN apk --no-cache add ca-certificates

# Copia o binário compilado
COPY --from=builder /app/stress-test /usr/local/bin/stress-test

# Define o entrypoint
ENTRYPOINT ["stress-test"]