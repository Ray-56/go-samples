package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Todo struct {
	ID       uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CreateAt time.Time `gorm:"column:createat;autoCreateTime" json:"crateAt"`
	UpdateAt time.Time `gorm:"column:updateat;autoUpdateTime" json:"updateAt"`
	Title    string    `gorm:"column:title" json:"title"`
	Desc     string    `gorm:"column:desc" json:"desc"`
	Status   bool      `gorm:"column:status" json:"status"`
}

var (
	db *gorm.DB
)

func init() {
	var err error
	dsn := "root@tcp(127.0.0.1:3306)/gin_todolist?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Todo{})
}

func main() {
	router := gin.Default()

	v1 := router.Group("/v1/todo")
	{
		v1.POST("/", createTodo)
		v1.GET("/", fetchAllTodo)
		v1.GET("/:id", fetchSingleTodo)
		v1.PUT("/:id", updateTodo)
		v1.DELETE("/:id", deleteTodo)
	}
	router.Run(":9000")
}

func createTodo(c *gin.Context) {
	var todo Todo
	c.BindJSON(&todo)
	if err := db.Create(&todo).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusCreated, todo)
	}
}

func fetchAllTodo(c *gin.Context) {
	var todoList []Todo
	err := db.Find(&todoList).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error})
		return
	}
	if len(todoList) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not todo found!"})
		return
	}
	c.JSON(http.StatusOK, todoList)
}

func fetchSingleTodo(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusOK, gin.H{"error": "invalid ID"})
		return
	}
	var todo Todo
	db.First(&todo, id)

	c.JSON(http.StatusOK, todo)
}

func updateTodo(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusOK, gin.H{"error": "invalid ID"})
		return
	}
	var todo Todo
	db.First(&todo, id)
	requestBody := make(map[string]interface{})
	c.BindJSON(&requestBody)
	db.Model(&todo).Update("title", requestBody["title"])
	db.Model(&todo).Update("desc", requestBody["desc"])
	c.JSON(http.StatusOK, gin.H{"message": "Todo updated successfully"})
}

func deleteTodo(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusOK, gin.H{"error": "invalid ID"})
		return
	}
	var todo Todo
	db.First(&todo, id)
	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		return
	}

	db.Delete(&todo)
	c.JSON(http.StatusOK, gin.H{"message": "The Todo has been removed"})
}
