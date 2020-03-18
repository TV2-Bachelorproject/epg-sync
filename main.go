package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/TV2-Bachelorproject/fetcher/loader"
	"github.com/TV2-Bachelorproject/server/model/public"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

const (
	HOST        = "127.0.0.1"
	PORT        = "5432"
	DB_USER     = "root"
	DB_PASSWORD = "root"
	DB_NAME     = "root"
)

func main() {

	initDB()
	dataToDB()
	closeDB()
}

func dataToDB() {
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
			p := &public.Program{}

			p.ProgramID = program.ProgramID

			if program.EpisodeTitle == "" {
				p.Title = program.Title
			} else {
				p.Title = program.EpisodeTitle
			}

			p.Teaser = program.Teaser
			p.Description = program.Description
			p.Cast = program.CastRaw
			p.Category.Name = program.Category
			p.Serie.Title = program.Title
			p.Season.RawSeasonID = program.SeasonID
			p.Season.Title = program.OriginalTitle
			p.SeasonEpisodeNumber = program.SeasonEpisodeNumber
			p.LinearEpisodeNumber = program.LinearEpisodeNumber
			p.ProductionID = program.ProductionID
			p.Production.Country = program.Production.Country
			p.Production.Year = program.Production.Year
			p.Production.ProducedBy = program.Production.ProducedBy
			p.Production.ProducedFor = program.Production.ProducedFor
			p.Production.Editor = program.Production.Editor
			p.AirtimeFrom = program.Airtime.From
			p.AirtimeTo = program.Airtime.To

			db.Create(p)
		}
	}
}

func initDB() {
	var err error

	//connect to database
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		HOST, PORT, DB_USER, DB_NAME, DB_PASSWORD)
	db, err = gorm.Open("postgres", dbinfo)

	if err != nil {
		panic(err)
	}

	//Migrate tables from models
	//db.AutoMigrate(&loader.Production{}, &loader.Program{})

	// //Add foreignKeys
	//db.Model(&Program{}).AddForeignKey("production_id", "productions(id)", "RESTRICT", "RESTRICT")
}

func closeDB() error {
	return db.Close()
}
