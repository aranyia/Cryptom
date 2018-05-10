package api

import "time"

type Portfolio struct {
	Valuations []Valuation `json:"valuations,omitempty"`
	Stakes     []Stake     `json:"stakes"`
}

type Stake struct {
	Unit       string      `json:"unit"`
	Amount     float64     `json:"amount"`
	Valuations []Valuation `json:"valuations,omitempty"`
}

type Valuation struct {
	Currency  string  `json:"currency"`
	Value     float64 `json:"value"`
	ValueUnit float64 `json:"valueUnit,omitempty"`
}

type ActiveTradePerformance struct {
	StartedAt     time.Time `json:"startedAt,string"`
	Unit          string    `json:"unit"`
	Amount        float64   `json:"amount,omitempty"`
	Currency      string    `json:"currency"`
	ValuePaid     float64   `json:"valuePaid"`
	ValueCurrent  float64   `json:"valueCurrent"`
	ValueChange   float64   `json:"valueChange"`
	PercentChange float64   `json:"percentChange"`
}

type LedgerEntry struct {
	CreatedAt time.Time `json:"createdAt,string"`
	Amount    float64   `json:"amount,string"`
	Balance   float64   `json:"balance,string"`
	Type      string    `json:"type"`
	OrderID   string    `json:"orderId"`
	ProductID string    `json:"productId"`
}

type AuthenticationRequest struct {
	ApiKey string `json:"apiKey"`
}

type AuthenticationResponse struct {
	AccessToken string `json:"accessToken"`
}
