package main

import (
	"log"

	"github.com/kilianmandscharo/work_hours/database"
	"github.com/kilianmandscharo/work_hours/server"
)

func main() {
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal("ERROR: could not open database", err)
	}
	defer db.Close()
	err = db.Init()
	if err != nil {
		log.Fatal("ERROR: could not initialize database", err)
	}

	router := server.NewRouter(db)
	router.Run(":8080")
}
