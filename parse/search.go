package link

import (
	"fmt"
	"strings"
	"sync"
)

var wg sync.WaitGroup
var start, end string
var correctRoute []string
var shutdown = make(chan struct{})
var pause = make(chan int, 1)

// GetRoute begins the search for the end path
func GetRoute(url string, startLink string, endLink string) []string {
	start = startLink
	end = endLink
	links, err := parse(url)
	if err != nil {
		return nil
	}
	route := make([]string, 0)
	route = append(route, start)
	searchLinks(links, route)
	wg.Wait()
	return correctRoute
}

func searchLinks(links []Link, route []string) {
	var done = false
	foundPath, endPath := foundPath(links)

	if foundPath {
		fmt.Println("You found it " + endPath)
		correctRoute = append(route, endPath)
		close(shutdown)
		return
	} else if len(route) > 6 { // dont search past 7 links deep
		fmt.Printf("You have gone farther than 6 %v\n ", route)
		return
	} else {
		for _, link := range links {
			path := "https://en.wikipedia.org/wiki/" + link.Href
			// gets all links from the link path
			newLinks, err := parse(path)
			if err != nil {
				fmt.Println(err)
				return
			}
			// new goroutine for each link
			wg.Add(1)
			go func(link Link, route []string, links []Link) {

				newRoute := append(route, link.Href)
				defer wg.Done()
				select {
				case _ = <-shutdown:
					done = true
					return
				default:
					searchLinks(links, newRoute)
					return
				}
			}(link, route, newLinks)
			if done == true {
				return
			}
		}
	}
}

// checks if the path is the one we were searching for
func foundPath(links []Link) (bool, string) {
	for _, link := range links {
		linkPath := link.Href
		linkPath = strings.ToLower(linkPath)
		if linkPath == end {
			return true, linkPath

		}
	}
	return false, ""
}
