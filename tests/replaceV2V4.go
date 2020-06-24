package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	start := time.Now()
	url := "http://localhost:8000/addVertex"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open("vertex.json")
	defer file.Close()
	part1, errFile1 := writer.CreateFormFile(
		"vertex",
		filepath.Base("vertex.json"),
	)
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		fmt.Println(errFile1)
	}
	file, errFile2 := os.Open("vertex4.zip")
	defer file.Close()
	part2, errFile2 := writer.CreateFormFile(
		"vertex4", filepath.Base(
		"vertex4.zip"),
	)
	_, errFile2 = io.Copy(part2, file)
	if errFile2 != nil {
		fmt.Println(errFile2)
	}
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	_, err = client.Do(req)

	fmt.Println("Created Vertex4 in",time.Since(start))
	time.Sleep(60*time.Second)
	start = time.Now()

	url = "http://localhost:8000/removeVertex/vertex2"
	method = "GET"
	client = &http.Client {
	}
	req, err = http.NewRequest(method, url, nil)
	if err != nil {
	  fmt.Println(err)
	}
	_, err = client.Do(req)
	fmt.Println("Removed Vertex2 in",time.Since(start))


}
