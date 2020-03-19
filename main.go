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

		creatingSupportingData(programs)

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
			p.Serie.Title = program.Title
			p.SeasonEpisodeNumber = program.SeasonEpisodeNumber
			p.LinearEpisodeNumber = program.LinearEpisodeNumber
			p.AirtimeFrom = program.Airtime.From
			p.AirtimeTo = program.Airtime.To
			p.Season.SerieID = p.Serie.ID

			setupRelations(program, p)
			db.Create(p)
		}
	}
}

func creatingSupportingData(programs []loader.Program) {
	for _, program := range programs {
		c := public.Category{Name: program.Category}
		sr := public.Serie{Title: program.Title}
		s := public.Season{Title: program.OriginalTitle,
			RawSeasonID: program.SeasonID,
		}
		p := public.Production{
			Country:     program.Production.Country,
			Year:        program.Production.Year,
			ProducedBy:  program.Production.ProducedBy,
			ProducedFor: program.Production.ProducedFor,
			Editor:      program.Production.Editor,
		}

		db.Where(p).FirstOrCreate(&p)

		if sr.Title != "" {
			db.Where(sr).FirstOrCreate(&sr)
		}
		//The  && removes movies and TV2 intern programs from the season table
		if s.Title != "" && s.RawSeasonID != "" {
			db.Where(s).FirstOrCreate(&s)
		}

		if c.Name != "" {
			db.Where(c).FirstOrCreate(&c)
		}

		for _, genre := range program.Genres {
			g := public.Genre{Name: genre}

			if g.Name != "" {
				db.Where(g).FirstOrCreate(&g)
			}

		}
	}
}

func setupRelations(program loader.Program, p *public.Program) {
	//Setup relations between tables
	db.Where(public.Category{Name: program.Category}).First(&p.Category)
	db.Where("name IN (?)", program.Genres).Find(&p.Genres)
	db.Where(public.Season{RawSeasonID: program.SeasonID}).First(&p.Season)
	db.Where(public.Serie{Title: program.Title}).First(&p.Serie)
	db.Where(public.Production{
		Year:        program.Production.Year,
		Country:     program.Production.Country,
		ProducedBy:  program.Production.ProducedBy,
		ProducedFor: program.Production.ProducedFor,
		Editor:      program.Production.Editor,
	}).First(&p.Production)
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
