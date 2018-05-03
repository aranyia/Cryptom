package main

import (
	"./integration"
	"./web"
	"./web/auth"

	"net/http"

	"google.golang.org/appengine"
)

func main() {
	http.HandleFunc("/", web.HandlePageIndex)
	http.HandleFunc("/api/auth", web.HandlerAPIGetAuthToken)
	http.HandleFunc("/api/portfolio", auth.SecurityMiddleware(web.HandleAPIPortfolio))
	http.HandleFunc("/api/performance", auth.SecurityMiddleware(web.HandleAPIPerformance))
	http.HandleFunc("/api/performance/btc", auth.SecurityMiddleware(web.HandleAPIPerformanceBtc))
	http.HandleFunc("/api/performance/bch", auth.SecurityMiddleware(web.HandleAPIPerformanceBch))
	http.HandleFunc("/api/history", auth.SecurityMiddleware(web.HandleAPIAccountHistory))

	integration.InitializeCurrencyRateUpdates()

	appengine.Main()
}
