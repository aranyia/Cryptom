package integration

import "../api"

// Exchange is the representation of a digital currency exchange.
// All exchanges integrated should implement the functions defined here.
type Exchange interface {
	GetBaseCurrency() string

	GetAccountHistory(string) []api.LedgerEntry

	GetPortfolio() api.Portfolio

	GetCurrentStakePerformance(string) api.LastTradePerformance

	GetCurrentStakePerformanceSummary() []api.LastTradePerformance
}
