package comicvine

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/antmyth/comix-lib/view"
)

const apiKey = "abef181d68b7a432d1438bbdddada81849521d24"
const format = "format=json"
const baseURL = "http://comicvine.gamespot.com/api"

type ComicVine struct{}

func (cv ComicVine) GetIssue(id string) IssueData {
	resp := executeRequest("issue", id)
	info := IssueResponse{}
	json.Unmarshal([]byte(resp), &info)

	return info.Results
}

func (cv ComicVine) GetVolume(id string) VolumeData {
	resp := executeRequest("volume", id)
	info := VolumeResponse{}
	json.Unmarshal([]byte(resp), &info)

	return info.Results
}

func buildRequest(resource, key string) string {
	return fmt.Sprintf("%s/%s/%s/?%s&api_key=%s", baseURL, resource, key, format, apiKey)
}

func executeRequest(resource, id string) string {
	req := buildRequest(resource, id)
	// log.Println(req)
	resp, err := http.Get(req)

	// log.Println(err)
	// log.Println(resp)

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
	}

	return bodyString
}

func (cvImg Image) FromComicVine() view.Image {
	return view.Image{
		SmallUrl:    cvImg.SmallUrl,
		ThumbUrl:    cvImg.ThumbUrl,
		TinyUrl:     cvImg.TinyUrl,
		OriginalUrl: cvImg.OriginalUrl,
	}
}
