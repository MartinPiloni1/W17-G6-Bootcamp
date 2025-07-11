package main

import (
	"fmt"
	config "github.com/aaguero_meli/W17-G6-Bootcamp/cmd/db"
	"io/ioutil"
	"log"
)

func main() {
	db := config.MustOpenDB()
	defer db.Close()

	dumpPath := "docs/db/seed/dump.sql"
	data, err := ioutil.ReadFile(dumpPath)
	if err != nil {
		log.Fatalf("Could not read the dump: %v", err)
	}

	fmt.Println("Executing example data dump...")

	_, err = db.Exec(string(data))
	if err != nil {
		log.Fatalf("Error executing the dump: %v", err)
	}

	fmt.Println("Example data load completed successfully!")
}
