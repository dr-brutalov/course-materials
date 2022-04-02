package scrape

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"scrape/logger"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/exp/slices"
)

//==========================================================================\\

// Helper function walk function, modified from Chap 7 BHG to enable passing in of
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

	fmt.Fprintf(w, `<html>
        <body>
            <H1>Welcome to my awesome File scraping API main page!</H1>
            <p>
                This API can be used to examine a filesystem. Neat! <br>
                The following endpoints are can be reached: <br>
                - /api-status ---> Check status of the API. '"ready": true' indicates the API is running properly. <br>
                - /indexer ---> Indexes files according to the active regular expressions. 
								As an option, pass in a query string with a regex to only search that one. <br>
                - /search ---> Pass in a query string to search for a specific file name. <br>
                - /addsearch/{regex} ---> Add the URL encoded regular expression to the searches. <br> 
                - /clear ---> Clear all stored regex values. <br>
                - /reset ---> Reset to the default regex values.
            </p>
        </body>`)

}

func FindFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering %s end point", r.URL.Path)
	q, ok := r.URL.Query()["q"]

	w.WriteHeader(http.StatusOK)
	if ok && len(q[0]) > 0 {

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
		// Did pass in a search term, show all that you've found
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

	regex, regexOK := r.URL.Query()["regex"]
	if regexOK {
		filepath.Walk(queryPath, walkFn2(w, `(i?)`+regex[0]))
	} else if err := filepath.Walk(queryPath, walkFn(w)); err != nil {
		log.Panicln(err)
	}

	//wrapper to make "nice json"
	w.Write([]byte(` "status": "completed"} `))

}

func AddSearch(w http.ResponseWriter, r *http.Request) {
	logger.Log("Entering"+r.URL.Path+" end point", 1)
	w.Header().Set("Content-Type", "application/json")
	regex := mux.Vars(r)["regex"]
	if regex != "" {
		addRegEx("(?i)" + regex)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}

func ClearSearch(w http.ResponseWriter, r *http.Request) {
	logger.Log("Entering"+r.URL.Path+" end point", 1)
	w.Header().Set("Content-Type", "application/json")
	clearRegEx()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}

func ResetSearch(w http.ResponseWriter, r *http.Request) {
	logger.Log("Entering"+r.URL.Path+" end point", 1)
	w.Header().Set("Content-Type", "application/json")
	resetRegEx()
	Files = make([]FileInfo, 0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status: "ok"}`))
}
