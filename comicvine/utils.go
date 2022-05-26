package comicvine

import (
	"strconv"

	"github.com/antmyth/comix-lib/view"
)

func BuildSeriesFromIssueAndVine(i view.Issue, vine ComicVine) view.Series {
	s := view.Series{}
	s.Publisher = i.Publisher
	s.Series = i.Series
	s.Volume = i.Volume
	s.Location = i.SeriesLocation

	volKey := vine.ExtractIdFromSiteUrl(i.VolumeAPI)
	volData := vine.GetVolume(volKey)
	s.Images = volData.Image.FromComicVine()
	s.TotalCount = volData.CountOfIssues
	s.ID, _ = strconv.Atoi(vine.ExtractIdFromCompoundId(volKey))
	s.Web = volData.SiteDetailUrl
	s.Description = volData.Description

	return s
}
