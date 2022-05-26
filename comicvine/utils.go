package comicvine

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/antmyth/comix-lib/view"
)

func BuildSeriesFromIssueAndVine(i view.Issue, vine ComicVine) view.Series {
	s := view.Series{}
	s.Publisher = i.Publisher
	s.Series = i.Series
	s.Volume = i.Volume
	s.Location = i.SeriesLocation

	volKey := ExtractIdFromSiteUrl(i.VolumeAPI)
	volData := vine.GetVolume(volKey)
	s.Images = volData.Image.FromComicVine()
	s.TotalCount = volData.CountOfIssues
	s.ID = ExtractNumIdFromSiteUrl(i.VolumeAPI)
	s.Web = volData.SiteDetailUrl
	s.Description = volData.Description

	return s
}

func ExtractIdFromSiteUrl(url string) string {
	r, _ := regexp.Compile("\\d+-\\d+")
	issueId := r.FindString(url)
	return issueId
}

func ExtractIdFromCompoundId(compound string) string {
	return strings.Split(compound, "-")[1]
}

func ExtractNumIdFromSiteUrl(url string) int {
	cid := ExtractIdFromSiteUrl(url)
	sid := ExtractIdFromCompoundId(cid)
	id, _ := strconv.Atoi(sid)
	return id

}
