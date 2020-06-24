package main

import (
	"github.com/hamzam15comp/vertex"
	"log"
	"os"
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

func main() {
	logInit()
	err := vertex.LaunchApp()
	if err != nil {
		logger.Println("Error: ", err)
	}
	logger.Println("Created IN and OUT pipe")
	for {
		datatype, data, err := vertex.ReadData()
		if err != nil {
			logger.Println(err)
		}

                logger.Println(
                        "$R$",
                        len(string(data)+datatype),
                        "$",
                        datatype,
                )
		fileName := "images/" + datatype + ".jpg"
		file, err := os.Create(fileName)
		if err != nil {
		        logger.Println(err)
		        continue
		}
		file.Close()
		err = ioutil.WriteFile(fileName, data, 0666)
		if err != nil {
		        logger.Println(err)
		        continue
		}

	}
}
