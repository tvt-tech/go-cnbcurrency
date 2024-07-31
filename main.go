package main

import "C"

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gocolly/colly"
	"gopkg.in/ini.v1"
)

const (
	CACHE_FILENAME      string = "currencies.ini"
	CNB_ALLOWED_DOMAINS string = "www.cnb.cz"
	NBU_ALLOWED_DOMAINS string = "bank.gov.ua"
	CNB_URL             string = "https://www.cnb.cz/en/financial-markets/foreign-exchange-market/central-bank-exchange-rate-fixing/central-bank-exchange-rate-fixing/"
	NBU_URL             string = "https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange/"
)

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

func main() {}

func GetCacheUpdateTime() int64 {
	if time, err := getFileModTime(CACHE_FILENAME); err == nil {
		return time.Unix()
	}
	return -1
}

func UpdateCurrencies() error {

	cnbCurrencies, err := FetchCNBCurrencies()
	if err != nil {
		fmt.Println("CNB currencies update Error", err)
	} else {
		err = CacheCurrencies("CNB", cnbCurrencies)
		if err != nil {
			fmt.Println("Error on CNB cache", err)
		}
	}

	nbuCurrencies, err := FetchNBUCurrencies()
	if err != nil {
		fmt.Println("NBU currencies update Error", err)
	} else {
		err = CacheCurrencies("NBU", nbuCurrencies)
		if err != nil {
			fmt.Println("Error on NBU cache", err)
		}
	}
	return err
}

func FetchNBUCurrencies() (map[string]float64, error) {
	output := make(map[string]float64)
	var scrapeErr error

	c := colly.NewCollector(colly.AllowedDomains(NBU_ALLOWED_DOMAINS)) // Adjust domain if needed

	c.OnXML("//currency", func(x *colly.XMLElement) {
		cc := x.ChildText("cc")
		rateStr := x.ChildText("rate")

		rate, err := strconv.ParseFloat(rateStr, 64)
		if err == nil {
			output[cc] = rate
		} else {
			fmt.Println("Can't parse %w=%w", cc, rateStr)
		}
	})

	c.OnError(func(_ *colly.Response, err error) {
		scrapeErr = err
	})

	err := c.Visit(NBU_URL)
	if err != nil {
		scrapeErr = err
	}

	if scrapeErr != nil {
		return output, scrapeErr
	}

	return output, nil

}

func FetchCNBCurrencies() (map[string]float64, error) {
	output := make(map[string]float64)

	var scrapeErr error

	c := colly.NewCollector(
		colly.AllowedDomains(CNB_ALLOWED_DOMAINS),
	)

	c.OnHTML(".currency-table", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, tr *colly.HTMLElement) {
			cc := tr.ChildText("td:nth-of-type(4)")
			rateStr := tr.ChildText("td:nth-of-type(5)")
			rate, err := strconv.ParseFloat(rateStr, 64)
			if err == nil {
				output[cc] = rate
			} else {
				fmt.Println("Can't parse %w=%w", cc, rateStr)
			}
		})
	})

	c.OnError(func(_ *colly.Response, err error) {
		scrapeErr = err
	})

	err := c.Visit(CNB_URL)
	if err != nil {
		scrapeErr = err
	}

	if scrapeErr != nil {
		return output, scrapeErr
	}

	return output, nil
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
	err = cfg.SaveTo(CACHE_FILENAME)
	if err != nil {
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
		return 0, fmt.Errorf("currencies.ini has not %w key in %w section", cc, bankID)
	}

	strValue := section.Key(cc)
	floatValue, err := strValue.Float64()
	if err != nil {
		return 0, fmt.Errorf("Parse value error")
	}
	return floatValue, nil
}
