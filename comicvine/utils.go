package comicvine

import (
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/antmyth/comix-lib/viewmodel"
)

func BuildSeriesFromIssueAndVine(i viewmodel.Issue, vine ComicVine) (viewmodel.Series, error) {
	s := viewmodel.Series{}
	s.Publisher = i.Publisher
	s.Series = i.Series
	s.Volume = i.Volume
	s.Location = i.SeriesLocation

	volKey := ExtractIdFromSiteUrl(i.VolumeAPI)
	volData := vine.GetVolume(volKey)
	s.Images = volData.Image.FromComicVine()
	s.TotalCount = volData.CountOfIssues
	sid, err := ExtractNumIdFromSiteUrl(i.VolumeAPI)
	if err != nil {
		return s, err
	}
	s.ID = sid
	s.Web = volData.SiteDetailUrl
	s.Description = volData.Description

	return s, nil
}

func ExtractIdFromSiteUrl(url string) string {
	r, _ := regexp.Compile("\\d{4}-\\d{3,}")
	issueId := r.FindString(url)
	return issueId
}

func extractIdFromCompoundId(compound string) (string, error) {
	split := strings.Split(compound, "-")
	if len(split) < 2 {
		return split[0], errors.New("No compound id found.")
	}
	return split[1], nil
}

func ExtractNumIdFromSiteUrl(url string) (int, error) {
	cid := ExtractIdFromSiteUrl(url)
	sid, err := extractIdFromCompoundId(cid)
	if err != nil {
		log.Printf("Error extracting id from:%v \n%+v\n", url, err)
		return 0, err
	}
	id, _ := strconv.Atoi(sid)
	return id, nil

}
