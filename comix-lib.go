package main

import (
	"fmt"
	"io/ioutil"
	"log"

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
	chunkSize := cfg.Import.ChunkSize
	importSize := cfg.Import.MaxImport

	cbz := cbz.CBZ{}
	issues := make([]view.Issue, 0)
	newIssues := make([]view.Issue, 0)

	//Start reading CBZ info
	inputFiles, err := ioutil.ReadDir(cfg.Path)
	for err != nil {
		fmt.Println(err)
		panic(err)
	}
	index := 0
	for _, file := range inputFiles {
		if file.IsDir() {
			ifn := fmt.Sprintf("%s/%s", cfg.Path, file.Name())
			log.Printf("Reading: %s\n", ifn)
			infiles, err := ioutil.ReadDir(ifn)
			for err != nil {
				fmt.Println(err)
				panic(err)
			}
			for _, f := range infiles {
				maybeIssue := cbz.BuildIssueFromCBZ(f.Name(), ifn)
				if maybeIssue != nil {
					issues = append(issues, *maybeIssue)
					index++
				}
				if (index % chunkSize) == 0 {
					FilterOutExistingIssues(&newIssues, issues)
					issues = make([]view.Issue, 0)
					if len(newIssues) >= importSize {
						break
					}
				}
			}
		}
		if len(newIssues) >= importSize {
			break
		}
	}

	//enrich issue data
	for i, issue := range newIssues {
		newIssues[i] = EnrichIssueWithVineData(issue)
	}

	// Group issues by Series and generate Series records
	m := view.AsSeriesMap(newIssues)
	series := make([]view.Series, 0)
	for _, v := range m {
		s := comicvine.BuildSeriesFromIssueAndVine(v[0], vine)
		series = append(series, s)
	}

	// store series on DB
	for _, v := range series {
		s := dao.FromSeriesView(v)
		res := db.Create(&s)
		log.Printf("Inserted Series %v with success?%v\n", s.ID, res.RowsAffected)
	}

	//store issues on DB
	for _, v := range newIssues {
		dbIss := dao.FromIssueView(v)

		res := db.Create(&dbIss)
		log.Printf("Inserted Issue %v with success?%v\n", v.ID, res.RowsAffected)
	}

}

func FilterOutExistingIssues(newIssues *[]view.Issue, issues []view.Issue) {
	//check if issue exists on DB
	for _, v := range issues {
		var iss []dao.Issue
		compoundId := vine.ExtractIdFromSiteUrl(v.Web)
		id := vine.ExtractIdFromCompoundId(compoundId)
		db.Find(&iss, id)
		if len(iss) == 0 {
			*newIssues = append(*newIssues, v)
			log.Printf("New issue to add to the DB:%v\n", id)
		}
	}
}

func EnrichIssueWithVineData(issue view.Issue) view.Issue {
	log.Printf("Extracting images for %s", issue.Web)
	issueKey := vine.ExtractIdFromSiteUrl(issue.Web)
	cvIssue := vine.GetIssue(issueKey)
	issue.Images = cvIssue.Image.FromComicVine()
	issue.ID = cvIssue.ID
	issue.VolumeAPI = cvIssue.Volume.ApiDetailUrl
	return issue
}
