package main

import (
	"./integration"
	"./web"
	"./web/auth"
	"google.golang.org/appengine"

	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./views/")))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/api/auth", web.HandlerAPIGetAuthToken)
	http.HandleFunc("/api/portfolio", auth.SecurityMiddleware(web.HandleAPIPortfolio))
	http.HandleFunc("/api/performance", auth.SecurityMiddleware(web.HandleAPIPerformance))
	http.HandleFunc("/api/performance/btc", auth.SecurityMiddleware(web.HandleAPIPerformanceBtc))
	http.HandleFunc("/api/performance/bch", auth.SecurityMiddleware(web.HandleAPIPerformanceBch))
	http.HandleFunc("/api/history", auth.SecurityMiddleware(web.HandleAPIAccountHistory))

	integration.InitializeCurrencyRateUpdates()

	appengine.Main()
}
