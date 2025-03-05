# CEP Weather Lookup Application

## Descrição do Projeto

Esta aplicação é um serviço web desenvolvido em Go que permite consultar informações meteorológicas para uma localidade brasileira usando um Código de Endereçamento Postal (CEP). O sistema integra duas APIs externas para fornecer informações precisas de temperatura.

## Funcionalidades

- Validação de CEP (8 dígitos)
- Consulta de localidade através do CEP
- Recuperação de temperatura atual
- Conversão de temperatura para múltiplas escalas:
  - Celsius
  - Fahrenheit
  - Kelvin

## Requisitos Técnicos

### Pré-requisitos

- Go 1.23 ou superior
- Docker (opcional, para containerização)
- Conta na WeatherAPI (para obtenção de chave de API)

### Dependências

- Gorilla Mux (Roteamento HTTP)
- GoDotEnv (Gerenciamento de variáveis de ambiente)

## Instalação

### Configuração Local

1. Clone o repositório:

   ```bash
   git clone https://github.com/duda-ariede/fullcycle-go-expert/cep-weather.git
   cd cep-weather
   ```

2. Instale as dependências:

   ```bash
   go mod download
   ```

3. Configure o arquivo `.env`:

   ```bash
   cp .env.example .env
   ```

   Edite o `.env` e adicione sua chave da WeatherAPI

4. Execute a aplicação:
   ```bash
   go run main.go
   ```

### Execução com Docker

1. Construa a imagem:

   ```bash
   docker-compose build
   ```

2. Inicie o container:
   ```bash
   docker-compose up
   ```

## Endpoints

### Consulta de Temperatura por CEP

- **URL**: `/weather/{cep}`
- **Método**: GET
- **Parâmetros**: CEP (8 dígitos)

### Códigos de Resposta

- `200 OK`: Consulta bem-sucedida
- `422 Unprocessable Entity`: CEP inválido
- `404 Not Found`: CEP não encontrado
- `500 Internal Server Error`: Erro na consulta de temperatura

### Exemplo de Resposta

```json
{
	"city": "São Paulo",
	"temp_C": 22.5,
	"temp_F": 72.5,
	"temp_K": 295.5
}
```

## Testes

Execute os testes com:

```bash
go test
```

### Casos de Teste Implementados

- Validação de CEP válido
- Validação de formato de CEP inválido
- Consulta de CEP inexistente

## Configuração de Ambiente

### Variáveis de Ambiente

- `WEATHERAPI_KEY`: Chave de API da WeatherAPI
- `PORT`: Porta para execução do serviço (padrão: 8080)

## Deployment

### Google Cloud Run

1. Construa a imagem Docker
2. Faça o push para o Google Container Registry
3. Crie um serviço no Cloud Run
4. Configure as variáveis de ambiente

## Considerações de Segurança

- Nunca commite chaves de API no repositório
- Use sempre `.env.example` como template
- Em produção, use variáveis de ambiente do provedor de nuvem

## Limitações Conhecidas

- Dependência de APIs externas
- Possíveis variações de temperatura
- Cobertura limitada de CEPs

## Contribuição

1. Faça um fork do projeto
2. Crie uma branch para sua feature
3. Commit suas alterações
4. Abra um Pull Request

# Acesso no Google Cloud Run

https://cep-weather-272445865940.us-central1.run.app/weather/CEP
