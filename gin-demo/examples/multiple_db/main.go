package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

var db1 *gorm.DB
var db2 *gorm.DB

func init() {
	var err error
	dsn1 := "root@tcp(127.0.0.1:3306)/cdr1?charset=utf8mb4&parseTime=True&loc=Local"
	db1, err = gorm.Open(mysql.Open(dsn1), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	dsn2 := "root@tcp(127.0.0.1:3306)/cdr2?charset=utf8mb4&parseTime=True&loc=Local"
	db2, err = gorm.Open(mysql.Open(dsn2), &gorm.Config{})
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
	BizLabel  string    `json:"bizlabel" form:"bizlabel"`
	Callee    string    `json:"callee" form:"callee"`
	StartTime time.Time `json:"starttime" form:"starttime" binding:"required"`
	EndTime   time.Time `json:"endtime" form:"endtime" binding:"required"`
}

func fetchAll(c *gin.Context) {
	// 处理查询参数中的 + 号
	c.Request.URL.RawQuery = strings.ReplaceAll(c.Request.URL.RawQuery, "+", "%2b")

	var queryParams QueryParams
	if err := c.ShouldBind(&queryParams); err != nil {
		fmt.Println("error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var allCdrList []Postbrother
	crd1List, err := getTablesCrd(db1, queryParams)
	fmt.Println("crd1List =>", crd1List)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	allCdrList = append(allCdrList, crd1List...)

	crd2List, err := getTablesCrd(db2, queryParams)
	fmt.Println("crd2List =>", crd2List)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	allCdrList = append(allCdrList, crd2List...)

	if len(allCdrList) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not todo found!"})
		return
	}
	c.JSON(http.StatusOK, allCdrList)
}

// 按照传入日期查询多表数据
func getTablesCrd(db *gorm.DB, queryParams QueryParams) ([]Postbrother, error) {
	var allCdrList []Postbrother

	tables, err := genTables(uint(queryParams.StartTime.Month()), uint(queryParams.EndTime.Month()))
	fmt.Println("tables =>", tables)
	if err != nil {
		return allCdrList, err
	}

	for _, t := range tables {
		var cdrlist []Postbrother
		// fmt.Sprintf("%%%s%%", queryParams.Callee)
		if err := db.Table(t).Where("callee LIKE ?", "%"+queryParams.Callee+"%").Where("bizlabel LIKE ?", fmt.Sprintf("%%%s%%", queryParams.BizLabel)).Where("starttime >= ?", queryParams.StartTime).Where("starttime <= ?", queryParams.EndTime).Find(&cdrlist).Error; err != nil {
			return allCdrList, err
		}
		allCdrList = append(allCdrList, cdrlist...)
	}

	return allCdrList, nil
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
