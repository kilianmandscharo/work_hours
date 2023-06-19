package main

import (
	"log"
)

func main() {
	db, err := newDatabase()
	if err != nil {
		log.Fatal("ERROR: could not open database", err)
	}
	defer db.close()
	err = db.init()
	if err != nil {
		log.Fatal("ERROR: could not initialize database", err)
	}

	router := newRouter(db)
	router.Run(":8080")
}
