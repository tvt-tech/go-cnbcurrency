package main

// import "C"

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

const (
	CNB_ALLOWED_DOMAINS string = "www.cnb.cz"
	NBU_ALLOWED_DOMAINS string = "bank.gov.ua"
	CNB_URL             string = "https://www.cnb.cz/en/financial-markets/foreign-exchange-market/central-bank-exchange-rate-fixing/central-bank-exchange-rate-fixing/"
	NBU_URL             string = "https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange/"
)

//export GetUsd
func GetUsd() float64 {
	return GetByCode("USD")
}

//export GetEur
func GetEur() float64 {
	return GetByCode("EUR")
}

// //export GetCurrency
// func GetCurrency(code *C.char) float64 {
// 	return GetByCode(C.GoString(code))
// }

func main() {
	//fmt.Println(GetUsd())
	//fmt.Println(GetEur())
	fmt.Println(GetNBUbyCode("USD"))
	fmt.Println(GetNBUbyCode("EUR"))
}

func GetNBUbyCode(input string) float64 {
	var output float64 = -1.0

	// Create a new collector
	c := colly.NewCollector(colly.AllowedDomains(NBU_ALLOWED_DOMAINS)) // Adjust domain if needed

	// Set up XML parsing for currency elements within exchange
	c.OnXML("//currency", func(x *colly.XMLElement) {
		code := x.ChildText("cc")
		rateStr := x.ChildText("rate")

		// Debug output
		// fmt.Printf("Found currency element with cc: %s, rate: %s\n", code, rateStr)

		// Check if the 'cc' element matches the input
		if code == input {
			// Convert rate to float64
			rate, err := strconv.ParseFloat(rateStr, 64)
			if err == nil {
				output = rate
				// fmt.Printf("Rate for %s: %f\n", input, output)
			} else {
				// fmt.Printf("Error converting rate: %v\n", err)
			}
		}
	})

	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Printf("Visiting URL: %s\n", r.URL)
	// })

	// c.OnError(func(r *colly.Response, err error) {
	// 	fmt.Printf("Request failed with status code %d and error: %v\n", r.StatusCode, err)
	// })

	// Start the request
	err := c.Visit(NBU_URL)
	if err != nil {
		fmt.Println("Error visiting URL:", err)
	}

	return output
}

func GetByCode(input string) float64 {
	var output float64 = -1.0

	c := colly.NewCollector(colly.AllowedDomains(CNB_ALLOWED_DOMAINS))
	c.OnHTML(".currency-table", func(e *colly.HTMLElement) {

		e.ForEach("tr", func(i int, tr *colly.HTMLElement) {
			code := tr.ChildText("td:nth-of-type(4)")
			rate, err := strconv.ParseFloat(tr.ChildText("td:nth-of-type(5)"), 64)
			if err == nil {
				if strings.EqualFold(code, input) {
					output = rate
				}
			}
		})

	})
	c.Visit(CNB_URL)
	return output
}
