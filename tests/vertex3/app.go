package main

import (
	"github.com/hamzam15comp/vertex"
	"log"
	"os"
	"fmt"
)

var logger *log.Logger

func logInit() {
	f, err := os.OpenFile("apperrors.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	logger = log.New(f, "[INFO]", log.LstdFlags)
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
		logger.Println(datatype, string(data))
                d := string(data)
                d = d + "\n Vertex3 says Hello!"
                data = []byte(d)
		logger.Println(datatype, string(data))
		fmt.Println(string(data))
		fmt.Println("----------------------------------")

	}
}
