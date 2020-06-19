package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var pwd string

var logger *log.Logger

func logInit() {
        f, err := os.OpenFile("ctrlerrors.log",
                os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
                log.Println(err)
        }

        logger = log.New(f, "[INFO]", log.LstdFlags)
}


func parseVertex(s string) (letters, numbers string) {
    var l, n []rune
    for _, r := range s {
        switch {
        case r >= 'A' && r <= 'Z':
            l = append(l, r)
        case r >= 'a' && r <= 'z':
            l = append(l, r)
        case r >= '0' && r <= '9':
            n = append(n, r)
        }
    }
    return string(l), string(n)
}

func parseMultiPart(w http.ResponseWriter, r *http.Request) {
	parseErr := r.ParseMultipartForm(32 << 20)
	if parseErr != nil {
		http.Error(
			w,
			"failed to parse multipart message",
			http.StatusBadRequest,
		)
		return
	}
	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		http.Error(
			w,
			"expecting multipart form file",
			http.StatusBadRequest,
		)
		return
	}
}

func getFile(r *http.Request, filekey string) error {
	file, data, err := r.FormFile(filekey)
	if err != nil {
		return fmt.Errorf("failed to get file", err)
	}
	defer file.Close()

	tmpfile, _ := os.Create(pwd + data.Filename)
	defer tmpfile.Close()

	io.Copy(tmpfile, file)
	logger.Println("Fetched file:", data.Filename)

	//Unzip if .zip file
	if filepath.Ext(data.Filename) == ".zip" {
		_, err := Unzip(pwd+data.Filename, pwd)
		if err != nil {
			return fmt.Errorf(
				"Failed to extract %s",
				data.Filename,
			)
		}
		logger.Println("Extracted", data.Filename)
	}
	return nil
}

func Unzip(src string, dest string) ([]string, error) {
	var filenames []string
	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()
	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}
		filenames = append(filenames, fpath)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}
		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

func FindInSlice(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

/*
func main(){
	fmt.Println("Pouring rain")
}
*/
