package request

type ScrapReq struct {
	ResourceId int64  `json:"resourceId"`
	Keyword    string `json:"keyword"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
}
