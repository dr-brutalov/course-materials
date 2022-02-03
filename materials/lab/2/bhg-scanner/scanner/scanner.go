// bhg-scanner/scanner.go modified from Black Hat Go > CH2 > tcp-scanner-final > main.go
// Code : https://github.com/blackhat-go/bhg/blob/c27347f6f9019c8911547d6fc912aa1171e6c362/ch-2/tcp-scanner-final/main.go
// License: {$RepoRoot}/materials/BHG-LICENSE
// Useage:
// A general use port scanner, yo.

package scanner

import (
	"fmt"
	"net"
	"sort"
	"time"
)

// Set a duration for the DialTimeout function (1 second seems fine for now)
var dur = 1 * time.Second

func worker(ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("scanme.nmap.org:%d", p)
		conn, err := net.DialTimeout("tcp", address, dur)
		if err != nil {
			results <- -p
			continue
		}
		conn.Close()
		results <- p
	}
}

// for Part 5 - consider
// easy: taking in a variable for the ports to scan (int? slice? ); a target address (string?)?
// med: easy + return  complex data structure(s?) (maps or slices) containing the ports.
// hard: restructuring code - consider modification to class/object
// No matter what you do, modify scanner_test.go to align; note the single test currently fails
// Returns the number of open ports and the total number of ports scanned.
func PortScanner(numPorts int) (int, int) {

	var openports []int   // notice the capitalization here. access limited!
	var closedports []int // var for tracking the closed ports

	ports := make(chan int, numPorts)
	results := make(chan int)

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	go func() {
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
	}()

	for i := 0; i < 1024; i++ {
		port := <-results
		if port > 0 {
			openports = append(openports, port)
		} else {
			negPort := port * (-1)
			closedports = append(closedports, negPort)
		}
	}

	close(ports)
	close(results)
	sort.Ints(openports)
	sort.Ints(closedports)

	//TODO 5 : Enhance the output for easier consumption, include closed ports

	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
	//portList := append(openports, closedports...)

	//writeToCSV(portList)

	totalPorts := len(openports) + len(closedports) // TODO 6 : Return total number of ports scanned (number open, number closed);
	numOpenPorts := len(openports)
	//response := fmt.Sprintf("We successfully scanned %d ports!", totalPorts)

	return numOpenPorts, totalPorts //response
	//you'll have to modify the function parameter list in the defintion and the values in the scanner_test
}

/*
// A helper function to reduce repetitive error code logic (thanks to Andey for the inspiration
// and https://golangcode.com/write-data-to-a-csv-file/ for a simple extension (previous version panics w/out a msg.)).
func checkErr(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

// A function for writing to a .csv file, breaking this piece out to help keep the PortScanner function readable
// I used this source as a reference for the function: https://golangcode.com/write-data-to-a-csv-file/

func writeToCSV(fileName string, portList []int) {
	file, err := os.Create("fileName" + ".csv")
	checkErr("Cannot create file, yo.", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range portList {
		err := writer.Write(value)
		checkErr("Cannot write to file, yo.", err)
	}
}
*/
