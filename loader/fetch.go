package loader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

//EPG struct
type EPG struct {
	Schedule []Program
}

//Program struct
type Program struct {
	gorm.Model
	ProgramID           string `json:"programId" gorm:"primary_key"`
	Title               string
	OriginalTitle       string
	EpisodeTitle        string
	Teaser              string
	Description         string
	CastRaw             string `json:"cast"`
	Category            string
	Genres              pq.StringArray `gorm:"type:varchar(100)[]"`
	SeasonID            string         `gorm:"type:varchar(100)"`
	SeasonEpisodeNumber int
	LinearEpisodeNumber int
	ProductionID        uint
	Production          Production `gorm:"foreignkey:production_id"`
	AirtimeFrom         int        `gorm:"type:bigint"`
	AirtimeTo           int        `gorm:"type:bigint"`
	Airtime             struct {
		From int `gorm:"-"`
		To   int `gorm:"-"`
	}
}

type Production struct {
	gorm.Model
	Country     string
	Year        int
	ProducedBy  string
	ProducedFor string
	Editor      string
}

// FetchURL data
func FetchURL(url, method string) string {

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	return string(body)
}

//FetchFile data
func FetchFile(fileName string) ([]Program, error) {
	//Open json file
	jsonFile, err := os.Open(fileName)

	if err != nil {
		return []Program{}, err
	}

	//read data from jsonfile
	byteValue, _ := ioutil.ReadAll(jsonFile)

	data := EPG{}

	json.Unmarshal([]byte(byteValue), &data)

	//close file so we can parse later on
	defer jsonFile.Close()

	return data.Schedule, nil
}

//CreateFile for json
func CreateFile(fileName string, data []byte) {

	_ = ioutil.WriteFile(fileName, data, 0644)

}
