package main

import (
	"github.com/hamzam15comp/vertex"
	"log"
	"image"
	"os"
	"io/ioutil"
	"fmt"
	"gocv.io/x/gocv"
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
			continue
		}

		//datatype := "img1"
		//data, err := ioutil.ReadFile("images/img.jpg")
		//if err != nil {
		//	logger.Println("Read Error")
		//	continue
		//}

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

		classifier := gocv.NewCascadeClassifier()
		defer classifier.Close()
		xmlFile := "images/haarcascade_frontalface_default.xml"
		if !classifier.Load(xmlFile) {
			logger.Println(
				"Error reading cascade file: %v\n",
				xmlFile,
			)
			continue
		}
		img := gocv.IMRead(fileName, gocv.IMReadUnchanged)
		if img.Empty(){
			logger.Println("Empty Mat")
			continue
		}
		rects := classifier.DetectMultiScale(img)
		fmt.Printf("Found %d faces in %s\n", len(rects), fileName)
		for _, r := range rects {
		        imgFace := img.Region(r)
		        gocv.GaussianBlur(
		                imgFace,
		                &imgFace,
		                image.Pt(75, 75),
		                0,
		                0,
		                gocv.BorderDefault,
		        )
		        imgFace.Close()
		}
		blurName := "images/blur" + datatype + ".jpg"
		f, err := os.Create(blurName)
		if err != nil {
			logger.Println(err)
			continue
		}
		f.Close()
		gocv.IMWrite(blurName, img)


		datatype = "b" + datatype
		data, err = ioutil.ReadFile(blurName)
		if err != nil {
			logger.Println(err)
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

	}
}
