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
		log.Fatalf("No se pudo leer el dump: %v", err)
	}

	fmt.Println("Ejecutando dump de datos de ejemplo...")

	_, err = db.Exec(string(data))
	if err != nil {
		log.Fatalf("Error ejecutando el dump: %v", err)
	}

	fmt.Println("¡Carga de datos de ejemplo completada exitosamente!")
}
