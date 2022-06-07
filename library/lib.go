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
	err := db.AutoMigrate(&dao.Series{}, &dao.Issue{}, dao.Publisher{})
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

func (lib ComicsLib) CountIssues() int {
	var count int64
	db.Model(&dao.Issue{}).Count(&count)
	log.Printf("Found %v issues;\n", count)
	return int(count)
}

func (lib ComicsLib) CountIssuesFor(series viewmodel.Series) int {
	var count int64
	db.Model(&dao.Issue{}).Where("series_id = ?", series.ID).Count(&count)
	log.Printf("Found %v issues for %v;\n", count, series.Series)
	return int(count)
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

func (lib ComicsLib) UpdateSeries(s viewmodel.Series) viewmodel.Series {
	sdb := dao.FromSeriesviewmodel(s)
	db.Save(&sdb)
	return sdb.Asviewmodel()
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
	result := db.Limit(pageSize).Offset(page * pageSize).Order("series").Find(&seriesList)
	log.Printf("Found %v series;\n", result.RowsAffected)
	res := make([]viewmodel.Series, len(seriesList))
	for i, v := range seriesList {
		res[i] = v.Asviewmodel()
	}
	return res
}

func (lib ComicsLib) CountSeries() int {
	var count int64
	db.Model(&dao.Series{}).Count(&count)
	log.Printf("Found %v series;\n", count)
	return int(count)
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

func (lib ComicsLib) GetSeriesByPublisher(pub viewmodel.Publisher) []viewmodel.Series {
	var seriesList []dao.Series
	result := db.Order("series").Where("publisher_id = ?", pub.ID).Find(&seriesList)
	log.Printf("Found %v series for publisher %v;\n", result.RowsAffected, pub.Name)
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

// ---- Publisher methods
func (lib ComicsLib) GetAllPublishers() []viewmodel.Publisher {
	var publisherList []dao.Publisher
	result := db.Order("name").Find(&publisherList)
	log.Printf("Found %v publishers;\n", result.RowsAffected)
	res := make([]viewmodel.Publisher, len(publisherList))
	for i, v := range publisherList {
		res[i] = v.AsViewmodel()
	}
	return res
}

func (lib ComicsLib) CountPublishers() int {
	var count int64
	db.Model(&dao.Publisher{}).Count(&count)
	log.Printf("Found %v publishers;\n", count)
	return int(count)
}

func (lib ComicsLib) InsertPublisher(publisher viewmodel.Publisher) int {
	if publisher.ID == 0 {
		return 0
	}
	dbSeries := dao.FromPublisherViewmodel(publisher)
	res := db.Create(&dbSeries)
	return int(res.RowsAffected)
}

func (lib ComicsLib) GetPublisherByID(id int) *viewmodel.Publisher {
	var ser []dao.Publisher
	db.Find(&ser, id)
	if len(ser) == 0 {
		return nil
	}
	res := ser[0].AsViewmodel()
	return &res
}

func (lib ComicsLib) UpdatePublisherCounters() error {
	result := db.Exec("update publishers as ss set series_count = (SELECT count(1) from series where series.publisher_id = ss.id)")
	if result.Error != nil {
		return result.Error
	}
	log.Printf("Updated %v publishers;\n", result.RowsAffected)
	return nil
}
