package main

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

type StressTestResults struct {
	TotalRequests        int
	SuccessfulRequests   int
	HttpStatusCounts     map[int]int
	TotalTime            time.Duration
}

var (
	url           string
	totalRequests int
	concurrency   int
)

func performStressTest(url string, totalRequests, concurrency int) StressTestResults {
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := StressTestResults{
		HttpStatusCounts: make(map[int]int),
	}

	start := time.Now()
	semaphore := make(chan struct{}, concurrency)

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		semaphore <- struct{}{}

		go func() {
			defer wg.Done()
			defer func() { <-semaphore }()

			resp, err := http.Get(url)
			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				results.HttpStatusCounts[-1]++
				return
			}
			defer resp.Body.Close()

			results.TotalRequests++
			results.HttpStatusCounts[resp.StatusCode]++

			if resp.StatusCode == http.StatusOK {
				results.SuccessfulRequests++
			}
		}()
	}

	wg.Wait()
	results.TotalTime = time.Since(start)

	return results
}

func printReport(results StressTestResults) {
	seconds := int(results.TotalTime.Seconds())
	milliseconds := results.TotalTime.Milliseconds() % 1000

	fmt.Println("\n--- Load Test Report ---")
	fmt.Printf("Total Execution Time: %d.%03d seconds\n", seconds, milliseconds)
	fmt.Printf("Total Requests: %d\n", results.TotalRequests)
	fmt.Printf("Successful Requests (200 OK): %d\n", results.SuccessfulRequests)

	fmt.Println("\nHTTP Status Distribution:")

	var sortedStatuses []int
	for status := range results.HttpStatusCounts {
		sortedStatuses = append(sortedStatuses, status)
	}
	sort.Ints(sortedStatuses)

	for _, status := range sortedStatuses {
		count := results.HttpStatusCounts[status]

		switch status {
		case -1:
			fmt.Printf("  Connection/Request Errors: %d\n", count)
		case http.StatusOK:
			fmt.Printf("  Status %d (OK): %d requests\n", status, count)
		case http.StatusNotFound:
			fmt.Printf("  Status %d (Not Found): %d requests\n", status, count)
		case http.StatusInternalServerError:
			fmt.Printf("  Status %d (Internal Server Error): %d requests\n", status, count)
		case http.StatusBadRequest:
			fmt.Printf("  Status %d (Bad Request): %d requests\n", status, count)
		case http.StatusUnauthorized:
			fmt.Printf("  Status %d (Unauthorized): %d requests\n", status, count)
		case http.StatusForbidden:
			fmt.Printf("  Status %d (Forbidden): %d requests\n", status, count)
		default:
			fmt.Printf("  Status %d: %d requests\n", status, count)
		}
	}

	fmt.Println("\nStatus Percentages:")
	for _, status := range sortedStatuses {
		count := results.HttpStatusCounts[status]
		percentage := (float64(count) / float64(results.TotalRequests)) * 100

		switch status {
		case -1:
			fmt.Printf("  Connection/Request Errors: %.2f%%\n", percentage)
		case http.StatusOK:
			fmt.Printf("  Status %d (OK): %.2f%%\n", status, percentage)
		case http.StatusNotFound:
			fmt.Printf("  Status %d (Not Found): %.2f%%\n", status, percentage)
		case http.StatusInternalServerError:
			fmt.Printf("  Status %d (Internal Server Error): %.2f%%\n", status, percentage)
		case http.StatusBadRequest:
			fmt.Printf("  Status %d (Bad Request): %.2f%%\n", status, percentage)
		case http.StatusUnauthorized:
			fmt.Printf("  Status %d (Unauthorized): %.2f%%\n", status, percentage)
		case http.StatusForbidden:
			fmt.Printf("  Status %d (Forbidden): %.2f%%\n", status, percentage)
		default:
			fmt.Printf("  Status %d: %.2f%%\n", status, percentage)
		}
	}
}

func runStressTest(cmd *cobra.Command, args []string) {
	// Validação da URL
	if url == "" {
		fmt.Println("Erro: URL é um parâmetro obrigatório")
		os.Exit(1)
	}

	// Executa o teste de carga
	results := performStressTest(url, totalRequests, concurrency)

	// Imprime o relatório
	printReport(results)
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "stress-test",
		Short: "Ferramenta de teste de carga para serviços web",
		Long:  `Uma CLI para realizar testes de carga em serviços web`,
		Run:   runStressTest,
	}

	// Definição das flags
	rootCmd.Flags().StringVarP(&url, "url", "u", "", "URL do serviço a ser testado (obrigatório)")
	rootCmd.Flags().IntVarP(&totalRequests, "requests", "r", 100, "Número total de requests")
	rootCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 10, "Número de chamadas simultâneas")

	// Torna a flag URL obrigatória
	rootCmd.MarkFlagRequired("url")

	// Executa o comando root
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}