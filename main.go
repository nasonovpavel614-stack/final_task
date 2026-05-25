package main

import (
	"log"
	"os"

	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
)

func main() {
	if err := db.Init(db.File()); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := server.Run("web"); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
