package accounts

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
)

func GetCsv(file string) [][]string {
	f, _ := os.Open(file)

	r := csv.NewReader(bufio.NewReader(f))

	result, _ := r.ReadAll()

	log.Println(result)

	return result

	// for {
	// 	record, err := r.Read()
	// 	// Stop at EOF.
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	// Display record.
	// 	// ... Display record length.
	// 	// ... Display all individual elements of the slice.
	// 	fmt.Println(record)
	// 	fmt.Println(len(record))
	// 	for value := range record {
	// 		fmt.Printf("  %v\n", record[value])
	// 	}
	// }
}
