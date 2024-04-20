package response

import "strings"

type ScrapResp struct {
	Url   string `json:"url"`
	Title string `json:"title"`
	Lyric string `json:"lyric"`
}

func (s *ScrapResp) Render() *ScrapResp {
	s.Lyric = strings.TrimSpace(s.Lyric)
	ss := strings.Split(s.Lyric, "Songwriter")
	s.Lyric = ss[0]
	return s
}
