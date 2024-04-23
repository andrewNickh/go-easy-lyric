package liriklagurohani

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
	ResourceName  = "lirik_lagu_rohani"
	baseUrl       = "https://liriklagukristen.id/"
	searchBaseURL = "https://liriklagukristen.id/?s="
)

var LirikLagu = new(_lirikLagu)

type _lirikLagu struct {
}

func (l *_lirikLagu) Scrape(req request.ScrapReq) ([]*response.ScrapResp, int, error) {
	searchUrl := l.generateSearchUrl(req.Keyword, req.Page)

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

	links := l.getAllSearchResultUrl(doc, req.Limit)

	if len(links) == 0 {
		return nil, 0, errors.New("song not found")
	}

	// scrape it
	var songs []*response.ScrapResp
	for _, link := range links {
		song, err := l.startScrape(link)
		if err != nil {
			log.Error(err)
		}
		songs = append(songs, song)
	}

	return songs, len(songs), nil
}

func (l *_lirikLagu) generateSearchUrl(keyword string, page int) string {
	inputs := strings.Split(strings.ToLower(strings.TrimSpace(keyword)), " ")
	input := strings.Join(inputs, "+")
	if page > 1 {
		return fmt.Sprintf("%spage/%d/?s=%s", baseUrl, page, input)
	}
	return searchBaseURL + input
}

func (l *_lirikLagu) getAllSearchResultUrl(n *html.Node, limit int) []string {
	var lyricLinks []string

	// declare extractLinks func
	var extractLinks func(*html.Node)
	var stop bool
	extractLinks = func(n *html.Node) {
		if n.Data == "section" {
			if len(n.Attr) > 0 {
				if n.Attr[0].Key == "class" && n.Attr[0].Val == "error-404 not-found" {
					stop = true
				}
			}
		}
		if n.Data == "div" {
			if len(n.Attr) > 0 {
				if n.Attr[0].Key == "class" && n.Attr[0].Val == "ast-pagination" {
					stop = true
				}
			}
		}

		if stop {
			return
		}

		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if strings.Contains(attr.Val, baseUrl) && attr.Val != baseUrl {
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

	// exec
	extractLinks(n)
	l.removeDuplicates(&lyricLinks)
	return lyricLinks
}

func (l *_lirikLagu) removeDuplicates(strSlice *[]string) {
	encountered := map[string]bool{}
	result := []string{}

	for _, v := range *strSlice {
		if encountered[v] == true {
			// Do not add duplicate string to result
			continue
		} else {
			// Add string to map and result slice
			encountered[v] = true
			result = append(result, v)
		}
	}

	*strSlice = result
}

func (l *_lirikLagu) startScrape(lyricLink string) (song *response.ScrapResp, err error) {
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
