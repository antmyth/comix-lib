package cbz

import (
	"encoding/xml"

	"github.com/antmyth/comix-lib/comicvine"
	"github.com/antmyth/comix-lib/viewmodel"
)

type ComicInfo struct {
	XMLName     xml.Name   `xml:"ComicInfo,omitempty"`
	Text        string     `xml:",chardata"`
	Xsd         string     `xml:"xsd,attr"`
	Xsi         string     `xml:"xsi,attr"`
	Title       string     `xml:"Title,omitempty"`
	Series      string     `xml:"Series,omitempty"`
	Number      string     `xml:"Number,omitempty"`
	Volume      string     `xml:"Volume,omitempty"`
	Summary     string     `xml:"Summary,omitempty"`
	Notes       string     `xml:"Notes,omitempty"`
	Year        string     `xml:"Year,omitempty"`
	Month       string     `xml:"Month,omitempty"`
	Day         string     `xml:"Day,omitempty"`
	Writer      string     `xml:"Writer,omitempty"`
	Penciller   string     `xml:"Penciller,omitempty"`
	Inker       string     `xml:"Inker,omitempty"`
	Colorist    string     `xml:"Colorist,omitempty"`
	Letterer    string     `xml:"Letterer,omitempty"`
	CoverArtist string     `xml:"CoverArtist,omitempty"`
	Editor      string     `xml:"Editor,omitempty"`
	Publisher   string     `xml:"Publisher,omitempty"`
	Genre       string     `xml:"Genre,omitempty"`
	Web         string     `xml:"Web,omitempty"`
	PageCount   string     `xml:"PageCount,omitempty"`
	LanguageISO string     `xml:"LanguageISO,omitempty"`
	Characters  string     `xml:"Characters,omitempty"`
	Teams       string     `xml:"Teams,omitempty"`
	Pages       ComicPages `xml:"Pages,omitempty"`
}
type ComicPages struct {
	Text string      `xml:",chardata"`
	Page []ComicPage `xml:"Page"`
}
type ComicPage struct {
	Text        string `xml:",chardata"`
	Image       string `xml:"Image,attr"`
	ImageWidth  string `xml:"ImageWidth,attr"`
	ImageHeight string `xml:"ImageHeight,attr"`
	Type        string `xml:"Type,attr"`
}

func (ci ComicInfo) ToSeriesDB() viewmodel.Series {
	return viewmodel.Series{
		Count:     1,
		Publisher: ci.Publisher,
		Series:    ci.Series,
		Volume:    ci.Volume,
	}
}

func (ci ComicInfo) ToIssueDB() viewmodel.Issue {
	id, _ := comicvine.ExtractNumIdFromSiteUrl(ci.Web)
	return viewmodel.Issue{
		ID:        id,
		Title:     ci.Title,
		Number:    ci.Number,
		Publisher: ci.Publisher,
		Series:    ci.Series,
		Volume:    ci.Volume,
		Web:       ci.Web,
	}
}
