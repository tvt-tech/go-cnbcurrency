package main

import "C"

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"gopkg.in/ini.v1"
)

const (
	CACHE_FILENAME string = "currencies.ini"
	CNB_JSON_URI   string = "https://api.cnb.cz/cnbapi/exrates/daily/"
	NBU_JSON_URI   string = "https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange?json"
)

type NBURate struct {
	R030         int     `json:"r030"`
	Txt          string  `json:"txt"`
	Rate         float64 `json:"rate"`
	CC           string  `json:"cc"`
	ExchangeDate string  `json:"exchangedate"`
}

type NBURatesResponse []NBURate

type CNBRate struct {
	ValidFor     string  `json:"validFor"`
	Order        int     `json:"order"`
	Country      string  `json:"country"`
	Currency     string  `json:"currency"`
	Amount       int     `json:"amount"`
	CurrencyCode string  `json:"currencyCode"`
	Rate         float64 `json:"rate"`
}

// Define a struct to represent the entire response
type CNBRatesResponse struct {
	Rates []CNBRate `json:"rates"`
}

//export GetCurrencyC
func GetCurrencyC(bankID *C.char, code *C.char) C.double {
	if bankID == nil || code == nil {
		return C.double(-2)
	}

	// Convert C strings to Go strings
	bankIDStr := C.GoString(bankID)
	codeStr := C.GoString(code)

	// Call the function and handle errors
	value, err := GetCurrency(bankIDStr, codeStr)
	if err != nil {
		return C.double(-1)
	}
	return C.double(value)
}

//export UpdateCurrenciesC
func UpdateCurrenciesC() C.int {
	err := UpdateCurrencies()
	if err != nil {
		return C.int(-1)
	}
	return C.int(0)
}

//export GetCacheUpdateTimeC
func GetCacheUpdateTimeC() C.longlong {
	return C.longlong(GetCacheUpdateTime())
}

func main() {
	fmt.Println("CurUpd: ", UpdateCurrencies())
	fmt.Println("UpdTime: ", GetCacheUpdateTime())
	if rate, err := GetCurrency("NBU", "EUR"); err == nil {
		fmt.Println("NBU EUR: ", rate)
	} else {
		fmt.Println(err)
	}
	if rate, err := GetCurrency("CNB", "USD"); err == nil {
		fmt.Println("CNB USD: ", rate)
	} else {
		fmt.Println(err)
	}
}

func GetCacheUpdateTime() int64 {
	if time, err := getFileModTime(CACHE_FILENAME); err == nil {
		return time.Unix()
	}
	return -1
}

func UpdateCurrencies() error {
	var updateErr error
	cnbCurrencies, err := FetchCNBCurrencies()
	fmt.Println("CNB currencies update Error", err)
	if err != nil {
		updateErr = err
		fmt.Println("CNB currencies update Error", err)
	} else {
		err := CacheCurrencies("CNB", cnbCurrencies)
		if err != nil {
			updateErr = err
			fmt.Println("Error on CNB cache", err)
		}
	}

	nbuCurrencies, err := FetchNBUCurrencies()
	if err != nil {
		fmt.Println("NBU currencies update Error", err)
		updateErr = err
	} else {
		err := CacheCurrencies("NBU", nbuCurrencies)
		if err != nil {
			updateErr = err
			fmt.Println("Error on NBU cache", err)
		}
	}
	return updateErr
}

func FetchAPIEndpoint(url string) ([]byte, error) {
	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %w", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status code %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	return body, nil
}

func FetchNBUCurrenciesAPI() (NBURatesResponse, error) {
	body, err := FetchAPIEndpoint(NBU_JSON_URI)
	if err != nil {
		return nil, err
	}

	// Create a variable to hold the JSON data
	var data NBURatesResponse

	// Unmarshal (decode) the JSON into the struct
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w", err)
	}

	// Return the parsed data
	return data, nil
}

func FetchCNBCurrenciesAPI() (*CNBRatesResponse, error) {
	body, err := FetchAPIEndpoint(CNB_JSON_URI)
	if err != nil {
		return nil, err
	}

	// Create an instance of the struct to hold the JSON data
	var data CNBRatesResponse

	// Unmarshal (decode) the JSON into the struct
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w", err)
	}
	// Return the parsed data
	return &data, nil
}

func FetchNBUCurrencies() (map[string]float64, error) {
	output := make(map[string]float64)
	data, err := FetchNBUCurrenciesAPI()
	if err == nil {
		for _, rate := range data {
			cc := rate.CC
			rate := rate.Rate
			output[cc] = rate
		}
		return output, nil
	}
	return output, err
}

func FetchCNBCurrencies() (map[string]float64, error) {
	output := make(map[string]float64)
	data, err := FetchCNBCurrenciesAPI()
	if err == nil {
		for _, rate := range data.Rates {
			cc := rate.CurrencyCode
			rate := rate.Rate
			output[cc] = rate
		}
		return output, nil
	}
	return output, err
}

func CacheCurrencies(bankID string, ccs map[string]float64) error {
	var cfg *ini.File
	var err error

	// Check if the INI file exists
	if _, err := os.Stat(CACHE_FILENAME); os.IsNotExist(err) {
		// If file does not exist, create a new configuration
		cfg = ini.Empty()
	} else {
		// Load the existing INI file
		cfg, err = ini.Load(CACHE_FILENAME)
		if err != nil {
			return fmt.Errorf("failed to load INI file: %w", err)
		}
	}

	// Get or create the section
	section, err := cfg.GetSection(bankID)
	if err != nil {
		section, err = cfg.NewSection(bankID)
		if err != nil {
			return fmt.Errorf("failed to get section: %w", err)
		}
	}

	for cc, rate := range ccs {
		strValue := strconv.FormatFloat(rate, 'f', -1, 64)

		// Set or update the key's value in the section
		_, err = section.NewKey(cc, strValue)
		if err != nil {
			return fmt.Errorf("failed to set key value: %w", err)
		}
	}

	// Save the updated configuration back to the INI file
	if cfg.SaveTo(CACHE_FILENAME) != nil {
		return fmt.Errorf("failed to save INI file: %w", err)
	}

	return nil
}

func getFileModTime(filePath string) (time.Time, error) {
	// Check if file exists and get its info
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return time.Time{}, fmt.Errorf("file does not exist")
	} else if err != nil {
		return time.Time{}, err
	}

	// Return the modification time
	return fileInfo.ModTime(), nil
}

func GetCurrency(bankID, cc string) (float64, error) {
	var cfg *ini.File
	var err error

	// Check if the INI file exists
	if _, err := os.Stat(CACHE_FILENAME); os.IsNotExist(err) {
		// If file does not exist, return error
		return 0, err
	} else {
		// Load the existing INI file
		cfg, err = ini.Load(CACHE_FILENAME)
		if err != nil {
			return 0, fmt.Errorf("failed to load INI file: %w", err)
		}
	}

	if !cfg.HasSection(bankID) {
		return 0, fmt.Errorf("currencies.ini has not %w section", bankID)
	}

	section := cfg.Section(bankID)

	if !section.HasKey(cc) {
		return 0, fmt.Errorf("currencies.ini has not %s key in %s section", cc, bankID)
	}

	strValue := section.Key(cc)
	floatValue, err := strValue.Float64()
	if err != nil {
		return 0, fmt.Errorf("Parse value error")
	}
	return floatValue, nil
}
