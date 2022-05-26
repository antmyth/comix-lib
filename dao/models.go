package dao

import (
	"time"

	"github.com/antmyth/comix-lib/comicvine"
	"github.com/antmyth/comix-lib/viewmodel"
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

func (i Image) Asviewmodel() viewmodel.Image {
	return viewmodel.Image{
		SmallUrl:    i.SmallUrl,
		ThumbUrl:    i.ThumbUrl,
		TinyUrl:     i.TinyUrl,
		OriginalUrl: i.OriginalUrl,
	}
}

func FromImageviewmodel(v viewmodel.Image) Image {
	return Image{
		SmallUrl:    v.SmallUrl,
		ThumbUrl:    v.ThumbUrl,
		TinyUrl:     v.TinyUrl,
		OriginalUrl: v.OriginalUrl,
	}
}

func (v Issue) Asviewmodel() viewmodel.Issue {
	return viewmodel.Issue{
		ID:             int(v.ID),
		Title:          v.Title,
		Series:         v.Series,
		SeriesId:       int(v.SeriesId),
		Number:         v.Number,
		Volume:         v.Volume,
		Publisher:      v.Publisher,
		Web:            v.Web,
		VolumeAPI:      v.VolumeAPI,
		SeriesLocation: v.SeriesLocation,
		Location:       v.Location,
		Images:         v.Images.Asviewmodel(),
	}
}

func FromIssueviewmodel(v viewmodel.Issue) Issue {
	issue := Issue{}
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
	ssi := comicvine.ExtractNumIdFromSiteUrl(v.VolumeAPI)
	issue.SeriesId = uint(ssi)
	issue.Images = FromImageviewmodel(v.Images)
	return issue
}

func (s Series) Asviewmodel() viewmodel.Series {
	return viewmodel.Series{
		ID:          int(s.ID),
		VineId:      s.VineId,
		Series:      s.Series,
		Volume:      s.Volume,
		Publisher:   s.Publisher,
		Count:       s.Count,
		TotalCount:  s.TotalCount,
		Web:         s.Web,
		Images:      s.Images.Asviewmodel(),
		Location:    s.Location,
		Description: s.Description,
	}
}

func FromSeriesviewmodel(v viewmodel.Series) Series {
	s := Series{}
	s.ID = uint(v.ID)
	s.VineId = comicvine.ExtractIdFromSiteUrl(v.Web)
	s.Series = v.Series
	s.Volume = v.Volume
	s.Publisher = v.Publisher
	s.TotalCount = v.TotalCount
	s.Web = v.Web
	s.Location = v.Location
	s.Description = v.Description
	s.Images = FromImageviewmodel(v.Images)
	return s
}
