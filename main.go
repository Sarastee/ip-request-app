package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type RequestInfo struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
	Message     string  `json:"message"`
}

func main() {
	data := ReadFile()

	f, errCreate := os.Create("out.txt")
	if errCreate != nil {
		log.Fatal(errCreate)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	dataArray := make([]RequestInfo, 0, 1)
	var temp RequestInfo

	for i := 0; i < len(data); i++ {
		newData := IPRequest(data[i])

		err := json.Unmarshal(newData, &temp)
		if err != nil {
			log.Fatal(err)
		}

		dataArray = append(dataArray, temp)
	}
	WriteInFile(dataArray, f)
	fmt.Println("Done.")
}

func IPRequest(IPnum string) []byte {
	ApiURl := fmt.Sprintf("http://ip-api.com/json/%s", IPnum)
	response, err := http.Get(ApiURl)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(response.Body)

	bytes, errRead := io.ReadAll(response.Body)
	if errRead != nil {
		log.Fatal(errRead)
	}

	return bytes
}

func WriteInFile(data []RequestInfo, f *os.File) {
	file, _ := json.MarshalIndent(data, "", " ")
	_, errWrite := f.WriteString(string(file))
	if errWrite != nil {
		log.Fatal(errWrite)
	}
}

func ReadFile() []string {
	var data []string
	file, err := os.Open("in.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)

	i := 0
	for scanner.Scan() {
		data = append(data, scanner.Text())
		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return data
}
