package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/antmyth/comix-lib/cbz"
	"github.com/antmyth/comix-lib/comicvine"
	"github.com/antmyth/comix-lib/config"
	"github.com/antmyth/comix-lib/dao"
	"github.com/antmyth/comix-lib/view"
	"gorm.io/gorm"
)

var cfg config.Config
var db *gorm.DB
var vine comicvine.ComicVine

func main() {
	cfg = config.ReadConfig()
	db = dao.GetConnection(cfg)
	vine = comicvine.ComicVine{}

	//migrate the schema
	err := db.AutoMigrate(&dao.Series{}, &dao.Issue{})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	log.Println("Finished DB migrations")

	BuildLib()

}

func BuildLib() {
	cbz := cbz.CBZ{}
	//Start reading CBZ info
	issues := make([]view.Issue, 0)
	inputFiles, err := ioutil.ReadDir(cfg.Path)
	for err != nil {
		fmt.Println(err)
		panic(err)
	}
	for i, file := range inputFiles {
		// ifn := fmt.Sprintf("%s/%s", cfg.Path, file.Name())

		if file.IsDir() {
			ifn := fmt.Sprintf("%s/%s", cfg.Path, file.Name())
			log.Printf("Reading: %s\n", ifn)
			infiles, err := ioutil.ReadDir(ifn)
			for err != nil {
				fmt.Println(err)
				panic(err)
			}
			for ii, f := range infiles {
				maybeIssue := cbz.BuildIssueFromCBZ(f.Name(), ifn)
				if maybeIssue != nil {
					issues = append(issues, *maybeIssue)
				}
				if ii > 1 {
					break
				}
			}
		}

		if i > 2 {
			break
		}
	}

	//check if issue exists on DB
	newIssues := make([]view.Issue, 0)
	for _, v := range issues {
		var iss []dao.Issue
		compoundId := vine.ExtractIdFromSiteUrl(v.Web)
		id := vine.ExtractIdFromCompoundId(compoundId)
		db.Find(&iss, id)
		if len(iss) == 0 {
			newIssues = append(newIssues, v)
			log.Printf("New issue to add to the DB:%v\n", id)
		}
	}

	//enrich issue data
	for i, issue := range newIssues {
		log.Printf("Extracting images for %s", issue.Web)
		issueKey := vine.ExtractIdFromSiteUrl(issue.Web)
		cvIssue := vine.GetIssue(issueKey)
		issue.Images = cvIssue.Image.FromComicVine()
		issue.ID = cvIssue.ID
		issue.VolumeAPI = cvIssue.Volume.ApiDetailUrl
		newIssues[i] = issue
	}

	m := view.AsSeriesMap(newIssues)
	series := make([]view.Series, 0)
	for _, v := range m {
		s := view.Series{}
		i := v[0]
		s.Count = len(v)
		s.Publisher = i.Publisher
		s.Series = i.Series
		s.Volume = i.Volume
		s.Location = i.SeriesLocation
		s.Issues = v

		volKey := vine.ExtractIdFromSiteUrl(i.VolumeAPI)
		volData := vine.GetVolume(volKey)
		s.Images = volData.Image.FromComicVine()
		s.TotalCount = volData.CountOfIssues
		s.ID, _ = strconv.Atoi(vine.ExtractIdFromCompoundId(volKey))
		s.Web = volData.SiteDetailUrl
		s.Description = volData.Description

		series = append(series, s)
	}

	for _, v := range series {
		s := dao.Series{}
		s.ID = uint(v.ID)
		s.VineId = vine.ExtractIdFromSiteUrl(v.Web)
		s.Series = v.Series
		s.Volume = v.Volume
		s.Publisher = v.Publisher
		s.TotalCount = v.TotalCount
		s.Web = v.Web
		s.Location = v.Location
		s.Description = v.Description
		s.Images = dao.FromImageView(v.Images)

		res := db.Create(&s)
		log.Printf("Inserted Series %v with success?%v\n", s.ID, res.RowsAffected)
	}

	//store issues on DB
	for _, v := range newIssues {
		dbIss := dao.Issue{}
		dbIss = dbIss.FromView(v)

		res := db.Create(&dbIss)
		log.Printf("Inserted Issue %v with success?%v\n", v.ID, res.RowsAffected)
	}

}
