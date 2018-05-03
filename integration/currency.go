package integration

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// sellRateModifier to calculate realistic market value of the currency.
const sellRateModifier = 0.99

// updateCycle the time period to fire the update of cached currency rates
const updateCycle = 60 * time.Minute

var updateTicker *time.Ticker

var currencyRatesMap = make(map[string]float64)

// GetCurrencyRate returns the value of the currency rate between the from
// and to currency codes, i.e. from "EUR" to "USD".
// Reads from the cached values first, then tries to load the value from
// the designated source.
func GetCurrencyRate(from string, to string) float64 {
	mapKey := toMapKey(from, to)

	var value = currencyRatesMap[mapKey]
	if value == 0 {
		updateCurrencyRate(from, to, currencyRatesMap)
	}
	return value
}

// InitializeCurrencyRateUpdates sets up a go process to periodically update
// cached currency rates scheduled by a time.Ticker instance.
func InitializeCurrencyRateUpdates() {
	updateTicker = time.NewTicker(updateCycle)
	go func() {
		for {
			select {
			case <-updateTicker.C:
				updateCurrencyRates(currencyRatesMap)
			}
		}
	}()
	log.Println("Initialized currency rate updates")
}

// DisableCurrencyRateUpdates disables the go process which updates
// cached currency rates.
func DisableCurrencyRateUpdates() {
	updateTicker.Stop()
	log.Println("Disabled currency rate updates")
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

func toMapKey(from string, to string) string {
	return from + "-" + to
}

func updateCurrencyRate(from string, to string, currencyRatesMap map[string]float64) {
	value, err := loadCurrencyRate(from, to)
	if err == nil {
		mapKey := toMapKey(from, to)
		value = value * sellRateModifier
		currencyRatesMap[mapKey] = value
		log.Printf("Updated %s rate: %.2f", mapKey, value)
	}
}

func updateCurrencyRates(currencyRatesMap map[string]float64) {
	for currencyPairKey := range currencyRatesMap {
		currencyPair := strings.Split(currencyPairKey, "-")
		updateCurrencyRate(currencyPair[0], currencyPair[1], currencyRatesMap)
	}
}
