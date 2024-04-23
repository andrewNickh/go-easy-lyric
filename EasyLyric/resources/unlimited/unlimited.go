package unlimited

import (
	"easy-lyric/EasyLyric/model/request"
	"easy-lyric/EasyLyric/model/response"
	"easy-lyric/util/log"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

const (
	ResourceName  = "unlimited_worship"
	baseUrl       = "https://unlimitedworship.org/"
	searchBaseURL = "https://unlimitedworship.org/search?type=1&key="
)

var Unlimited = new(_unlimited)

type _unlimited struct {
}

func (u *_unlimited) Scrape(req request.ScrapReq) ([]*response.ScrapResp, int, error) {
	searchUrl := u.generateSearchUrl(req.Keyword, req.Page)

	resp, err := http.Get(searchUrl)
	log.Info(searchUrl)
	if err != nil {
		log.Error(err)
		return nil, 0, errors.New("search song error")
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, 0, errors.New("search song error")
	}

	links := u.getAllSearchResultUrl(doc, req.Limit)

	if len(links) == 0 {
		return nil, 0, errors.New("song not found")
	}

	// scrape it
	var songs []*response.ScrapResp
	for _, link := range links {
		song, err := u.startScrape(link)
		if err != nil {
			log.Error(err)
		}
		songs = append(songs, song)
	}

	return songs, len(songs), nil
}

func (u *_unlimited) generateSearchUrl(keyword string, page int) string {
	inputs := strings.Split(strings.ToLower(strings.TrimSpace(keyword)), " ")
	input := strings.Join(inputs, "+")
	input = searchBaseURL + input
	if page > 1 {
		input = fmt.Sprintf("%s&page=%d", input, page)
	}
	return input
}

func (u *_unlimited) getAllSearchResultUrl(n *html.Node, limit int) []string {
	var lyricLinks []string

	var extractLinks func(*html.Node)
	var stop bool
	extractLinks = func(n *html.Node) {
		if n.Data == "div" {
			if len(n.Attr) > 0 {
				if n.Attr[0].Key == "" && n.Attr[0].Val == "" {
					stop = true
				}
			}
		}

		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if strings.Contains(attr.Val, baseUrl+"songs/detail") {
					lyricLinks = append(lyricLinks, attr.Val)
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if !stop {
				if len(lyricLinks) == limit {
					break
				}
				extractLinks(c)
			} else {
				break
			}
		}
	}

	extractLinks(n)
	return lyricLinks
}

func (u *_unlimited) startScrape(lyricLink string) (song *response.ScrapResp, err error) {
	if lyricLink == "" {
		return nil, errors.New("url not found")
	}

	song = &response.ScrapResp{
		Url: lyricLink,
	}

	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{lyricLink},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			r.HTMLDoc.Find("div.now-playing-title").Each(func(i int, s *goquery.Selection) {
				title := s.Text()
				if title != "" {
					song.Title = title
				}
			})

			r.HTMLDoc.Find("div.lyric-content-container").Each(func(i int, s *goquery.Selection) {
				lyric := s.Find("pre.lyric-content").Text()
				if lyric != "" {
					song.Lyric = lyric
				}
			})
		},
	}).Start()

	return song.Render(), nil
}
