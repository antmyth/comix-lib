package library

import (
	"fmt"
	"log"

	"github.com/antmyth/comix-lib/config"
	"github.com/antmyth/comix-lib/dao"
	"github.com/antmyth/comix-lib/viewmodel"
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
func (lib ComicsLib) GetIssueByID(id int) *viewmodel.Issue {
	var iss []dao.Issue
	db.Find(&iss, id)
	if len(iss) == 0 {
		return nil
	}
	res := iss[0].Asviewmodel()
	return &res
}

func (lib ComicsLib) GetAllIssuesFor(series viewmodel.Series) []viewmodel.Issue {
	var issueList []dao.Issue
	result := db.Order("number").
		Where("series_id = ?", series.ID).
		Find(&issueList)
	log.Printf("Found %v issues for %v;\n", result.RowsAffected, series.Series)
	res := make([]viewmodel.Issue, len(issueList))
	for i, v := range issueList {
		res[i] = v.Asviewmodel()
	}
	return res
}

func (lib ComicsLib) InsertIssue(issue viewmodel.Issue) int {
	dbIss := dao.FromIssueviewmodel(issue)
	res := db.Create(&dbIss)
	return int(res.RowsAffected)
}

// ---- Series methods
func (lib ComicsLib) InsertSeries(series viewmodel.Series) int {
	dbSeries := dao.FromSeriesviewmodel(series)
	res := db.Create(&dbSeries)
	return int(res.RowsAffected)
}

func (lib ComicsLib) GetSeriesByID(id int) *viewmodel.Series {
	var ser []dao.Series
	db.Find(&ser, id)
	if len(ser) == 0 {
		return nil
	}
	res := ser[0].Asviewmodel()
	return &res
}

func (lib ComicsLib) GetSeriesByIDWithIssues(id int) *viewmodel.Series {
	series := lib.GetSeriesByID(id)
	if series != nil {
		issues := lib.GetAllIssuesFor(*series)
		series.Count = len(issues)
		series.Issues = issues
	}

	return series
}

func (lib ComicsLib) GetAllSeriesPaginated(page, pageSize int) []viewmodel.Series {
	var seriesList []dao.Series
	log.Printf("Finding series: page[%v] & pageSize[%v]\n", page, pageSize)
	result := db.Order("series").Find(&seriesList).Limit(pageSize).Offset(page * pageSize)
	log.Printf("Found %v series;\n", result.RowsAffected)
	res := make([]viewmodel.Series, len(seriesList))
	for i, v := range seriesList {
		res[i] = v.Asviewmodel()
	}
	return res
}

func (lib ComicsLib) GetAllSeries() []viewmodel.Series {
	var seriesList []dao.Series
	result := db.Order("series").Find(&seriesList)
	log.Printf("Found %v series;\n", result.RowsAffected)
	res := make([]viewmodel.Series, len(seriesList))
	for i, v := range seriesList {
		res[i] = v.Asviewmodel()
	}
	return res
}

func (lib ComicsLib) UpdateSeriesCounters() error {
	result := db.Exec("update series as ss set count = (SELECT count(1) from issues where issues.series_id = ss.id)")
	if result.Error != nil {
		return result.Error
	}
	log.Printf("Updated %v series;\n", result.RowsAffected)
	return nil
}
