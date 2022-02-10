package main

import "bhg-scanner/scanner"

func main() {
	scanner.PortScanner([2]int{1, 1024}, 100)
}
