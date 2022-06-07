package comicvine

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/antmyth/comix-lib/viewmodel"
)

const apiKey = "abef181d68b7a432d1438bbdddada81849521d24"
const format = "format=json"
const baseURL = "http://comicvine.gamespot.com/api"
const volumePrefix = "4050"
const issuePrefix = "4000"
const publisherPrefix = "4010"

type ComicVine struct{}

func (cv ComicVine) GetIssue(id string) (IssueData, error) {
	resp, err := executeRequest("issue", id)
	if err != nil {
		return IssueData{}, err
	}
	info := IssueResponse{}
	json.Unmarshal([]byte(resp), &info)

	return info.Results, nil
}

func (cv ComicVine) GetIssueBy(id int) (IssueData, error) {
	return cv.GetIssue(fmt.Sprintf("%v-%v", issuePrefix, id))
}

func (cv ComicVine) GetVolume(id string) (VolumeData, error) {
	resp, err := executeRequest("volume", id)
	if err != nil {
		return VolumeData{}, err
	}
	info := VolumeResponse{}
	json.Unmarshal([]byte(resp), &info)

	return info.Results, nil
}

func (cv ComicVine) GetVolumeBy(id int) (VolumeData, error) {
	return cv.GetVolume(fmt.Sprintf("%v-%v", volumePrefix, id))
}

func (cv ComicVine) GetPublisher(id int) (PublisherData, error) {
	resp, err := executeRequest("publisher", fmt.Sprintf("%v-%v", publisherPrefix, id))
	if err != nil {
		return PublisherData{}, err
	}
	info := PublisherResponse{}
	json.Unmarshal([]byte(resp), &info)

	return info.Results, nil
}

func buildRequest(resource, key string) string {
	return fmt.Sprintf("%s/%s/%s/?%s&api_key=%s", baseURL, resource, key, format, apiKey)
}

func executeRequest(resource, id string) (string, error) {
	req := buildRequest(resource, id)
	// log.Printf("comicvine request:%v\n", req)
	resp, err := http.Get(req)

	var bodyString string
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString = string(bodyBytes)
		// log.Printf("comicvine response:{{%+v}}\n", bodyString)
	} else {
		log.Printf("Failed to get data for %v with ID=%v from comicvine on %v\n%v\n", resource, id, req, resp.StatusCode)
		log.Println(resp.Body)
		errStr := fmt.Sprintf("Failed to get data, response not OK actual:%v.", resp.StatusCode)
		return "", errors.New(errStr)
	}
	//throtling access to comic vine API
	time.Sleep(2 * time.Second)
	return bodyString, nil
}

func (cvImg Image) FromComicVine() viewmodel.Image {
	return viewmodel.Image{
		SmallUrl:    cvImg.SmallUrl,
		ThumbUrl:    cvImg.ThumbUrl,
		TinyUrl:     cvImg.TinyUrl,
		OriginalUrl: cvImg.OriginalUrl,
	}
}

func (cvPub PublisherData) FromComicVinePublisher() viewmodel.Publisher {
	res := viewmodel.Publisher{}
	res.ID = cvPub.ID
	res.Name = cvPub.Name
	res.Description = cvPub.Description
	res.Images = cvPub.Image.FromComicVine()
	return res
}
