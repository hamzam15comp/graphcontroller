package main

import (
	"github.com/hamzam15comp/vertex"
	"log"
	"os"
	"time"
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

		datatype := "message"
		data := []byte("Vertex1 says Hello!")
		logger.Println(datatype, string(data))
		time.Sleep(10 * time.Second)
		err = vertex.WriteData("all", datatype, data)
		if err != nil {
			logger.Println(err)
		}
	}
}
