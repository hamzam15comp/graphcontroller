package main

import (
	"github.com/hamzam15comp/vertex"
	"log"
	"os"
)

var logger *log.Logger

func logInit() {
	f, err := os.OpenFile("launcherrors.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	logger = log.New(f, "[INFO]", log.LstdFlags)
}

func main() {
	logInit()
	logger.Println("Launching Vertex Agent")
	vertex.VertexAgentLaunch()
}
