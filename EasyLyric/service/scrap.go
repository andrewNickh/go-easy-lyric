package service

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os"
	"strings"
)

var BaseURL = "https://www.kidung.com/search/"

type Song struct {
	Title string `json:"title"`
}

func main() {
	// read input from user
	fmt.Println("Input Song:")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	inputs := strings.Split(strings.ToLower(strings.TrimSpace(line)), " ")
	input := strings.Join(inputs, "+")
	//log.Println(input)
	link := BaseURL + input
	log.Println(link)

	resp, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// get the readmore button
	var lyricLinks []string
	var extractLinks func(*html.Node)
	extractLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if strings.Contains(attr.Val, "https://www.kidung.com/") && strings.Contains(attr.Val, "#more-") {
					lyricLinks = append(lyricLinks, attr.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractLinks(c)
			//if lyricLinks != nil {
			//	break
			//}
		}
	}

	extractLinks(doc)
	//log.Println(lyricLinks)

	if lyricLinks == nil {
		log.Fatal("Lyric Not Found")
	} else {
		// scrap the lyric
		for _, lyricLink := range lyricLinks {
			s, err := scrape(lyricLink)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("")
			fmt.Println("===========================================")
			fmt.Println(render(s).Title)
			fmt.Println("===========================================")
		}
	}

	//url := "https://quotes.toscrape.com"
	//c := colly.NewCollector()
	//
	//c.OnRequest(func(r *colly.Request) {
	//	r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")
	//	fmt.Println("Visiting", r.URL)
	//})
	//
	//c.OnResponse(func(r *colly.Response) {
	//	fmt.Println("Response code", r.StatusCode)
	//})
	//
	//c.OnError(func(r *colly.Response, err error) {
	//	fmt.Println("Error", err.Error())
	//})
	//
	//c.OnHTML(".quote", func(h *colly.HTMLElement) {
	//	quotes := h.DOM.Find(".text").Text()
	//	author := h.DOM.Find(".author").Text()
	//	fmt.Printf("%v\n by %v\n", quotes, author)
	//})
	//
	//c.Visit(url)
}

func render(song *Song) *Song {
	song.Title = strings.TrimSpace(song.Title)
	s := strings.Split(song.Title, "Songwriter")
	song.Title = s[0]
	return song
}

func scrape(lyricLink string) (song *Song, err error) {
	song = &Song{}
	urls := []string{
		lyricLink,
	}

	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: urls,
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			if r.StatusCode == http.StatusNotFound {
				err = errors.New("Lyric not found")
			} else {
				r.HTMLDoc.Find("div.entry-content").Each(func(i int, s *goquery.Selection) {
					title := s.Find("p").Text()
					if title != "" {
						song.Title = title
					}
				})
			}
		},
	}).Start()

	if err != nil {
		song = nil
		return
	}
	return
}

func quotesParse(g *geziyor.Geziyor, r *client.Response) {
	if r.StatusCode == http.StatusNotFound {
		log.Println("Lyric not found")
		g.Exports <- nil
		return
	} else {
		r.HTMLDoc.Find("div.entry-content").Each(func(i int, s *goquery.Selection) {
			title := s.Find("p").Text()
			//log.Println("\n", title)
			if title != "" {
				g.Exports <- map[string]interface{}{
					"title": strings.Trim(title, " "),
				}
			}
		})
	}

	//r.HTMLDoc.Find("article.product_pod").Each(func(i int, s *goquery.Selection) {
	//	g.Exports <- map[string]interface{}{
	//		"title": s.Find("h3").Text(),
	//	}
	//})
	//if href, ok := r.HTMLDoc.Find("li.next > a").Attr("href"); ok {
	//	g.Get(r.JoinURL(href), quotesParse)
	//}
}
