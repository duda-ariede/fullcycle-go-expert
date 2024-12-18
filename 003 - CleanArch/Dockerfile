# Usar imagem oficial do Go
FROM golang:1.22-alpine

# Definir diretório de trabalho
WORKDIR /app

# Copiar go mod e sum files
COPY go.mod go.sum ./
COPY cmd/ordersystem/.env ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Compilar aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -o ordersystem ./cmd/ordersystem/main.go ./cmd/ordersystem/wire_gen.go

# Expor portas necessárias
EXPOSE 8000 50051 3000

# Dockerfile
RUN wget https://github.com/jwilder/dockerize/releases/download/v0.6.1/dockerize-linux-amd64-v0.6.1.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-v0.6.1.tar.gz \
    && rm dockerize-linux-amd64-v0.6.1.tar.gz

# No CMD
CMD dockerize -wait tcp://mysql:3306 -wait tcp://rabbitmq:5672 -timeout 60s ./ordersystem

# Comando para rodar a aplicação
#CMD ["./ordersystem"]