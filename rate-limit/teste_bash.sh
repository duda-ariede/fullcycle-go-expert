#!/bin/bash

if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
else
    echo "ERRO: Arquivo .env não encontrado na pasta configs."
    exit 1
fi

# Usar valores do .env ou definir padrão se não existirem
MAX_IP_REQUESTS=${RATE_LIMIT_IP_MAX_REQUESTS}
MAX_TOKEN_REQUESTS=${RATE_LIMIT_TOKEN_MAX_REQUESTS}

echo "Testando Rate Limiter por IP:"
echo "Limite configurado: $MAX_IP_REQUESTS requisições"

ip_blocked=0
for i in $(seq 1 $((MAX_IP_REQUESTS + 1))); do
    response=$(curl -s -o /dev/null -w "%{http_code}" http://app:8080)
    echo "Requisição IP $i - Código HTTP: $response"

    if [ $response -eq 429 ]; then
        ip_blocked=1
        break
    fi
done

if [ $ip_blocked -eq 0 ]; then
    echo "ERRO: Limite de IP não foi aplicado!"
fi

echo -e "\nTestando Rate Limiter por Token:"
echo "Limite configurado: $MAX_TOKEN_REQUESTS requisições"

token_blocked=0
for i in $(seq 1 $((MAX_TOKEN_REQUESTS + 1))); do
    response=$(curl -s -o /dev/null -w "%{http_code}" -H "API_KEY: test_token" http://app:8080)
    echo "Requisição Token $i - Código HTTP: $response"

    if [ $response -eq 429 ]; then
        token_blocked=1
        break
    fi
done

if [ $token_blocked -eq 0 ]; then
    echo "ERRO: Limite de Token não foi aplicado!"
fi

# Verificação de logs
#echo -e "\nVerificando logs para detalhes:"
#docker compose logs app | grep -E "Rate limit|Increment|Block"
