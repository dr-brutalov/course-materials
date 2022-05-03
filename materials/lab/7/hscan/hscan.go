package hscan

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

//==========================================================================\\

var shalookup map[string]string
var md5lookup map[string]string

func GuessSingle(sourceHash string, filename string) {

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		password := scanner.Text()

		if len(sourceHash) == 32 {
			// MD5 is a 32-bit hex encoding
			hash := fmt.Sprintf("%x", md5.Sum([]byte(password)))
			if hash == sourceHash {
				fmt.Printf("[+] Password found (MD5): %s\n", password)
			}
		} else if len(sourceHash) == 64 {
			hash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
			if hash == sourceHash {
				fmt.Printf("[+] Password found (SHA-256): %s\n", password)
			}
		} else {
			fmt.Println("Unexpected hash length.")
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}

func GenHashMaps(filename string) {

	shalookup = make(map[string]string)
	md5lookup = make(map[string]string)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	fmt.Println("Creating some channels")

	// Create the structure for workers

	numberOfWorkers := 400
	shaPass := make(chan string, numberOfWorkers)
	md5Pass := make(chan string, numberOfWorkers)
	shaHash := make(chan string, numberOfWorkers)
	md5Hash := make(chan string, numberOfWorkers)

	fmt.Println("The channels are set up and ready to be used. Spawn the workers!")

	// Start the goroutines for the workers
	for i := 0; i < numberOfWorkers; i++ {
		// Call the helper functions for hashing
		go shaWorker(shaPass, shaHash)
		go md5Worker(md5Pass, md5Hash)
	}

	fmt.Println("The channels are ready to go, let's spawn some workers!")
	// Build out aggregators for each of the hashing algorithms

	// SHA256
	go func(result chan string) {
		workingWorkers := numberOfWorkers
		for workingWorkers > 0 {
			res := <-result
			if res == "I have finished my hasing!" {
				numberOfWorkers--
			} else {
				splitString := strings.Split(res, ",")
				hash, pass := splitString[0], splitString[1]
				shalookup[hash] = pass
			}
		}
	}(shaHash)

	// MD5
	go func(result chan string) {
		workingWorkers := numberOfWorkers
		for workingWorkers > 0 {
			res := <-result
			if res == "I have finished my hasing!" {
				numberOfWorkers--
			} else {
				splitString := strings.Split(res, ",")
				hash, pass := splitString[0], splitString[1]
				md5lookup[hash] = pass
			}
		}
	}(md5Hash)

	// Splits up the work among the workers
	for scanner.Scan() {
		password := scanner.Text()

		// Pass the value in the channel to a reference.
		shaPass <- password
		md5Pass <- password
	}

	fmt.Println("Finished designating passwords to workers.")

	shaPass <- "Time to stop SHA256 hashing!"
	md5Pass <- "Time to stop MD5 hashing!"
}

func shaWorker(in, out chan string) {
	for {
		password := <-in

		if password == "Time to stop SHA256 hashing!" {
			out <- "I have finished my hashing!"
			return
		}

		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
		out <- fmt.Sprintf("%s, %s", hash, password)
	}
}

func md5Worker(in, out chan string) {
	for {
		password := <-in

		if password == "Time to stop MD5 hashing!" {
			out <- "I have finished my hashing!"
			return
		}

		hash := fmt.Sprintf("%x", md5.Sum([]byte(password)))
		out <- fmt.Sprintf("%s, %s", hash, password)
	}
}

func GetSHA(hash string) (string, error) {
	password, ok := shalookup[hash]
	if ok {
		return password, nil

	} else {

		return "", errors.New("this password does not exist")

	}
}

func GetMD5(hash string) (string, error) {
	password, ok := md5lookup[hash]
	if ok {
		return password, nil
	} else {
		return "", errors.New("this password does not exist")
	}
}
