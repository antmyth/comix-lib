package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/antmyth/comix-lib/cbz"
	"github.com/antmyth/comix-lib/comicvine"
	"github.com/antmyth/comix-lib/config"
	"github.com/antmyth/comix-lib/library"
	"github.com/antmyth/comix-lib/view"
)

var cfg config.Config

var vine comicvine.ComicVine
var lib *library.ComicsLib

func main() {
	cfg = config.ReadConfig()
	vine = comicvine.ComicVine{}
	libz, err := library.New()
	if err != nil {
		panic(err)
	}
	lib = libz

	ss := lib.GetAllSeries()

	for i, v := range ss {
		log.Printf("Series[%v] - %v\n", i, v.Series)
	}
	si := lib.GetAllIssuesFor(ss[2])
	for i, v := range si {
		log.Printf("Issue[%v] - %v:%v\n", i, v.Number, v.Title)
	}
	log.Println(lib.GetSeriesByIDWithIssues(ss[16].ID))
	// BuildLib()
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
		res := lib.InsertSeries(v)
		log.Printf("Inserted Series %v with success?%v\n", v.ID, res)
	}

	//store issues on DB
	for _, v := range newIssues {
		res := lib.InsertIssue(v)
		log.Printf("Inserted Issue %v with success?%v\n", v.ID, res)
	}

}

func FilterOutExistingIssues(newIssues *[]view.Issue, issues []view.Issue) {
	//check if issue exists on DB
	for _, v := range issues {
		id := comicvine.ExtractNumIdFromSiteUrl(v.Web)
		if lib.GetIssueByID(id) == nil {
			*newIssues = append(*newIssues, v)
			log.Printf("New issue to add to the DB:%v\n", id)

		}
	}
}

func EnrichIssueWithVineData(issue view.Issue) view.Issue {
	log.Printf("Extracting images for %s", issue.Web)
	issueKey := comicvine.ExtractIdFromSiteUrl(issue.Web)
	cvIssue := vine.GetIssue(issueKey)
	issue.Images = cvIssue.Image.FromComicVine()
	issue.ID = cvIssue.ID
	issue.VolumeAPI = cvIssue.Volume.ApiDetailUrl
	return issue
}
