package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/TV2-Bachelorproject/fetcher/loader"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {

	//connect to database
	db, err := gorm.Open("postgres", "host=127.0.0.1 port=5432 user=root dbname=root password=root sslmode=disable")

	if err != nil {
		panic(err)
	}
	//if error close connection
	defer db.Close()

	//Migrate tables from models
	db.AutoMigrate(&loader.Production{}, &loader.Program{})

	// //Add foreignKeys
	//db.Model(&Program{}).AddForeignKey("production_id", "productions(id)", "RESTRICT", "RESTRICT")

	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "Missing argument")
		os.Exit(1)
	}

	argument := os.Args[1]

	if strings.Contains(argument, "https") {
		fmt.Println(loader.FetchURL(argument, "GET"))
	} else {
		programs, err := loader.FetchFile(argument)

		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}

		for _, program := range programs {
			program.AirtimeFrom = program.Airtime.From
			program.AirtimeTo = program.Airtime.To
			db.Create(&program)
		}
	}
}
