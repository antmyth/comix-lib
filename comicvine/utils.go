package comicvine

import (
	"errors"
	"log"
	"path"
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

	if len(i.VolumeAPI) == 0 {
		// extract volume data from vine
		issueId, err := ExtractNumIdFromSiteUrl(i.Web)
		if err != nil {
			log.Printf("Error extracting issue id %v : %v\n", i.ID, i.Web)
			return s, err
		}
		vi, err := vine.GetIssueBy(issueId)
		if err != nil {
			log.Printf("Error getting issue info %v : %v\n", issueId, vi)
			return s, err
		}
		i.VolumeAPI = vi.Volume.ApiDetailUrl
	}
	volKey := ExtractIdFromSiteUrl(i.VolumeAPI)
	volData, err := vine.GetVolume(volKey)
	if err != nil {
		log.Printf("Error extracting Volume data for %v : %v\n", i.VolumeAPI, i.Series)
		return s, err
	}
	s.Images = volData.Image.FromComicVine()
	s.TotalCount = volData.CountOfIssues
	sid, err := ExtractNumIdFromSiteUrl(i.VolumeAPI)
	if err != nil {
		return s, err
	}
	s.ID = sid
	s.Web = volData.SiteDetailUrl
	s.Description = volData.Description
	s.PublisherId = volData.Publisher.ID

	return s, nil
}

func ExtractIdFromSiteUrl(url string) string {
	issueId := path.Base(url)
	// r, _ := regexp.Compile("\\d{4}-\\d{3,}")
	// issueId := r.FindString(url)
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
