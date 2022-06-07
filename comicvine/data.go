package comicvine

type IssueResponse struct {
	Results IssueData `json:"results,omitempty"`
}

type VolumeResponse struct {
	Results VolumeData `json:"results,omitempty"`
}

type PublisherResponse struct {
	Results PublisherData `json:"results,omitempty"`
}

type IssueData struct {
	ID            int         `json:"id"`
	ApiDetailUrl  string      `json:"api_detail_url"`
	SiteDetailUrl string      `json:"site_detail_url"`
	Volume        VolumeShort `json:"volume"`
	Image         Image       `json:"image"`
}

type VolumeData struct {
	ID            int            `json:"id"`
	ApiDetailUrl  string         `json:"api_detail_url"`
	SiteDetailUrl string         `json:"site_detail_url"`
	CountOfIssues int            `json:"count_of_issues"`
	Image         Image          `json:"image"`
	Description   string         `json:"description,omitempty"`
	Publisher     PublisherShort `json:"publisher"`
}

type VolumeShort struct {
	ID            string `json:"id"`
	ApiDetailUrl  string `json:"api_detail_url"`
	SiteDetailUrl string `json:"site_detail_url"`
	Name          string `json:"name"`
}

type Image struct {
	SmallUrl    string `json:"small_url,omitempty"`
	ThumbUrl    string `json:"thumb_url,omitempty"`
	TinyUrl     string `json:"tiny_url,omitempty"`
	OriginalUrl string `json:"original_url,omitempty"`
}

type PublisherShort struct {
	ID           int    `json:"id"`
	ApiDetailUrl string `json:"api_detail_url"`
	Name         string `json:"name"`
}

type PublisherData struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Image       Image  `json:"image"`
}
