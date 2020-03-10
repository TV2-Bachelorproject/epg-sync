package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/TV2-Bachelorproject/fetcher/loader"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
)

//Program struct
type Program struct {
	gorm.Model
	ProgramID           string `json:"programId" gorm:"primary_key"`
	Title               string
	OriginalTitle       string
	Teaser              string
	Description         string
	CastRaw             string `json:"cast"`
	Category            string
	Genres              pq.StringArray `gorm:"type:varchar(100)[]"`
	SeasonID            string         `gorm: "type:varchar(100)"`
	SeasonEpisodeNumber string
	Production          Production `gorm:"foreignkey:ID"`
	//Airtime             Airtime    `gorm:"foreignkey:From"`
}

type Airtime struct {
	gorm.Model
	From int
	To   int
}

type Production struct {
	gorm.Model
	Country     string
	Year        int
	ProducedBy  string
	ProducedFor string
	Editor      string
}

func main() {

	//connect to database
	db, err := gorm.Open("postgres", "host=127.0.0.1 port=5432 user=root dbname=root password=root sslmode=disable")

	if err != nil {
		panic(err)
	}
	//if error close connection
	defer db.Close()

	program := Program{
		ProgramID:     "JPVSYOQcfsrv",
		Title:         "Badehotellet",
		OriginalTitle: "Badehotellet - år 7",
		Teaser:        "Dansk tv-serie - år 7 - (4:6) ",
		Description:   "Da udenrigsminister Scavenius udtrykker beundring for de store tyske sejre, mener Lydia Vetterstrøm, at tiden er kommet til at vise det nationale sindelag ved at tage til Alsang. Det er dog ikke helt nemt at få alle med. Fru Frigh er kun optaget af Johan Ramsing, efter de i al hemmelighed har tilbragt natten sammen. Weyse øver sin filmrolle til 'Hærværk' med anti-nazisten Gerhard Flügelhorn, samtidig med at han febrilsk forsøger at skjule, han har optrådt for tyskerne i Skagen. Og Madsen kæmper med at få konen med til banketten hos Værnemagten i Aalborg, mens Morten forsøger at finde ud af, hvem der hjælper tyskerne med at udbygge lufthavnen dernede.",
		CastRaw:       "Medvirkende: Amanda: Amalie Dollerup, Molly: Bodil Jørgensen og Georg Madsen: Lars Ranthe. Desuden medvirker: Therese Madsen: Anne Louise Hassing. Edward Weyse: Jens Jacob Tychsen. Helene Weyse: Cecilie Stenspil. Alice Frigh: Anette Støvelbæk. Fru Fjeldsø: Birthe Neumann. Lydia Vetterstrøm Ploug: Sonja Oppenhagen. Gerhard Flügelhorn: Thure Lindhardt. Oberst Konrad Fuchs: Joachim Kappl. Hjalmar Aurland: Peter Hesse Overgaard. Philip Dupont: Sigurd Holmen Le Dous. Morten: Morten Aaskov Hemmingsen. Grev Ditmar: Mads Wille. Otilia: Merete Mærkedahl. Edith: Ulla Vejby. Ane: Mia Helene Højgaard. Nana: Laura Kronborg Kjær. Peter Andreas: Kristian Halken. Bertha: Lucia Vinde Dirchsen. Leslie: Lue Støvelbæk. August Molin: Søren Sætter-Lassen. Johan Ramsing: Rasmus Botoft. Folmer Gregersen: Troels Malling Thaarup. Jan Larsen: Aske Bang. Uwe Kiessling: Anton Rubtsov.  Original episodetitel: Alsang og hærværk.  Forfattere: Stig Thorsboe og Hanna Lundblad. Instruktion: Fabian Wullenweber. Producer: Michael Bille Frandsen. ",
		Category:      "TV-Serie",
		Genres: []string{
			"Serie",
			"Drama \u0026 Fiktion",
		},
		SeasonID:            "JPVNyQrqeCTv",
		SeasonEpisodeNumber: "",
		Production: Production{
			Country:     "Danmark",
			Year:        2020,
			ProducedBy:  "SF Film Production i co-produktion med Nitrat Film og Thorsboe \u0026 Lundblad",
			ProducedFor: "TV 2 DANMARK",
			Editor:      "Pernille Bech Christensen",
		},
		// Airtime: Airtime{
		// 	From: 1582457700000,
		// 	To:   1582461000000,
		// },
	}

	//db.AutoMigrate(&Production{}, &Program{})
	//db.CreateTable(&Program{})
	db.Create(&program)
	db.Save(&program)

	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "Missing argument")
		os.Exit(1)
	}

	argument := os.Args[1]

	if strings.Contains(argument, "https") {
		fmt.Println(loader.FetchURL(argument, "GET"))
	} else {
		data, err := loader.FetchFile(argument)

		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Println(data)
	}
}
