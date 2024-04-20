package request

type ScrapReq struct {
	Keyword string `json:"keyword"`
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
}
