# Serviço de Rate Limiter

Um serviço de limitação de taxa (rate limiting) implementado em Go, utilizando Redis para armazenamento. O serviço suporta limitação tanto por endereço IP quanto por tokens de API, com limites e durações de bloqueio configuráveis.

## Funcionalidades

- Duas estratégias de limitação de taxa:
  - Limitação baseada em IP
  - Limitação baseada em token
- Limites de requisições e durações de bloqueio configuráveis
- Armazenamento em Redis para ambientes distribuídos
- Atualizações de configuração em tempo real via variáveis de ambiente
- Suporte a Docker para fácil implantação

## Arquitetura

O rate limiter segue um padrão de arquitetura limpa com os seguintes componentes:

- **Middleware**: Gerencia requisições HTTP e extrai informações do cliente
- **Service**: Implementa a lógica principal de limitação de taxa
- **Repository**: Gerencia o armazenamento de dados no Redis
- **Configuration**: Lida com o gerenciamento dinâmico de configurações

## Configuração

O serviço pode ser configurado usando variáveis de ambiente em um arquivo `.env`:

```env
# Habilitar/desabilitar limitação de taxa
RATE_LIMIT_ENABLED=true

# Limitação baseada em IP
RATE_LIMIT_IP_MAX_REQUESTS=10       # Máximo de requisições permitidas por minuto
RATE_LIMIT_IP_BLOCK_DURATION=5      # Duração do bloqueio em minutos

# Limitação baseada em token
RATE_LIMIT_TOKEN_MAX_REQUESTS=100   # Máximo de requisições permitidas por minuto
RATE_LIMIT_TOKEN_BLOCK_DURATION=10  # Duração do bloqueio em minutos

# Configuração do Redis
REDIS_HOST=redis
REDIS_PORT=6379
```

## Executando com Docker

1. Clone o repositório
2. Crie um arquivo `.env` com sua configuração
3. Execute os serviços usando Docker Compose:

```bash
docker compose up -d
```

## Testes

O projeto inclui um container de teste que verifica automaticamente a funcionalidade de limitação de taxa. Para executar os testes:

```bash
docker compose up tester
```

O testador irá:

1. Testar a limitação baseada em IP
2. Testar a limitação baseada em token
3. Verificar o comportamento correto de bloqueio
4. Registrar resultados no console

## Como Funciona

### Processo de Limitação de Taxa

1. **Processamento de Requisição**:

   - O middleware extrai o IP do cliente ou token de API
   - A extração de IP segue uma ordem de prioridade: X-Forwarded-For > X-Real-IP > Remote-Addr

2. **Verificação de Limites**:

   - O serviço verifica se o cliente está atualmente bloqueado
   - Se não estiver bloqueado, incrementa o contador de requisições
   - Contadores expiram após 1 minuto (janela deslizante)

3. **Mecanismo de Bloqueio**:
   - Quando os limites são excedidos, o cliente é bloqueado
   - A duração do bloqueio é configurada separadamente para limites de IP e token
   - O status de bloqueio é armazenado no Redis com expiração automática

### Códigos de Resposta

- `200`: Requisição bem-sucedida
- `429`: Muitas Requisições (limite excedido)

## Uso da API

Para usar a limitação baseada em token, inclua a chave de API no cabeçalho da requisição:

```bash
curl -H "API_KEY: seu_token" http://localhost:8080
```

Para limitação baseada em IP, simplesmente faça requisições sem uma chave de API:

```bash
curl http://localhost:8080
```

## Configuração Dinâmica

O serviço suporta atualizações de configuração em tempo real:

1. Modifique o arquivo `.env`
2. O serviço detecta automaticamente as mudanças
3. Novas configurações de limite são aplicadas imediatamente
4. Bloqueios existentes permanecem até sua duração expirar

## Monitoramento e Depuração

O serviço fornece logs detalhados para monitorar o comportamento da limitação de taxa:

- Rastreamento de requisições
- Eventos de limite excedido
- Mudanças de configuração
- Detecção de IP do cliente

Para visualizar os logs:

```bash
docker compose logs app
```

## Tratamento de Erros

O serviço inclui tratamento robusto de erros para:

- Problemas de conexão com Redis
- Configurações inválidas
- Falhas na extração de IP
- Erros de gerenciamento de contador
