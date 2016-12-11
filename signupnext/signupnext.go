package main

import (
	"bufio"
	"log"
	"os"
)

func main() {

	if file, err := os.Open("signuplinks.csv"); err == nil {

		// make sure it gets closed
		defer file.Close()

		// create a new scanner and read the file line by line
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			log.Println(scanner.Text())
		}

		// check for errors
		if err = scanner.Err(); err != nil {
			log.Fatal(err)
		}

	} else {
		log.Fatal(err)
	}

}
