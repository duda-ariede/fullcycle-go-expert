package config

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func LoadConfig() {
	  // Verifica se há um caminho de env file definido
    envFile := os.Getenv("ENV_FILE")
    if envFile != "" {
        // Se ENV_FILE estiver definido, usa esse caminho específico
        log.Printf("Carregando configuração do arquivo: %s", envFile)
        viper.SetConfigFile(envFile)
    } else {
        // Configurações padrão de busca
        viper.SetConfigName(".env")
        viper.SetConfigType("env")
        viper.AddConfigPath("/app")  // Caminho dentro do container
        viper.AddConfigPath(".")
        viper.AddConfigPath("./configs")
    }

    // Definir valores padrão
    setDefaultConfigurations()

    // Tentar ler o arquivo de configuração
    if err := viper.ReadInConfig(); err != nil {
        log.Println("Aviso: Não foi possível ler o arquivo de configuração. Usando valores padrão.")
        log.Println("Erro detalhado:", err)
    } else {
        log.Printf("Configuração carregada do arquivo: %s", viper.ConfigFileUsed())
    }

    // Ativar monitoramento de mudanças no arquivo de configuração
    viper.WatchConfig()

    // Configurar callback para quando o arquivo mudar
    viper.OnConfigChange(func(e fsnotify.Event) {
        log.Println("Arquivo de configuração modificado:", e.Name)
				log.Println("Caminho completo:", e.Name)
    })
}

func setDefaultConfigurations() {
    viper.SetDefault("RATE_LIMIT_ENABLED", true)
    viper.SetDefault("RATE_LIMIT_IP_MAX_REQUESTS", 10)
    viper.SetDefault("RATE_LIMIT_IP_BLOCK_DURATION", 5)
    viper.SetDefault("RATE_LIMIT_TOKEN_MAX_REQUESTS", 100)
    viper.SetDefault("RATE_LIMIT_TOKEN_BLOCK_DURATION", 10)
    viper.SetDefault("REDIS_HOST", "localhost")
    viper.SetDefault("REDIS_PORT", "6379")
}

// Função para obter as configurações de rate limiter
func GetRateLimiterConfig() (int, int, int, int) {
    ipMaxRequests := viper.GetInt("RATE_LIMIT_IP_MAX_REQUESTS")
    ipBlockDuration := viper.GetInt("RATE_LIMIT_IP_BLOCK_DURATION")
    tokenMaxRequests := viper.GetInt("RATE_LIMIT_TOKEN_MAX_REQUESTS")
    tokenBlockDuration := viper.GetInt("RATE_LIMIT_TOKEN_BLOCK_DURATION")

    return ipMaxRequests, ipBlockDuration, tokenMaxRequests, tokenBlockDuration
}