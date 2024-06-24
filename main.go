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
	e = dbc.AutoMigrate(&db.Book{}, &db.Question{}, &db.Answer{})
	if e != nil {
		log.Panic(e)
		return
	}

	r.GET("/", home)
	r.POST("/books", addbooks)
	r.DELETE("/books/:bid", deleteBook)
	r.GET("/books/:bid", getBook)
	r.GET("/books/:bid/questions", questions)
	r.POST("/books/:bid/questions", addQuestion)
	r.DELETE("/books/:bid/questions/:qid", deleteQuestion)
	r.GET("/books/:bid/questions/:qid", getQuestion)
	r.GET("/books/:bid/questions/:qid/answers", getAnswers)
	r.POST("/books/:bid/questions/:qid/answers", addAnswers)
	r.POST("/books/:bid/questions/:qid/correct", setCorrectAnswer)
	r.GET("/books/:bid/questions/:qid/choose/:aid", chooseAnswer)
	r.Run()
}
func chooseAnswer(c *gin.Context) {
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

	var a db.Answer
	e = dbc.First(&a, c.Param("aid")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}
	if a.QuestionID != q.ID {
		c.String(400, "invalid question ID for answer")
		return
	}

	if !a.IsCorrect {
		c.String(400, "Wrong answer")
		return
	}

	var q2 db.Question
	e = dbc.Limit(1).Where("book_id=? and id>?", b.ID, q.ID).First(&q2).Error
	if e != nil {
		if e == gorm.ErrRecordNotFound {
			c.Redirect(http.StatusSeeOther, "/books/"+c.Param("bid")+"/questions")
			return
		}
		c.AbortWithError(500, e)
		return
	}

	c.Redirect(http.StatusSeeOther, "/books/"+c.Param("bid")+"/questions/"+strconv.FormatUint(uint64(q2.ID), 10))
}
func setCorrectAnswer(c *gin.Context) {
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

	var a db.Answer
	e = dbc.First(&a, c.Request.FormValue("aid")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}
	if a.QuestionID != q.ID {
		c.String(400, "invalid question ID for answer")
		return
	}

	a.IsCorrect = true
	e = dbc.Model(&a).Where("question_id=?", q.ID).Update("is_correct", nil).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	e = dbc.Model(&a).Update("is_correct", true).Error
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

	e = dbc.Create(&db.Answer{
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

	var vs []db.Answer
	e = dbc.Where("question_id = ?", q.ID).Find(&vs).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	c.HTML(200, "answers.html", D{
		"Book":     b,
		"Question": q,
		"Answers":  vs,
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

	var vs []db.Answer
	e = dbc.Where("question_id = ?", q.ID).Find(&vs).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	c.HTML(200, "question.html", D{
		"Book":     b,
		"Question": q,
		"Answers":  vs,
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

	var vs []db.Question
	const size = 20
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 1 {
		page = 1
	}
	e = dbc.Offset((page - 1) * size).Limit(size).Find(&vs).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	c.HTML(200, "questions.html", D{
		"Book":      b,
		"Questions": vs,
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
