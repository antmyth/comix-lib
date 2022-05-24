package cbz

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"strings"

	"github.com/antmyth/comix-lib/view"
)

const (
	cbrSuffix = ".cbr"
	cbzSuffix = ".cbz"
)

type CBZ struct{}

func (cbz CBZ) BuildIssueFromCBZ(fname, parent string) *view.Issue {
	if !strings.HasPrefix(fname, ".") && strings.HasSuffix(fname, cbzSuffix) {
		fullPath := fmt.Sprintf("%s/%s", parent, fname)
		log.Printf("Reading: %s\n", fullPath)
		issue := readCBZFileData(fullPath)
		issue.SeriesLocation = parent
		return &issue
	}
	return nil
}

func readCBZFileData(ifn string) view.Issue {
	read, err := zip.OpenReader(ifn)
	if err != nil {
		msg := "Failed to open: %s"
		log.Fatalf(msg, err)
	}
	defer read.Close()
	log.Printf("Reading : %s \n", ifn)

	var ci ComicInfo
	for _, file := range read.File {
		if file.FileHeader.Name == "ComicInfo.xml" {
			str, err := readStringFile(file)
			if err != nil {
				log.Fatalf("Failed to read %s from zip: %s", ifn, err)
			}
			if err := xml.Unmarshal([]byte(str), &ci); err != nil {
				panic(err)
			}
		}
	}
	res := ci.ToIssueDB()
	res.Location = ifn
	log.Printf("Read : %+v \n", res)

	return res
}

func readStringFile(file *zip.File) (string, error) {
	res := ""
	fileread, err := file.Open()
	if err != nil {
		msg := "Failed to open zip %s for reading: %s"
		return res, fmt.Errorf(msg, file.Name, err)
	}
	defer fileread.Close()
	var buffer bytes.Buffer
	for {
		readdata := make([]byte, 1024)
		n, err := fileread.Read(readdata)
		if n < 1024 {
			if n > 0 {
				buffer.Write(readdata)
			}
			break
		} else {
			buffer.Write(readdata)
		}
		if err != nil {
			panic(err)
		}
	}
	res = string(buffer.Bytes())

	if err != nil {
		msg := "Failed to read zip %s for reading: %s"
		return res, fmt.Errorf(msg, file.Name, err)
	}

	return res, nil
}
