package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func init() {
	// Carrega variáveis de ambiente do arquivo .env
	if err := godotenv.Load(); err != nil {
		log.Println("Nenhum arquivo .env encontrado. Usando variáveis de ambiente do sistema.")
	}
}

// normalizeCityName remove acentuação, substitui espaços e converte para minúsculas
func normalizeCityName(cityName string) string {
	// Remove acentuação
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, cityName)

	// Converte para minúsculas
	result = strings.ToLower(result)

	// Remove caracteres especiais e substitui espaços por hífens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	result = reg.ReplaceAllString(result, "-")

	// Remove hífens no início e no fim
	result = strings.Trim(result, "-")

	return result
}

// ViaCEPResponse represents the structure of the ViaCEP API response
type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
	Erro       bool   `json:"erro,omitempty"`
}

// WeatherAPIResponse represents the structure of the WeatherAPI response
type WeatherAPIResponse struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

// TemperatureResponse represents the response with temperatures in different scales
type TemperatureResponse struct {
	City   string  `json:"city"`
	TempC  float64 `json:"temp_C"`
	TempF  float64 `json:"temp_F"`
	TempK  float64 `json:"temp_K"`
}

// convertTemperatures converts Celsius to Fahrenheit and Kelvin
func convertTemperatures(celsius float64) TemperatureResponse {
	return TemperatureResponse{
		TempC: roundToOneDecimal(celsius),
		TempF: roundToOneDecimal(celsius*1.8 + 32),
		TempK: roundToOneDecimal(celsius + 273),
	}
}

// roundToOneDecimal rounds a float to one decimal place
func roundToOneDecimal(num float64) float64 {
	return float64(int(num*10)) / 10
}

// validateCEP checks if the CEP is valid (8 digits)
func validateCEP(cep string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, cep)
	return match
}

// getCityFromCEP retrieves city information from ViaCEP
func getCityFromCEP(cep string) (string, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var viaCEPResp ViaCEPResponse
	err = json.NewDecoder(resp.Body).Decode(&viaCEPResp)
	if err != nil {
		return "", err
	}

	if viaCEPResp.Erro || viaCEPResp.Localidade == "" {
		return "", fmt.Errorf("CEP not found")
	}

	return viaCEPResp.Localidade, nil
}

// getWeatherForCity retrieves weather information for a given city
func getWeatherForCity(city string) (WeatherAPIResponse, error) {
	apiKey := os.Getenv("WEATHERAPI_KEY")
	if apiKey == "" {
		return WeatherAPIResponse{}, fmt.Errorf("WeatherAPI key not set")
	}

	normalizedCity := normalizeCityName(city)

	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, normalizedCity)

	resp, err := http.Get(url)
	if err != nil {
		return WeatherAPIResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return WeatherAPIResponse{}, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	var weatherResp WeatherAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&weatherResp)
	if err != nil {
		return WeatherAPIResponse{}, err
	}

	return weatherResp, nil
}

// weatherHandler handles the main temperature lookup logic
func weatherHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cep := vars["cep"]

	// Validate CEP format
	if !validateCEP(cep) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid zipcode"))
		return
	}

	// Get city from CEP
	city, err := getCityFromCEP(cep)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("can not find zipcode"))
		return
	}

	// Get weather for city
	weatherResp, err := getWeatherForCity(city)
	if err != nil {
		log.Printf("Weather API error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error retrieving weather"))
		return
	}

	// Convert temperatures
	tempResponse := convertTemperatures(weatherResp.Current.TempC)
	tempResponse.City = weatherResp.Location.Name

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tempResponse)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/weather/{cep}", weatherHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %v", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}