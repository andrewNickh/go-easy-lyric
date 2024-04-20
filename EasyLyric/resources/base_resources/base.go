package base_resources

import (
	"easy-lyric/EasyLyric/model/request"
	"easy-lyric/EasyLyric/model/response"
)

type Source interface {
	Scrape(req request.ScrapReq) ([]*response.ScrapResp, int, error)
}
