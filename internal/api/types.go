package api

type QuotesResponse struct {
	Status  string `json:"s"`
	Data    []Data `json:"d,omitempty"`
	Message string `json:"message,omitempty"`
}
type Cmd struct {
	T  int     `json:"t"`
	O  float64 `json:"o"`
	H  float64 `json:"h"`
	L  float64 `json:"l"`
	C  float64 `json:"c"`
	V  int     `json:"v"`
	Tf string  `json:"tf"`
}
type Values struct {
	Change           float64 `json:"ch"`
	ChangePercentage float64 `json:"chp"`
	LastPrice        float64 `json:"lp"`
	Spread           float64 `json:"spread"`
	Ask              float64 `json:"ask"`
	Bid              float64 `json:"bid"`
	OpenPrice        float64 `json:"open_price"`
	HighPrice        float64 `json:"high_price"`
	LowPrice         float64 `json:"low_price"`
	PrevClosePrice   float64 `json:"prev_close_price"`
	Volume           int     `json:"volume"`
	ShortName        string  `json:"short_name"`
	Exchange         string  `json:"exchange"`
	Description      string  `json:"description"`
	OriginalName     string  `json:"original_name"`
	Symbol           string  `json:"symbol"`
	FyToken          string  `json:"fyToken"`
	Time             int     `json:"tt"`
	Cmd              Cmd     `json:"cmd"`
}
type Data struct {
	Name   string `json:"n"`
	Status string `json:"s"`
	Values Values `json:"v"`
}

type AuthCodeResponse struct {
	Status      string      `json:"status"`
	Code        int         `json:"code"`
	Message     interface{} `json:"message,omitempty"`
	AccessToken string      `json:"access_token"`
}
