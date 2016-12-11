package main

import (
	"bufio"
	"log"
	"os"

	"github.com/remotejob/gojobextractor/signupnext/compliteSignUp"
)

func main() {

	if file, err := os.Open("signuplinks.csv"); err == nil {

		// make sure it gets closed
		defer file.Close()

		// create a new scanner and read the file line by line
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			log.Println(scanner.Text())
			link := scanner.Text()
			compliteSignUp.Complite(link)

		}

		// check for errors
		if err = scanner.Err(); err != nil {
			log.Fatal(err)
		}

	} else {
		log.Fatal(err)
	}

}
