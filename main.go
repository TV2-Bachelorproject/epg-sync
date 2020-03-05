package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/TV2-Bachelorproject/fetcher/fetcher"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "Missing argument")
		os.Exit(1)
	}

	argument := os.Args[1]

	if strings.Contains(argument, "https") {
		fmt.Println(fetcher.FetchURL(argument, "GET"))
	} else {
		data, err := fetcher.FetchFile(argument)

		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Println(data)
	}
}
