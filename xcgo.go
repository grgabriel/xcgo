package main

import (
	"encoding/json"
	"fmt"
	"bufio"
	"os"
	"net/http"
	"strings"
	"flag"
	"time"
)

// API key file must exist and contain only the api key string
const apiKeyFile string = "/home/corian/.config/exchange/api.key"
// the cache will be created if it does not exist
const jsonCache string = "/home/corian/.config/exchange/exchange.json"

var exchangeData = JsonResponse {}

// Structure to store the JSON response from the api
type JsonResponse struct {
	Metadata struct {
		LastUpdate string `json:"last_updated_at"`
	} `json:"meta"`
	Coindata struct {
		EUR struct {
			Code string `json:"code"`
			Value float64 `json:"value"`
		} `json:"EUR"`
	} `json:"data"`
}

// Checks for error, kills program if error encountered
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Load the API key from file and return as string
func loadApi() string {
	key, err := os.ReadFile(apiKeyFile)
	if err != nil {
		fmt.Println("Cannot load API key")
		panic(err)
	}
	return strings.ReplaceAll(string(key), "\n", "")
}

// Get current exchange data from API. Returns a JSON string
func callApi() {
 	apiKey := loadApi()
	api := "https://api.currencyapi.com/v3/latest?apikey=" + apiKey + "&currencies=EUR&base_currency=GBP"
	
	resp, err := http.Get(api)
	check(err)
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	json := ""
	for scanner.Scan() {
		json += scanner.Text()
	}

	setExchange(json)

}


// Populate exchangeData with json data from jsonData string
func setExchange(jsonData string) {
	err := json.Unmarshal([]byte(jsonData), &exchangeData)
	if err != nil {
		fmt.Println("Error parsing JSON data!")
		panic(err)
	}
}

// Populate the exchange data
// If the cache file does not exist, call the API for data and make a new one
// If the cache file is older than one day, call API and make a new one
func getExchange(cur string, val float64) {

	// try to open cache file
	f, err := os.ReadFile(jsonCache)
	
	// file does not exist, call api and make new file
	if err != nil {
		fmt.Println("Cache file not found. Loading API")
		callApi()
	}else{
		setExchange(string(f))	
		
		// Transform "2024-05-15T16:04:30Z" into "2024-05-15 16:04:30"
		timeStr := strings.ReplaceAll(strings.ReplaceAll(exchangeData.Metadata.LastUpdate, "Z", ""), "T", " ")
		
		// Check if the cache update is latest
		// Cache is updated every day at 23:59:59 so the latest is always yesterday
		if !isYesterday(timeStr) {
			fmt.Println("Cache out of date. Updating")
			callApi()
		}
	}
}

// Check if givenDate is yesterday
func isYesterday(givenDate string) bool {
	t, err := time.Parse("2006-01-02 15:04:05", givenDate)
	if err != nil {
		fmt.Println("Error!")
		return false
	}
	yesterday := time.Now().AddDate(0,0,-1)
	return yesterday.Format("2006-02-01") == t.Format("2006-02-01")
}


func main() {

	// parse flags
	// usage is xcgo -c e|p -v 123.45

	cur := flag.String("c", "p", "Currency, p or e")
	val := flag.Float64("v", 1.0, "Value to convert")
	flag.Parse()
	
	calc := 0.0

	if *cur != "p" && *cur != "e" {
		fmt.Println("Invalid argument.")
		fmt.Println("Usage: xce -c e|p -v 123.45")
		fmt.Println("	e: convert Euro to Pound\n	p: convert Pound to Euro")
	}else {
		getExchange(*cur, *val)
		switch *cur {
			case "p":
				calc = exchangeData.Coindata.EUR.Value * *val
				fmt.Println(*val, "£ is ", calc, "€")
			case "e":
				calc = *val / exchangeData.Coindata.EUR.Value
				fmt.Println(*val, "€ is ", calc, "£")
			default:
				panic("how did you get here")
		}
		
	}
}

