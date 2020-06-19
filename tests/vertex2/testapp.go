package main

import (
	"log"
	"time"

	"github.com/hamzam15comp/vertex"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	log.Printf("Initializing Vertex V11 Publisher")
	pub1 := vertex.InitVertex(1, 1, "pub")
	for {
		log.Printf("Sending data to B1")
		//		err := vertex.SendData(0, "ALL", []byte("dataALL"))
		//		failOnError(err, "Sending Data Failed")
		err := vertex.SendData(pub1, "For2", []byte("data2"))
		failOnError(err, "Sending Data Failed")
		//		err = vertex.SendData(3, "For3", []byte("data3"))
		//		failOnError(err, "Sending Data Failed")
		//		err = vertex.SendData(4, "For4", []byte("data4"))
		//		failOnError(err, "Sending Data Failed")
		time.Sleep(2 * time.Second)
	}
}
