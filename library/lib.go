package library

import (
	"fmt"
	"log"

	"github.com/antmyth/comix-lib/config"
	"github.com/antmyth/comix-lib/dao"
	"github.com/antmyth/comix-lib/view"
	"gorm.io/gorm"
)

var cfg config.Config
var db *gorm.DB

func New() (*ComicsLib, error) {
	log.Println("Starting Lib init")
	cfg = config.ReadConfig()
	db = dao.GetConnection(cfg)

	//migrate the schema
	err := db.AutoMigrate(&dao.Series{}, &dao.Issue{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	log.Println("Finished DB migrations")

	res := ComicsLib{}
	return &res, nil
}

type ComicsLib struct {
}

// ---- Issues methods
func (lib ComicsLib) GetIssueByID(id int) *view.Issue {
	var iss []dao.Issue
	db.Find(&iss, id)
	if len(iss) == 0 {
		return nil
	}
	res := iss[0].AsView()
	return &res
}

func (lib ComicsLib) GetAllIssuesFor(series view.Series) []view.Issue {
	var issueList []dao.Issue
	result := db.Order("number").
		Where("series_id = ?", series.ID).
		Find(&issueList)
	log.Printf("Found %v issues for %v;\n", result.RowsAffected, series.Series)
	res := make([]view.Issue, len(issueList))
	for i, v := range issueList {
		res[i] = v.AsView()
	}
	return res
}

func (lib ComicsLib) InsertIssue(issue view.Issue) int {
	dbIss := dao.FromIssueView(issue)
	res := db.Create(&dbIss)
	return int(res.RowsAffected)
}

// ---- Series methods
func (lib ComicsLib) InsertSeries(series view.Series) int {
	dbSeries := dao.FromSeriesView(series)
	res := db.Create(&dbSeries)
	return int(res.RowsAffected)
}

func (lib ComicsLib) GetSeriesByID(id int) *view.Series {
	var ser []dao.Series
	db.Find(&ser, id)
	if len(ser) == 0 {
		return nil
	}
	res := ser[0].AsView()
	return &res
}

func (lib ComicsLib) GetSeriesByIDWithIssues(id int) *view.Series {
	series := lib.GetSeriesByID(id)
	if series != nil {
		issues := lib.GetAllIssuesFor(*series)
		series.Count = len(issues)
		series.Issues = issues
	}

	return series
}

func (lib ComicsLib) GetAllSeries() []view.Series {
	var seriesList []dao.Series
	result := db.Order("series").Find(&seriesList)
	log.Printf("Found %v series;\n", result.RowsAffected)
	res := make([]view.Series, len(seriesList))
	for i, v := range seriesList {
		res[i] = v.AsView()
	}
	return res
}
