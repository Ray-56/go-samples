package main

import (
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {

	file, _ := exec.LookPath(os.Args[0])
	p, _ := filepath.Abs(file)
	index := strings.LastIndex(p, string(os.PathSeparator))
	println(index, p[:index])

	_, absolutePath, _, _ := runtime.Caller(0)
	curDir := filepath.Dir(absolutePath)
	templatesPath := path.Join(curDir, "/templates/*")

	r := gin.Default()
	r.LoadHTMLGlob(templatesPath)
	// r.LoadHTMLGlob("templates/*")
	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})
	r.Run(":8080")
}
