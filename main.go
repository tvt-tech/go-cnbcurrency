package main

import "C"

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gocolly/colly"
	"gopkg.in/ini.v1"
)

const (
	CNB_ALLOWED_DOMAINS string = "www.cnb.cz"
	NBU_ALLOWED_DOMAINS string = "bank.gov.ua"
	CNB_URL             string = "https://www.cnb.cz/en/financial-markets/foreign-exchange-market/central-bank-exchange-rate-fixing/central-bank-exchange-rate-fixing/"
	NBU_URL             string = "https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange/"
)

//

//export GetCurrencyC
func GetCurrencyC(bankID *C.char, code *C.char) float64 {
	value, err := GetCurrency(C.GoString(bankID), C.GoString(code))
	if err != nil {
		return 0
	}
	return value
}

//export UpdateCurrenciesC
func UpdateCurrenciesC() int {
	err := UpdateCurrencies()
	if err != nil {
		return -1
	}
	return 0
}

func main() {
	// UpdateCurrencies()
	// fmt.Println(GetCurrency("NBU", "USD"))
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
	if _, err := os.Stat("currencies.ini"); os.IsNotExist(err) {
		// If file does not exist, create a new configuration
		cfg = ini.Empty()
	} else {
		// Load the existing INI file
		cfg, err = ini.Load("currencies.ini")
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
	err = cfg.SaveTo("currencies.ini")
	if err != nil {
		return fmt.Errorf("failed to save INI file: %w", err)
	}

	return nil
}

func GetCurrency(bankID, cc string) (float64, error) {
	var cfg *ini.File
	var err error

	// Check if the INI file exists
	if _, err := os.Stat("currencies.ini"); os.IsNotExist(err) {
		// If file does not exist, return error
		return 0, err
	} else {
		// Load the existing INI file
		cfg, err = ini.Load("currencies.ini")
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
