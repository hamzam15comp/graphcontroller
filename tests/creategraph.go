package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func main() {

	url := "http://localhost:8000/createGraph"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open("graph.json")
	defer file.Close()
	part1, errFile1 := writer.CreateFormFile(
		"graph",
		filepath.Base("graph.json"),
	)
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		fmt.Println(errFile1)
	}
	file, errFile2 := os.Open("vertex1.zip")
	defer file.Close()
	part2, errFile2 := writer.CreateFormFile(
		"vertex1",
		filepath.Base("vertex1.zip"),
	)
	_, errFile2 = io.Copy(part2, file)
	if errFile2 != nil {
		fmt.Println(errFile2)
	}
	file, errFile3 := os.Open("vertex2.zip")
	defer file.Close()
	part3, errFile3 := writer.CreateFormFile(
		"vertex2",
		filepath.Base("vertex2.zip"),
	)
	_, errFile3 = io.Copy(part3, file)
	if errFile3 != nil {
		fmt.Println(errFile3)
	}
	file, errFile4 := os.Open("edge1.zip")
	defer file.Close()
	part4, errFile4 := writer.CreateFormFile(
		"edge1",
		filepath.Base("edge1.zip"),
	)
	_, errFile4 = io.Copy(part4, file)
	if errFile4 != nil {
		fmt.Println(errFile4)
	}
	file, errFile5 := os.Open("vertex3.zip")
	defer file.Close()
	part5, errFile5 := writer.CreateFormFile(
		"vertex3",
		filepath.Base("vertex3.zip"),
	)
	_, errFile5 = io.Copy(part5, file)
	if errFile5 != nil {
		fmt.Println(errFile5)
	}
	file, errFile6 := os.Open("edge2.zip")
	defer file.Close()
	part6, errFile6 := writer.CreateFormFile(
		"edge2",
		filepath.Base("edge2.zip"),
	)
	_, errFile6 = io.Copy(part6, file)
	if errFile6 != nil {
		fmt.Println(errFile6)
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
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}
