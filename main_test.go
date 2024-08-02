package main

import (
	"fmt"
	"testing"
)

func TestGetRate(t *testing.T) {
	var rate float64
	var err error

	err = UpdateCurrencies()
	if err != nil {
		t.Error(err)
	}

	rate, err = GetCurrency("NBU", "USD")
	if err != nil {
		t.Error(err)
	}
	if rate <= 0 {
		t.Error("Rate <= 0")
	}
	fmt.Println("OK")

	rate, err = GetCurrency("CNB", "USD")
	if err != nil {
		t.Error(err)
	}
	if rate <= 0 {
		t.Error("Rate <= 0")
	}
	fmt.Println("OK")

}
