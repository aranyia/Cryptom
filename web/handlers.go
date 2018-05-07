package web

import (
	"../api"
	"../integration"
	"./auth"

	"encoding/json"
	"net/http"
)

var exchange integration.Exchange = &integration.GDAXExchange{BaseCurrency: "EUR", ValueCurrencies: []string{"EUR", "HUF", "AED"}}

func HandleAPIPortfolio(respWriter http.ResponseWriter, request *http.Request) {
	response := exchange.GetPortfolio()

	respBody, _ := json.Marshal(response)

	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Write(respBody)
}

func HandleAPIPerformance(respWriter http.ResponseWriter, request *http.Request) {
	response := exchange.GetCurrentStakePerformanceSummary()

	respBody, _ := json.Marshal(response)

	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Write(respBody)
}

func HandleAPIPerformanceBtc(respWriter http.ResponseWriter, request *http.Request) {
	response := exchange.GetCurrentStakePerformance("BTC-EUR")

	respBody, _ := json.Marshal(response)

	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Write(respBody)
}

func HandleAPIPerformanceBch(respWriter http.ResponseWriter, request *http.Request) {
	response := exchange.GetCurrentStakePerformance("BCH-EUR")

	respBody, _ := json.Marshal(response)

	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Write(respBody)
}

func HandleAPIAccountHistory(respWriter http.ResponseWriter, request *http.Request) {
	response := exchange.GetAccountHistory("BTC")

	respBody, _ := json.Marshal(response)

	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Write(respBody)
}

func HandlerAPIGetAuthToken(respWriter http.ResponseWriter, request *http.Request) {
	var authRequest api.AuthenticationRequest

	requestBody := make([]byte, request.ContentLength)
	request.Body.Read(requestBody)
	if json.Unmarshal(requestBody, &authRequest) != nil || authRequest.ApiKey == "" {
		http.Error(respWriter, "invalid request", http.StatusBadRequest)
		return
	}

	if auth.IsAuthenticationValid(authRequest.ApiKey) {
		tokenString, _ := auth.GenerateAuthToken()

		respBody, _ := json.Marshal(api.AuthenticationResponse{AccessToken: tokenString})

		respWriter.Header().Add("Content-Type", "application/json")
		respWriter.Write(respBody)
	} else {
		http.Error(respWriter, "invalid API key", http.StatusUnauthorized)
	}
}
