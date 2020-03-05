package loader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

//EPG struct
type EPG struct {
	Schedule []Program
}

//Program struct
type Program struct {
	ID                  string `json:"programId"`
	Title               string
	OriginalTitle       string
	Teaser              string
	Description         string
	CastRaw             string `json:"cast"`
	Category            string
	Genres              []string
	SeasonID            string
	SeasonEpisodeNumber string
	Production          struct {
		Country     string
		Year        int
		ProducedBy  string
		ProducedFor string
		Editor      string
	}
	Airtime struct {
		From int
		To   int
	}
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
func FetchFile(fileName string) (string, error) {
	//Open json file
	jsonFile, err := os.Open(fileName)

	if err != nil {
		return "", err
	}

	//read data from jsonfile
	byteValue, _ := ioutil.ReadAll(jsonFile)

	data := EPG{}

	json.Unmarshal([]byte(byteValue), &data)

	json, err := json.MarshalIndent(data, "", "\t")

	if err != nil {
		return "", err
	}

	//close file so we can parse later on
	defer jsonFile.Close()

	return string(json), nil
}

//CreateFile for json
func CreateFile(fileName string, data []byte) {

	_ = ioutil.WriteFile(fileName, data, 0644)

}
