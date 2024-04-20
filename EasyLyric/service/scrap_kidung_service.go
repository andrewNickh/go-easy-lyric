package service

import (
	"easy-lyric/EasyLyric/model/response"
	"easy-lyric/config"
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

var ScrapService = new(scrapService)

type scrapService struct {
}

func (s *scrapService) GetScrapService(keyword string, page, limit int) ([]*response.ScrapResp, int, error) {
	searchUrl := s.generateSearchUrl(keyword, page)

	// get search result
	resp, err := http.Get(searchUrl)
	if err != nil {
		log.Error(err)
		return nil, 0, errors.New("search song error")
	}
	defer resp.Body.Close()

	// get html node
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, 0, errors.New("search song error")
	}

	// get all search result url by html node
	links := s.getAllSearchResultUrl(doc, limit)

	if len(links) == 0 {
		return nil, 0, errors.New("song not found")
	}

	// scrape it
	var songs []*response.ScrapResp
	for _, link := range links {
		song, err := s.scrape(link)
		if err != nil {
			log.Error(err)
		}
		songs = append(songs, song)
	}

	return songs, len(songs), nil
}

func (s *scrapService) generateSearchUrl(keyword string, page int) string {
	inputs := strings.Split(strings.ToLower(strings.TrimSpace(keyword)), " ")
	input := strings.Join(inputs, "+")
	if page > 0 {
		return fmt.Sprintf("%s/page/%d", config.Instance.SearchBaseURL+input, page)
	}
	return config.Instance.SearchBaseURL + input
}

func (s *scrapService) getAllSearchResultUrl(n *html.Node, limit int) []string {
	var lyricLinks []string

	// declare extractLinks func
	var extractLinks func(*html.Node)
	extractLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if strings.Contains(attr.Val, config.Instance.BaseURL) && strings.Contains(attr.Val, "#more-") {
					lyricLinks = append(lyricLinks, attr.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if len(lyricLinks) == limit {
				break
			}
			extractLinks(c)
		}
	}

	// exec
	extractLinks(n)
	return lyricLinks
}

func (s *scrapService) scrape(lyricLink string) (song *response.ScrapResp, err error) {
	if lyricLink == "" {
		return nil, errors.New("url not found")
	}

	song = &response.ScrapResp{
		Url: lyricLink,
	}
	//log.Infof("=== %s ===", lyricLink)

	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{lyricLink},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			r.HTMLDoc.Find("h1.entry-title").Each(func(i int, s *goquery.Selection) {
				title := s.Text()
				if title != "" {
					song.Title = strings.Replace(title, "Lirik & Chord Lagu ", "", 1)
				}
			})
			r.HTMLDoc.Find("div.entry-content").Each(func(i int, s *goquery.Selection) {
				lyric := s.Find("p").Text()
				if lyric != "" {
					song.Lyric = lyric
				}
			})
		},
	}).Start()

	return song.Render(), nil
}
