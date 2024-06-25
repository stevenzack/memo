package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stevenzack/memo/db"
	"github.com/stevenzack/memo/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type D map[string]any

var (
	dbc *gorm.DB
)

func init() {
	if os.Getenv("GIN_MODE") == "release" {
		gin.DisableConsoleColor()
		fo, e := os.OpenFile("log.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if e != nil {
			log.Panic(e)
			return
		}
		gin.DefaultWriter = io.MultiWriter(fo)
		log.SetFlags(log.Lshortfile)
		log.SetOutput(fo)
	}
}

func main() {
	r := gin.Default()

	//html
	var t = template.New("views").Funcs(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	})
	t, e := t.ParseGlob("*.html")
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
	e = dbc.AutoMigrate(&db.Book{}, &db.Question{}, &db.Option{})
	if e != nil {
		log.Panic(e)
		return
	}

	r.GET("/", home)
	r.POST("/books", addbooks)
	r.DELETE("/books/:bid", deleteBook)
	r.GET("/books/:bid", getBook)
	r.POST("/books/:bid", updateBook)
	r.GET("/books/:bid/questions", questions)
	r.POST("/books/:bid/questions", addQuestion)
	r.DELETE("/books/:bid/questions/:qid", deleteQuestion)
	r.GET("/books/:bid/questions/:qid", getQuestion)
	r.POST("/books/:bid/questions/:qid", updateQuestion)
	r.GET("/books/:bid/questions/:qid/options", getAnswers)
	r.POST("/books/:bid/questions/:qid/options", addAnswers)
	r.POST("/books/:bid/questions/:qid/options/:oid", updateAnswer)
	r.DELETE("/books/:bid/questions/:qid/options/:oid", deleteAnswer)
	r.GET("/books/:bid/questions/:qid/options/:oid", getAnswer)
	r.Run()
}
func getAnswer(c *gin.Context) {
	var a db.Option
	e := dbc.First(&a, c.Param("oid")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}
	c.HTML(200, "option.html", a)
}
func deleteAnswer(c *gin.Context) {
	e := dbc.Delete(&db.Option{}, c.Param("oid")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}
	c.Redirect(http.StatusSeeOther, c.Request.Referer())
}
func updateAnswer(c *gin.Context) {
	e := dbc.Model(&db.Option{}).Where("id=?", c.Param("oid")).Update("text", c.Request.FormValue("text")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}
	c.Redirect(http.StatusSeeOther, c.Request.Referer()+"/..")
}
func updateQuestion(c *gin.Context) {
	var b db.Book
	e := dbc.First(&b, c.Param("bid")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	var q db.Question
	e = dbc.First(&q, c.Param("qid")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	if q.BookID != b.ID {
		c.String(400, "invalid book ID for question")
		return
	}

	e = dbc.Model(&q).Update("text", c.Request.FormValue("text")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}
	c.Redirect(http.StatusSeeOther, c.Request.Referer())
}
func updateBook(c *gin.Context) {
	var b db.Book
	e := dbc.First(&b, c.Param("bid")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	e = dbc.Model(&b).Update("name", c.Request.FormValue("name")).Update("desc", c.Request.FormValue("desc")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}
	c.Redirect(http.StatusSeeOther, c.Request.Referer())
}
func addAnswers(c *gin.Context) {
	var b db.Book
	e := dbc.First(&b, c.Param("bid")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	var q db.Question
	e = dbc.First(&q, c.Param("qid")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	if q.BookID != b.ID {
		c.String(400, "invalid book ID for question")
		return
	}

	e = dbc.Create(&db.Option{
		QuestionID: q.ID,
		Text:       c.Request.FormValue("text"),
	}).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}
	c.Redirect(http.StatusSeeOther, c.Request.Referer())
}
func getAnswers(c *gin.Context) {
	var b db.Book
	e := dbc.First(&b, c.Param("bid")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	var q db.Question
	e = dbc.First(&q, c.Param("qid")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	if q.BookID != b.ID {
		c.String(400, "invalid book ID for question")
		return
	}

	var vs []db.Option
	e = dbc.Where("question_id = ?", q.ID).Find(&vs).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	c.HTML(200, "options.html", D{
		"Book":     b,
		"Question": q,
		"Options":  vs,
	})
}
func getQuestion(c *gin.Context) {
	var b db.Book
	e := dbc.First(&b, c.Param("bid")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	var q db.Question
	e = dbc.First(&q, c.Param("qid")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	if q.BookID != b.ID {
		c.String(400, "invalid book ID for question")
		return
	}

	var vs []db.Option
	e = dbc.Where("question_id = ?", q.ID).Find(&vs).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	c.HTML(200, "question.html", D{
		"Book":     b,
		"Question": q,
		"Options":  vs,
	})
}
func deleteQuestion(c *gin.Context) {
	var b db.Book
	e := dbc.First(&b, c.Param("bid")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	e = dbc.Where("book_id = ?", b.ID).Delete(&db.Question{}, c.Param("qid")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}
}
func addQuestion(c *gin.Context) {
	var b db.Book
	e := dbc.First(&b, c.Param("bid")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	v := db.Question{
		BookID: b.ID,
		Text:   c.Request.FormValue("text"),
	}
	e = dbc.Create(&v).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}
	c.Redirect(http.StatusSeeOther, c.Request.Referer())
}
func questions(c *gin.Context) {
	var b db.Book
	e := dbc.First(&b, c.Param("bid")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	var total int64
	e = dbc.Model(&db.Question{}).Where("book_id=?", b.ID).Count(&total).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	var vs []db.Question
	const size = 10
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 1 {
		page = 1
	}
	e = dbc.Where("book_id=?", b.ID).Offset((page - 1) * size).Limit(size).Find(&vs).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	c.HTML(200, "questions.html", D{
		"Book":      b,
		"Questions": vs,
		"Total":     total,
		"Page":      page,
		"TotalPage": util.PageNum(total, size),
	})
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
