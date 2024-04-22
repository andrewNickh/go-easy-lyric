package request

type GetById struct {
	Id int64 `json:"id" form:"id"`
}

type PageInfo struct {
	Page  int `json:"page" form:"page"`
	Limit int `json:"limit" form:"limit"`
}

type ScrapReq struct {
	ResourceId int64  `json:"resourceId"`
	Keyword    string `json:"keyword"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
}

type SaveSongReq struct {
	ResourceId int64  `json:"resourceId"`
	Url        string `json:"url"`
	Title      string `json:"title"`
	Lyric      string `json:"lyric"`
}

type UpdateSongReq struct {
	GetById
	SaveSongReq
}

type CreateResourceReq struct {
	Name   string `json:"name"`
	Status bool   `json:"status"`
	Url    string `json:"url"`
}

type UpdateResourceReq struct {
	GetById
	Name   string `json:"name"`
	Status bool   `json:"status"`
	Url    string `json:"url"`
}
