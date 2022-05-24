package dao

import (
	"strconv"
	"time"

	"github.com/antmyth/comix-lib/comicvine"
	"github.com/antmyth/comix-lib/view"
)

var vine comicvine.ComicVine

// gorm.Model definition
type Series struct {
	ID          uint `gorm:"primaryKey"`
	VineId      string
	Series      string
	Volume      string
	Publisher   string `gorm:"index"`
	Count       int
	TotalCount  int
	Web         string
	Location    string
	Description string
	Images      Image `gorm:"embedded;embeddedPrefix:images_"`
	// Issues      []Issue `gorm:"foreignKey:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Issue struct {
	ID             uint `gorm:"primaryKey"`
	Title          string
	Series         string
	Number         string
	Volume         string
	Publisher      string
	Web            string
	VolumeAPI      string
	Images         Image `gorm:"embedded;embeddedPrefix:images_"`
	SeriesLocation string
	Location       string
	SeriesId       uint
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Image struct {
	SmallUrl    string
	ThumbUrl    string
	TinyUrl     string
	OriginalUrl string
}

func FromImageView(v view.Image) Image {
	return Image{
		SmallUrl:    v.SmallUrl,
		ThumbUrl:    v.ThumbUrl,
		TinyUrl:     v.TinyUrl,
		OriginalUrl: v.OriginalUrl,
	}
}

func (issue Issue) FromView(v view.Issue) Issue {
	issue.ID = uint(v.ID)
	issue.Title = v.Title
	issue.Series = v.Series
	issue.Number = v.Number
	issue.Volume = v.Volume
	issue.Publisher = v.Publisher
	issue.Web = v.Web
	issue.VolumeAPI = v.VolumeAPI
	issue.SeriesLocation = v.SeriesLocation
	issue.Location = v.Location
	ci := vine.ExtractIdFromSiteUrl(v.VolumeAPI)
	si := vine.ExtractIdFromCompoundId(ci)
	ssi, _ := strconv.Atoi(si)
	issue.SeriesId = uint(ssi)
	issue.Images = Image{
		SmallUrl:    v.Images.SmallUrl,
		ThumbUrl:    v.Images.ThumbUrl,
		TinyUrl:     v.Images.TinyUrl,
		OriginalUrl: v.Images.OriginalUrl,
	}
	return issue
}
