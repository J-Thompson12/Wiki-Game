package link

import (
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// CleanPath cleans the string to make it a valid wiki path
func CleanPath(s string) string {
	s = strings.TrimRight(s, "\n")
	s = strings.Replace(s, " ", "_", 1)
	s = strings.TrimRight(s, " ")
	return s
}

// CreateURL checks if the path entered is a full url or just the path and turns it into a url
func CreateURL(s string) string {
	if strings.Contains(s, "https://en.wikipedia.org/wiki/") {
		return s
	}
	path := "https://en.wikipedia.org/wiki/" + s
	return path
}

// IsValidURL checks if the paths entered are valid urls
func IsValidURL(url string) bool {
	if url != "" {
		if !strings.Contains(url, "https://en.wikipedia.org/wiki/") {
			return false
		}
		pageSource := read(url)
		s := strings.NewReader(pageSource)
		doc, err := html.Parse(s)
		if err != nil {
			panic(err)
		}
		var f func(*html.Node)
		isValid := true
		// Looks for the html element that shows its not a valid wiki url
		f = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "div" {
				for _, attr := range n.Attr {
					if attr.Key == "class" {
						if attr.Val == "noarticletext mw-content-ltr" {
							isValid = false
							return
						}
					}
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(doc)
		return isValid
	}
	return false
}

func read(url string) string {
	var s string
	resp, err := http.Get(url)
	// if an error is returned pause all threads for 20 seconds and resend
	if err != nil {
		pause <- 1
		fmt.Println(err)
		time.Sleep(20 * time.Second)
		<-pause
		return read(url)
	}
	if resp.Body == nil {
		fmt.Println("resp somehow is nil")
	}
	defer resp.Body.Close()
	html, err := ioutil.ReadAll(resp.Body)
	// if an error is returned pause all threads for 20 seconds and resend
	if err != nil {
		pause <- 1
		fmt.Println(err)
		time.Sleep(20 * time.Second)
		fmt.Println("its been 20 seconds http" + url)
		<-pause
		return read(url)
	}
	s = string(html)
	return s
}
