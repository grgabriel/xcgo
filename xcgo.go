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

// Load the API key from file
func loadApi() string {
	key, err := os.ReadFile(apiKeyFile)
	if err != nil {
		panic("Cannot load API key")
	}
	keyText := string(key)
	keyText = strings.ReplaceAll(keyText, "\n", "")
	return keyText 
}

// Get current exchange data from API. Returns a JSON string
func callApi() string {
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
	return json
}

// Calculate the exchange. Prints calculated value
// If the cache file does not exist, call the API for data and make a new one
// If the cache file is older than one day, call API and make a new one
func getExchange(cur string, val float64) {
	// try to open cache file
	f, err := os.ReadFile(jsonCache)
	str := ""
	// file does not exist, call api and make new file
	if err != nil {
		fmt.Println("Cache file not found. Loading API")
		str = callApi()
		ferr := os.WriteFile(jsonCache, []byte(str), 0644)
		if ferr != nil {
			panic("Could not write new json file!")
		}
	}else{
		str = string(f)
		resp := JsonResponse {}
		err := json.Unmarshal([]byte(str), &resp)
		if err != nil {
			fmt.Println("Error parsing JSON!")
			panic(err)
		}
		timeStr := resp.Metadata.LastUpdate
		timeStr = strings.ReplaceAll(timeStr, "T", " ")
		timeStr = strings.ReplaceAll(timeStr, "Z", "")
		if !isYesterday(timeStr) {
			fmt.Println("Cache out of date. Updating")
			str = callApi()
			ferr := os.WriteFile(jsonCache, []byte(str), 0644)
			if ferr != nil {
				panic(ferr)
			}
		}
		
		


		if cur == "p" {
			e := val * resp.Coindata.EUR.Value
			fmt.Println("Pound conversion: ", val, "£ is ", e, "€")
		} else if cur == "e" {
			p := val / resp.Coindata.EUR.Value
			fmt.Println("Pound conversion: ", val, "€ is ", p, "£")
		} else {
			fmt.Println("Invalid argument given")
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

	cur := flag.String("c", "p", "Currency, p or e")
	val := flag.Float64("v", 1.0, "Value to convert")
	flag.Parse()

	if *cur != "p" && *cur != "e" {
		panic("Invalid currency option. Options: p|e")
	}

	getExchange(*cur, *val)	

}

