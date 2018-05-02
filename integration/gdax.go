package integration

import (
	"fmt"
	"os"
	"strings"

	"../api"
	"github.com/preichenberger/go-gdax"
)

type GDAXExchange struct {
	BaseCurrency    string
	ValueCurrencies []string
}

func (exchange *GDAXExchange) GetBaseCurrency() string {
	return exchange.BaseCurrency
}

var client = initClient()

func initClient() *gdax.Client {
	return gdax.NewClient(os.Getenv("GDAX_SECRET"), os.Getenv("GDAX_KEY"), os.Getenv("GDAX_PASSPHRASE"))
}

func (exchange *GDAXExchange) getStakes() (stakes []api.Stake, sumValuations []api.Valuation, err error) {
	baseCurrency := exchange.GetBaseCurrency()
	accounts, err := client.GetAccounts()

	var sumValue float64
	for _, account := range accounts {
		stake := api.Stake{Unit: account.Currency, Amount: account.Balance}

		var valuations []api.Valuation

		var baseValue, unitPrice float64
		if account.Currency != baseCurrency {
			ticker, _ := client.GetTicker(account.Currency + "-" + baseCurrency)
			unitPrice = ticker.Price
			baseValue = account.Balance * unitPrice
		} else {
			unitPrice = 0
			baseValue = account.Balance
		}

		for _, currency := range exchange.ValueCurrencies {
			if currency == account.Currency {
				continue
			}
			valuations = append(valuations, getValuation(baseCurrency, currency, baseValue, unitPrice))
		}

		stake.Valuations = valuations
		sumValue += baseValue

		if stake.Amount > 0 {
			stakes = append(stakes, stake)
		}
	}

	for _, currency := range exchange.ValueCurrencies {
		sumValuations = append(sumValuations, getValuation(baseCurrency, currency, sumValue, 0))
	}
	return stakes, sumValuations, err
}

func getValuation(baseCurrency string, valueCurrency string, baseValue float64, baseUnitPrice float64) api.Valuation {
	var rate float64
	if baseCurrency == valueCurrency {
		rate = 1
	} else {
		rate = GetCurrencyRate(baseCurrency, valueCurrency)
	}
	return api.Valuation{Currency: valueCurrency, Value: baseValue * rate, ValueUnit: baseUnitPrice * rate}
}

func (exchange *GDAXExchange) GetPortfolio() api.Portfolio {
	stakes, portfolioValuations, _ := exchange.getStakes()
	portfolio := api.Portfolio{Stakes: stakes, Valuations: portfolioValuations}

	return portfolio
}

func getFilledOrders() (buyOrders []gdax.Order, sellOrders []gdax.Order) {
	params := gdax.ListOrdersParams{Status: "done", Pagination: gdax.PaginationParams{Limit: 50}}

	var orders []gdax.Order
	client.ListOrders(params).NextPage(&orders)

	for _, order := range orders {
		if order.DoneReason == "filled" {
			switch order.Side {
			case "buy":
				buyOrders = append(buyOrders, order)
			case "sell":
				sellOrders = append(sellOrders, order)
			}
			fmt.Printf("Filled %s order at %s\n", order.Side, order.CreatedAt.Time().String())
		}
	}
	return buyOrders, sellOrders
}

func (exchange *GDAXExchange) GetAccountHistory(accountUnit string) (accountLedgerEntries []gdax.LedgerEntry) {
	accounts, _ := client.GetAccounts()
	for _, account := range accounts {
		if account.Currency == accountUnit {
			client.ListAccountLedger(account.Id, gdax.GetAccountLedgerParams{Pagination: gdax.PaginationParams{Limit: 50}}).NextPage(&accountLedgerEntries)
			return accountLedgerEntries
		}
	}
	return nil
}

func (exchange *GDAXExchange) GetCurrentStakePerformanceSummary() (performanceIndicators []api.LastTradePerformance) {
	products := []string{"BTC-EUR", "BCH-EUR"}

	for _, product := range products {
		cryptoUnit := strings.Split(product, "-")[0]
		performance := exchange.GetCurrentStakePerformance(product, exchange.GetAccountHistory(cryptoUnit))
		performanceIndicators = append(performanceIndicators, performance)
	}
	return performanceIndicators
}

func (exchange *GDAXExchange) GetCurrentStakePerformance(productId string, ledgerEntries []gdax.LedgerEntry) api.LastTradePerformance {
	orderIDs := map[string]bool{}

	for _, entry := range ledgerEntries {
		if entry.Balance == 0 {
			break
		}
		orderIDs[entry.Details.OrderId] = true
	}

	var sumAmount float64
	var sumPayed float64
	for orderID := range orderIDs {
		order, _ := client.GetOrder(orderID)
		sumAmount += order.FilledSize
		sumPayed += order.ExecutedValue
	}
	ticker, _ := client.GetTicker(productId)

	sumCurrentValue := ticker.Price * sumAmount
	valueChange := sumCurrentValue - sumPayed
	percentChange := valueChange / sumPayed * 100

	productTypes := strings.Split(productId, "-")

	return api.LastTradePerformance{Unit: productTypes[0], Amount: sumAmount, Currency: productTypes[1],
		ValuePayed:    sumPayed,
		ValueCurrent:  sumCurrentValue,
		ValueChange:   valueChange,
		PercentChange: percentChange}
}
