package viewmodel

import "fmt"

func New() *Image {
	res := Image{}
	return &res
}

type ComicsLib struct {
	SeriesCount int      `json:"seriesCount"`
	SeriesList  []Series `json:"seriesList"`
}

type Series struct {
	ID          int     `json:"id,omitempty"`
	VineId      string  `json:"vineId,omitempty"`
	Series      string  `json:"series"`
	Volume      string  `json:"volume,omitempty"`
	Publisher   string  `json:"publisher,omitempty"`
	Count       int     `json:"count"`
	TotalCount  int     `json:"totalcount"`
	Web         string  `json:"web,omitempty"`
	Images      Image   `json:"image,omitempty"`
	Issues      []Issue `json:"issues"`
	Location    string  `json:"location,omitempty"`
	Description string  `json:"description,omitempty"`
}

type Issue struct {
	ID             int    `json:"id,omitempty"`
	Title          string `json:"title"`
	Series         string `json:"series"`
	SeriesId       int    `json:"seriesId"`
	Number         string `json:"number"`
	Volume         string `json:"volume,omitempty"`
	Publisher      string `json:"publisher,omitempty"`
	Web            string `json:"web,omitempty"`
	VolumeAPI      string `json:"volume_api,omitempty"`
	Images         Image  `json:"image,omitempty"`
	SeriesLocation string `json:"seriesLocation,omitempty"`
	Location       string `json:"location,omitempty"`
}

type Image struct {
	SmallUrl    string `json:"small_url,omitempty"`
	ThumbUrl    string `json:"thumb_url,omitempty"`
	TinyUrl     string `json:"tiny_url,omitempty"`
	OriginalUrl string `json:"original_url,omitempty"`
}

func (issue Issue) ToString() string {
	return fmt.Sprintf("%s|%s", issue.Series, issue.Volume)
}

func AsSeriesMap(issues []Issue) map[string][]Issue {
	res := make(map[string][]Issue, 0)
	for _, v := range issues {
		ss, e := res[v.ToString()]
		if e {
			ss = append(ss, v)
			res[v.ToString()] = ss
		} else {
			res[v.ToString()] = []Issue{v}
		}
	}
	return res
}
