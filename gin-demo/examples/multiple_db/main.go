package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Postbrother struct {
	Callid   int64     `gorm:"column:callid" db:"callid" json:"callid" form:"callid"`
	Callee   string    `gorm:"column:callee" db:"callee" json:"callee" form:"callee"`
	Bizlabel string    `gorm:"column:bizlabel" db:"bizlabel" json:"bizlabel" form:"bizlabel"`
	Startime time.Time `gorm:"column:starttime" db:"starttime" json:"starttime" form:"starttime"`
	Endtime  time.Time `gorm:"column:endtime" db:"endtime" json:"endtime" form:"endtime"`
}

func (Postbrother) TabName() string {
	return "postbrother"
}

var dbList [2]*gorm.DB

func init() {
	var err error
	var db1 *gorm.DB
	dsn1 := "root@tcp(127.0.0.1:3306)/cdr1?charset=utf8mb4&parseTime=True&loc=Local"
	db1, err = gorm.Open(mysql.Open(dsn1), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	dbList[0] = db1
	if err != nil {
		panic("failed to connect database")
	}

	var db2 *gorm.DB
	dsn2 := "root@tcp(127.0.0.1:3306)/cdr2?charset=utf8mb4&parseTime=True&loc=Local"
	db2, err = gorm.Open(mysql.Open(dsn2), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	dbList[1] = db2
	if err != nil {
		panic("failed to connect database")
	}

	// db.AutoMigrate(&Todo{})
}

func main() {
	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.GET("/crd", fetchAll)
	}
	router.Run(":9000")
}

type QueryParams struct {
	Callids   []string  `json:"callids" form:"callids"`
	BizLabel  string    `json:"bizlabel" form:"bizlabel"`
	Callees   []string  `json:"callees" form:"callees"`
	Callee    string    `json:"callee" form:"callee"`
	StartTime time.Time `json:"starttime" form:"starttime" binding:"required"`
	EndTime   time.Time `json:"endtime" form:"endtime" binding:"required"`
}

func fetchAll(c *gin.Context) {
	// 处理查询参数中的 + 号
	c.Request.URL.RawQuery = strings.ReplaceAll(c.Request.URL.RawQuery, "+", "%2b")

	var queryParams QueryParams
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		fmt.Println("error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var rsltCdrList []Postbrother
	tables, err := genTables(uint(queryParams.StartTime.Month()), uint(queryParams.EndTime.Month()))
	fmt.Println("tables =>", tables)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	cSize := 2 * len(tables)
	ch := make(chan []Postbrother, cSize)
	var wg sync.WaitGroup

	for _, db := range dbList {
		wg.Add(1)
		go func(db *gorm.DB) {
			defer func() {
				wg.Done()
			}()
			// getTablesCrd(db, queryParams)
			for _, t := range tables {
				wg.Add(1)
				go func(tName string) {
					defer func() {
						wg.Done()
					}()
					queryTable(db, tName, queryParams, ch)
				}(t)
			}
		}(db)
	}
	wg.Wait()
	close(ch)

	for received := range ch {
		rsltCdrList = append(rsltCdrList, received...)
	}

	if len(rsltCdrList) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not todo found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rsltCdrList, "atotal": len(rsltCdrList)})
}

func queryTable(db *gorm.DB, tableName string, queryParams QueryParams, ch chan<- []Postbrother) {
	var cdrlist []Postbrother
	tx := db.Table(tableName).Where("starttime between ? and ?", queryParams.StartTime, queryParams.EndTime)
	if len(queryParams.Callees) > 0 {
		tx = tx.Where("callee In ?", queryParams.Callees)
	}
	if len(queryParams.BizLabel) > 0 {
		tx = tx.Where("bizlabel LIKE ?", queryParams.BizLabel+"%")
	}
	if len(queryParams.Callee) > 0 {
		tx = tx.Where("callee LIKE ?", queryParams.Callee+"%")
	}
	if err := tx.Limit(5000).Find(&cdrlist).Error; err != nil {
		fmt.Println("error", err.Error())
	}

	ch <- cdrlist
}

func genTables(start uint, end uint) ([]string, error) {
	var ret []string

	if start <= 0 || end > 12 {
		return ret, errors.New("between 1 and 12")
	}

	var count uint = 0
	for i := start; i <= end; i++ {
		ret = append(ret, fmt.Sprintf("postbrother_%d", i))
		count++
	}

	return ret, nil
}

// 按照传入日期查询多表数据(同步 -> goroutine + channel)
/* func _getTablesCrd(db *gorm.DB, queryParams QueryParams) ([]Postbrother, error) {
	var allCdrList []Postbrother

	tables, err := genTables(uint(queryParams.StartTime.Month()), uint(queryParams.EndTime.Month()))
	fmt.Println("tables =>", tables)
	if err != nil {
		return allCdrList, err
	}

	var wg sync.WaitGroup
	wg.Add(len(tables))
	for _, t := range tables {
		go func(t string) {
			defer func() {
				wg.Done()
			}()
			var cdrlist []Postbrother
			tx := db.Table(t).Where("starttime between ? and ?", queryParams.StartTime, queryParams.EndTime)
			if len(queryParams.Callees) > 0 {
				tx = tx.Where("callee In ?", queryParams.Callees)
			}
			if len(queryParams.BizLabel) > 0 {
				tx = tx.Where("bizlabel LIKE ?", queryParams.BizLabel+"%")
			}
			if len(queryParams.Callee) > 0 {
				tx = tx.Where("callee LIKE ?", queryParams.Callee+"%")
			}
			if err := tx.Limit(5000).Find(&cdrlist).Error; err != nil {
				fmt.Println("error", err.Error())
			}

			allCdrList = append(allCdrList, cdrlist...)
		}(t)
	}
	wg.Wait()

	return allCdrList, nil
} */
