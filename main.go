package main

import (
	"GoProjects/wiki/parse"
	"html/template"
	"net/http"
)

type searchDetails struct {
	Start string
	End   string
	Route []string
}

func (d searchDetails) runSearch() []string {
	startURL := link.CreateURL(d.Start)
	return link.GetRoute(startURL, d.Start, d.End)
}

func main() {
	tmpl := template.Must(template.ParseFiles("template.gohtml"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		details := searchDetails{
			Start: r.FormValue("start"),
			End:   r.FormValue("end"),
		}

		// do something with details
		start := link.CleanPath(details.Start)
		startURL := link.CreateURL(start)
		startValid := link.IsValidURL(startURL)

		end := link.CleanPath(details.End)
		endURL := link.CreateURL(end)
		endValid := link.IsValidURL(endURL)

		success := startValid && endValid
		var route []string
		if success {
			route = details.runSearch()
		}

		tmpl.Execute(w, struct {
			Route   []string
			Success bool
		}{route, success})
	})

	// http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method != http.MethodPost {
	// 		tmpl.Execute(w, nil)
	// 		return
	// 	}

	// 	// do something with details
	// 	start := link.CleanPath(details.Start)
	// 	startURL := link.CreateURL(start)
	// 	startValid := link.IsValidURL(startURL)

	// 	end := link.CleanPath(details.End)
	// 	endURL := link.CreateURL(end)
	// 	endValid := link.IsValidURL(endURL)

	// 	success := startValid && endValid
	// 	route := details.runSearch()
	// 	tmpl.Execute(w, struct {
	// 		Route   []string
	// 		Success bool
	// 	}{route, success})
	// })

	http.ListenAndServe(":8080", nil)
}

// import (
// 	"GoProjects/wiki/parse"
// 	"bufio"
// 	"fmt"
// 	"github.com/gorilla/mux"
// 	"net/http"
// 	"os"
// 	"strings"
// )

// func main() {
// 	var start string
// 	var end string
// 	var startValid = false
// 	var endValid = false
// 	var startURL string

// 	r := mux.NewRouter()
// 	r.HandleFunc("/main", func(w http.ResponseWriter, r *http.Request) {
// 		// get the book
// 		// navigate to the page
// 	})

// reader := bufio.NewReader(os.Stdin)
// for !startValid {
// 	fmt.Println("where do you want to start?")
// 	start, _ = reader.ReadString('\n')
// 	start = link.CleanPath(start)
// 	startURL = link.CreateURL(start)
// 	startValid = link.IsValidURL(startURL)

// 	if !startValid {
// 		fmt.Println("Not valid url")
// 	}
// }
// for !endValid {
// 	fmt.Println("where do you need to go?")
// 	end, _ = reader.ReadString('\n')
// 	end = link.CleanPath(end)
// 	endURL := link.CreateURL(end)
// 	endValid = link.IsValidURL(endURL)

// 	if !endValid {
// 		fmt.Println("Not valid url")
// 	}
// }

// startPath := strings.TrimRight(start, "/wiki/")
// endPath := strings.TrimRight(end, "/wiki/")
// endPath = strings.ToLower(endPath)
// // fmt.Println(startPath)
// // fmt.Println(endPath)

// route := link.GetRoute(startURL, startPath, endPath)
// if route != nil {
// 	fmt.Println("start")
// 	for num, v := range route {
// 		num = num + 1
// 		fmt.Print(num)
// 		fmt.Print(" = ")
// 		fmt.Println(v)
// 	}
// 	fmt.Println("end")
// } else {
// 	fmt.Println("was not found")
// }

//}
