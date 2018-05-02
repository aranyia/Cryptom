package integration

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// sellRateModifier to calculate realistic market value of the currency.
const sellRateModifier = 0.99

var currencyRatesMap = make(map[string]float64)

// GetCurrencyRate returns the value of the currency rate between the from
// and to currency codes, i.e. from "EUR" to "USD".
// Reads from the cached values first, then tries to load the value from
// the designated source.
func GetCurrencyRate(from string, to string) float64 {
	mapKey := from + "-" + to

	var value = currencyRatesMap[mapKey]
	if value == 0 {
		value, err := loadCurrencyRate(from, to)
		if err == nil {
			value = value * sellRateModifier
			currencyRatesMap[mapKey] = value
		}
	}
	return value
}

func loadCurrencyRate(from string, to string) (float64, error) {
	const endpointURI = "https://free.currencyconverterapi.com/api/v3/convert"

	currencyPairKey := from + "_" + to
	endpointURIWithParams := endpointURI + "?compact=ultra&q=" + currencyPairKey

	resp, err := http.Get(endpointURIWithParams)
	if err != nil {
		return 0, err
	}
	var currencyRateResponse map[string]float64

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &currencyRateResponse)

	return currencyRateResponse[currencyPairKey], nil
}
