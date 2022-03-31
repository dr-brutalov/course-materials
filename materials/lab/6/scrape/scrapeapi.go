package scrape

// scrapeapi.go HAS TEN TODOS - TODO_5-TODO_14 and an OPTIONAL "ADVANCED" ASK

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"golang.org/x/exp/slices"
)

//==========================================================================\\

// Helper function walk function, modfied from Chap 7 BHG to enable passing in of
// additional parameter http responsewriter; also appends items to global Files and
// if responsewriter is passed, outputs to http

func walkFn(w http.ResponseWriter) filepath.WalkFunc {

	var fileCount int64 = 0

	return func(path string, f os.FileInfo, err error) error {
		w.Header().Set("Content-Type", "application/json")

		for _, r := range regexes {
			if r.MatchString(path) {
				var tfile FileInfo
				dir, filename := filepath.Split(path)
				tfile.Filename = string(filename)
				tfile.Location = string(dir)

				// Go 1.18 for the win.
				if !slices.Contains(Files, tfile) {
					Files = append(Files, tfile)
					fileCount++
				}

				if w != nil && len(Files) > 0 {

					// Added a basic counter.
					w.Write([]byte(`"` + (strconv.FormatInt(fileCount, 10)) + `":  `))
					fileCount++
					json.NewEncoder(w).Encode(tfile)
					w.Write([]byte(`,`))

				}

				log.Printf("[+] HIT: %s\n", path)

			}

		}
		return nil
	}

}

// TODONE
func walkFn2(w http.ResponseWriter, query string) filepath.WalkFunc {
	r := regexp.MustCompile(query)

	var fileCount int64 = 0

	return func(path string, f os.FileInfo, err error) error {
		if r.MatchString(path) {
			var tfile FileInfo
			dir, filename := filepath.Split(path)
			tfile.Filename = string(filename)
			tfile.Location = string(dir)

			// Go 1.18 for the win.
			if !slices.Contains(Files, tfile) {
				Files = append(Files, tfile)
				fileCount++
			}

			if w != nil && len(Files) > 0 {

				// Added a basic counter.
				w.Write([]byte(`"` + (strconv.FormatInt(fileCount, 10)) + `":  `))
				fileCount++
				json.NewEncoder(w).Encode(tfile)
				w.Write([]byte(`,`))

			}

			log.Printf("[+] HIT: %s\n", path)

		}
		return nil
	}
}

//==========================================================================\\

func APISTATUS(w http.ResponseWriter, r *http.Request) {

	log.Printf("Entering %s end point", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{ "status" : "API is up and running ",`))
	var regexstrings []string

	for _, regex := range regexes {
		regexstrings = append(regexstrings, regex.String())
	}

	w.Write([]byte(` "regexs" :`))
	json.NewEncoder(w).Encode(regexstrings)
	w.Write([]byte(`}`))
	log.Println(regexes)

}

func MainPage(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering %s end point", r.URL.Path)
	w.Header().Set("Content-Type", "text/html")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "<html><body><H1>Welcome to my awesome File scraping API main page!</H1></body>")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "This API does some stuff and then prints that stuff to the user. Neat")

}

func FindFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering %s end point", r.URL.Path)
	q, ok := r.URL.Query()["q"]

	w.WriteHeader(http.StatusOK)
	if ok && len(q[0]) > 0 {
		log.Printf("Entering search with query=%s", q[0])

		// ADVANCED: Create a function in scrape.go that returns a list of file locations; call and use the result here
		// e.g., func finder(query string) []string { ... }

		var FOUND = false
		for _, File := range Files {
			if File.Filename == q[0] {
				json.NewEncoder(w).Encode(File.Location)
				FOUND = true
			}
		}
		// Added a FOUND boolean, check it here and let the user know if the query failed.
		if !FOUND {
			w.Write([]byte(`{ "There is no file matching your query.",`))
		}
	} else {
		// didn't pass in a search term, show all that you've found
		w.Write([]byte(`"files":`))
		json.NewEncoder(w).Encode(Files)
	}
}

func IndexFiles(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering %s end point", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")

	location, locOK := r.URL.Query()["location"]

	// Restrict where the user can search from (I'm on Ubuntu so we'll start from /home/)
	const rootDir string = "/home/"
	var queryPath string = rootDir + location[0]

	if locOK && len(queryPath) > 0 {
		w.WriteHeader(http.StatusOK)

	} else {
		w.WriteHeader(http.StatusFailedDependency)
		w.Write([]byte(`{ "parameters" : {"required": "location",`))
		w.Write([]byte(`"optional": "regex"},`))
		w.Write([]byte(`"examples" : { "required": "/indexer?location=/xyz",`))
		w.Write([]byte(`"optional": "/indexer?location=/xyz&regex=(i?).md"}}`))
		return
	}

	//wrapper to make "nice json"
	w.Write([]byte(`{ `))

	// Looks
	regex, regexOK := r.URL.Query()["regex"]
	if regexOK {
		filepath.Walk(queryPath, walkFn2(w, `(i?)`+regex[0]))
	} else if err := filepath.Walk(queryPath, walkFn(w)); err != nil {
		log.Panicln(err)
	}

	//wrapper to make "nice json"
	w.Write([]byte(` "status": "completed"} `))

}

//TODO_12 create endpoint that calls resetRegEx AND *** clears the current Files found; ***
//TODO_12 Make sure to connect the name of your function back to the reset endpoint main.go!

//TODO_13 create endpoint that calls clearRegEx ;
//TODO_12 Make sure to connect the name of your function back to the clear endpoint main.go!

//TODO_14 create endpoint that calls addRegEx ;
//TODO_12 Make sure to connect the name of your function back to the addsearch endpoint in main.go!
// consider using the mux feature
// params := mux.Vars(r)
// params["regex"] should contain your string that you pass to addRegEx
// If you try to pass in (?i) on the command line you'll likely encounter issues
// Suggestion : prepend (?i) to the search query in this endpoint

func ExtendSearch(w http.ResponseWriter, r *http.Request) {

}

func ClearSearch(w http.ResponseWriter, r *http.Request) {

}

func ResetSearch(w http.ResponseWriter, r *http.Request) {

}
