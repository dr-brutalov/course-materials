// Build and Use this File to interact with the shodan package
// In this directory lab/3/shodan/main:
// go build main.go
// SHODAN_API_KEY=YOURAPIKEYHERE ./main <search term>

package main

import (
	"fmt"
	"log"
	"os"
	"shodan/shodan"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage: main <searchterm>")
	}
	apiKey := os.Getenv("SHODAN_API_KEY")
	s := shodan.New(apiKey)
	info, err := s.APIInfo()
	if err != nil {
		log.Panicln(err)
	}

	var nextPage string
	fmt.Println("Press Y and Enter Key to get the next page.")
	fmt.Scanln(&nextPage)
	pageIndex := 0

	for nextPage == "Y" {
		pageIndex++
		fmt.Printf(
			"Query Credits: %d\nScan Credits:  %d\n\n",
			info.QueryCredits,
			info.ScanCredits)

		hostSearch, err := s.HostSearch(os.Args[1], pageIndex)
		if err != nil {
			log.Panicln(err)
		}

		/*
			fmt.Printf("Host Data Dump\n")
			for _, host := range hostSearch.Matches {
				fmt.Println("==== start ", host.IPString, "====")
				h, _ := json.Marshal(host)
				fmt.Println(string(h))
				fmt.Println("==== end ", host.IPString, "====")
				//fmt.Println("Press the Enter Key to continue.")
				//fmt.Scanln()
			}*/

		fmt.Printf("IP, Port, AreaCode\n")

		for _, host := range hostSearch.Matches {
			fmt.Printf("%s, %d, %f Lat, %f Lon\n", host.IPString, host.Port, host.Location.Latitude, host.Location.Longitude)
		}
		fmt.Println("Press Y and Enter to get the next page.")
		fmt.Scanln(&nextPage)
	}
}
