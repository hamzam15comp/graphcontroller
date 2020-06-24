package main

import (
	"github.com/hamzam15comp/vertex"
	"log"
	"os"
	"time"
	"strconv"
	"io"
	"net/http"
	"io/ioutil"
)

var logger *log.Logger

func logInit() {
	f, err := os.OpenFile("apperrors.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	logger = log.New(f, "", log.Lmicroseconds | log.LUTC)
}

func getImage(name string)([]byte){
	fullUrlFile := "https://storage.needpix.com/thumbs/human-740259_1280.jpg"
	fileName := "images/" + name + ".jpg"
	file, err := os.Create(fileName)
	if err != nil {
		return []byte{}
	}
	resp, err := http.Get(fullUrlFile)
	if err != nil {
		return []byte{}
	}
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return []byte{}
	}
	file.Close()
	resp.Body.Close()

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return []byte{}
	}
	return data
}


func main() {
	logInit()
	err := vertex.LaunchApp()
	if err != nil {
		logger.Println("Error: ", err)
	}
	logger.Println("Created IN and OUT pipe")
	i := 1
	time.Sleep(120 * time.Second)

	for {
		datatype := "img" + strconv.Itoa(i)
		data := getImage(datatype)
		if len(data) == 0 {
			continue
		}

		logger.Println(
			"$S$",
			len(string(data)+datatype),
			"$",
			datatype,
		)

		err = vertex.WriteData("all", datatype, data)
		if err != nil {
			logger.Println(err)
		}
		i = i + 1
		time.Sleep(60 * time.Second)
	}
}
