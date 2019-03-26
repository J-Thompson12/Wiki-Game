package link

import (
	"golang.org/x/net/html"
	"strings"
	"sync"
)

var searched = make(map[string]string)
var searchedMutex = sync.RWMutex{}

// Link creates an object with the link path
type Link struct {
	Href string
}

// builds each link by parsing out the path
func buildLink(n *html.Node) Link {
	var ret Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			if strings.HasPrefix(attr.Val, "/wiki") && !strings.Contains(attr.Val, ":") && !strings.Contains(attr.Val, "%") {
				path := strings.TrimLeft(attr.Val, "/wiki/")
				// puts all searched paths into a map so we dont search the same link twice
				searchedMutex.RLock()
				_, ok := searched[path]
				searchedMutex.RUnlock()
				if !ok {
					ret.Href = path
					searchedMutex.Lock()
					searched[path] = path
					searchedMutex.Unlock()
					break
				}
			}
		}
	}
	return ret
}
