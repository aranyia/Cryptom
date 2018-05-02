package integration

import (
	"../api"
	gdax "github.com/preichenberger/go-gdax"
)

type Exchange interface {
	GetBaseCurrency() string

	GetAccountHistory(string) []gdax.LedgerEntry

	GetPortfolio() api.Portfolio

	GetCurrentStakePerformance(string, []gdax.LedgerEntry) api.LastTradePerformance

	GetCurrentStakePerformanceSummary() []api.LastTradePerformance
}
