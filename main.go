package main

import "C"

import (
	// "fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

//export GetUsd
func GetUsd() float64 {
	return GetByCode("USD")
}

//export GetEur
func GetEur() float64 {
	return GetByCode("EUR")
}

//export GetCurrency
func GetCurrency(code *C.char) float64 {
	return GetByCode(C.GoString(code))
}

func main() {
	// fmt.Println(GetUsd())
	// fmt.Println(GetEur())

}

func GetByCode(input string) float64 {
	var output float64 = -1.0

	c := colly.NewCollector(colly.AllowedDomains("www.cnb.cz"))
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
	c.Visit("https://www.cnb.cz/en/financial-markets/foreign-exchange-market/central-bank-exchange-rate-fixing/central-bank-exchange-rate-fixing/")
	return output
}
