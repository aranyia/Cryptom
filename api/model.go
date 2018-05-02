package api

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

type LastTradePerformance struct {
	Unit          string  `json:"unit"`
	Amount        float64 `json:"amount,omitempty"`
	Currency      string  `json:"currency"`
	ValuePayed    float64 `json:"valuePayed"`
	ValueCurrent  float64 `json:"valueCurrent"`
	ValueChange   float64 `json:"valueChange"`
	PercentChange float64 `json:"percentChange"`
}

type AuthenticationRequest struct {
	ApiKey string `json:"apiKey"`
}

type AuthenticationResponse struct {
	AccessToken string `json:"accessToken"`
}
