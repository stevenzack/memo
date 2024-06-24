package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stevenzack/memo/db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type D map[string]any

var (
	dbc *gorm.DB
)

func main() {
	r := gin.Default()

	//html
	t, e := template.ParseGlob("*.html")
	if e != nil {
		log.Panic(e)
		return
	}
	r.SetHTMLTemplate(t)

	//orm
	dbc, e = gorm.Open(mysql.Open(os.Getenv("MEMO_MYSQL")))
	if e != nil {
		log.Panic(e)
		return
	}
	e = dbc.AutoMigrate(&db.Book{})
	if e != nil {
		log.Panic(e)
		return
	}

	r.GET("/", home)
	r.POST("/books", addbooks)
	r.DELETE("/books/:bid", deleteBook)
	r.GET("/books/:bid", getBook)
	r.Run()
}
func getBook(c *gin.Context) {
	var v db.Book
	e := dbc.First(&v, c.Param("bid")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}
	c.HTML(200, "book.html", v)
}
func deleteBook(c *gin.Context) {
	e := dbc.Delete(&db.Book{}, c.Param("bid")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}
	c.Redirect(http.StatusSeeOther, c.Request.Referer())
}
func home(c *gin.Context) {
	var vs []db.Book
	e := dbc.Find(&vs).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}
	c.HTML(200, "index.html", vs)
}

func addbooks(c *gin.Context) {
	name := c.Request.FormValue("name")
	var v db.Book
	v.Name = name
	e := dbc.Create(&v).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}
	c.Redirect(http.StatusSeeOther, c.Request.Referer())
}
