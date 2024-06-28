package main

import (
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stevenzack/memo/db"
	"github.com/stevenzack/memo/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type D map[string]any

const (
	resourceDir = "static"
)

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
		"subAbs": util.SubAbs,
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
	r.Static("/"+resourceDir, "./"+resourceDir)
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
	os.Remove(resourceDir + "/books/" + c.Param("bid") + "/questions/" + c.Param("qid") + "/options/" + c.Param("oid"))

	e := dbc.Delete(&db.Option{}, c.Param("oid")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}
	c.Redirect(http.StatusSeeOther, c.Request.Referer())
}
func updateAnswer(c *gin.Context) {
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

	var o db.Option
	e = dbc.First(&o, c.Param("oid")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	dstDir := "books/" + c.Param("bid") + "/questions/" + c.Param("qid") + "/options/" + c.Param("oid")
	audio, e := readStaticFile(c, "audio", dstDir)
	if e != nil {
		c.AbortWithError(500, e)
		return
	}
	if audio != nil {
		os.Remove(o.Audio.String)
	}
	video, e := readStaticFile(c, "video", dstDir)
	if e != nil {
		c.AbortWithError(500, e)
		return
	}
	if video != nil {
		os.Remove(o.Video.String)
	}

	e = dbc.Model(&db.Option{}).Where("id=?", c.Param("oid")).Update("text", c.Request.FormValue("text")).Update("video", video).Update("audio", audio).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}
	c.Redirect(http.StatusSeeOther, c.Request.Referer()+"/..")
}

func readStaticFile(c *gin.Context, formName string, dstDir string) (*string, error) {
	fh, e := c.FormFile(formName)
	if e != nil {
		if e == http.ErrMissingFile {
			return nil, nil
		}
		log.Println(e)
		return nil, e
	}
	fi, e := fh.Open()
	if e != nil {
		log.Println(e)
		return nil, e
	}
	defer fi.Close()
	out := resourceDir + "/" + dstDir + "/" + strconv.Itoa(rand.Intn(100000)) + filepath.Ext(fh.Filename)
	os.MkdirAll(filepath.Dir(out), 0755)
	fo, e := os.OpenFile(out, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if e != nil {
		log.Println(e)
		return nil, e
	}
	defer fo.Close()

	_, e = io.Copy(fo, fi)
	if e != nil {
		log.Println(e)
		return nil, e
	}
	return &out, nil
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

	m := dbc.Model(&q)
	dstDir := "books/" + c.Param("bid") + "/questions/" + c.Param("qid")
	audio, e := readStaticFile(c, "audio", dstDir)
	if e != nil {
		c.AbortWithError(500, e)
		return
	}
	if audio != nil {
		os.Remove(q.Audio.String)
		m = m.Update("audio", audio)
	}
	video, e := readStaticFile(c, "video", dstDir)
	if e != nil {
		c.AbortWithError(500, e)
		return
	}
	if video != nil {
		os.Remove(q.Video.String)
		m = m.Update("video", video)
	}

	text := c.Request.FormValue("text")
	if text != "" {
		m = m.Update("text", text)
	}

	playNext := false
	incDone := false
	// slayed
	slayed := c.Request.FormValue("slayed")
	if slayed != "" {
		var i any = nil
		if slayed == "true" {
			i = 1
		}
		m = m.Update("slayed", i)
		playNext = true
		incDone = true
	}

	//wrong
	wrong := c.Request.FormValue("wrong")
	if wrong != "" {
		m = m.Update("wrong_count", q.WrongCount.Int16+1)
		playNext = true
	}

	// done
	done := c.Request.FormValue("done")
	if done != "" {
		var i any = nil
		if done == "true" {
			i = 1
		}
		m = m.Update("done", i)
		playNext = true
		incDone = true
	}

	e = m.Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	// inc done
	if incDone {
		e = dbc.Model(&b).Update("last_done_at", time.Now()).Error
		if e != nil {
			c.AbortWithError(500, e)
			return
		}

		if !q.FirstReview.Valid {
			e = dbc.Model(&b).Update("today_done", b.TodayDone+1).Error
			if e != nil {
				c.AbortWithError(500, e)
				return
			}
			e = dbc.Model(&q).Where("first_review is null").Update("first_review", time.Now()).Error
			if e != nil {
				c.AbortWithError(500, e)
				return
			}
		}

	}

	if playNext {
		c.Redirect(http.StatusSeeOther, "/books/"+c.Param("bid")+"?action=play")
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

	var total int64
	e = dbc.Model(&db.Question{}).Where("book_id=?", b.ID).Count(&total).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	dailyPlan, _ := strconv.ParseUint(c.Request.FormValue("daily_plan"), 10, 16)
	if dailyPlan > uint64(total) {
		dailyPlan = uint64(total)
	}

	e = dbc.Model(&b).Update("name", c.Request.FormValue("name")).Update("desc", c.Request.FormValue("desc")).Update("daily_plan", uint16(dailyPlan)).Error
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

	var q db.Question
	e = dbc.First(&q, c.Param("qid")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	os.RemoveAll(resourceDir + "/books/" + c.Param("bid") + "/questions/" + c.Param("qid"))

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

	q := db.Question{
		BookID: b.ID,
		Text:   c.Request.FormValue("text"),
	}
	e = dbc.Create(&q).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	//option
	opt := c.Request.FormValue("option")
	if opt != "" {
		e = dbc.Create(&db.Option{
			QuestionID: q.ID,
			Text:       opt,
		}).Error
		if e != nil {
			c.AbortWithError(500, e)
			return
		}
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
	const size = 30
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
	var b db.Book
	e := dbc.First(&b, c.Param("bid")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	switch c.Query("action") {
	case "play":
		// play
		if b.DailyPlan <= 0 {
			c.String(200, "Daily plan of this book is not set.")
			return
		}

		// review yesterday's questions
		if b.LastDoneAt.Time.Day() != time.Now().Day() {
			var args = []any{b.ID}
			const reviewWhere = "first_review between ? and ?"
			reviewTime := []string{}
			// yesterday
			reviewTime = append(reviewTime, reviewWhere)
			args = append(args, util.YesterdayAgo(time.Now())...)

			// 3 days ago
			reviewTime = append(reviewTime, reviewWhere)
			args = append(args, util.ThreeDaysAgo(time.Now())...)

			// 7 days ago
			reviewTime = append(reviewTime, reviewWhere)
			args = append(args, util.SevenDaysAgo(time.Now())...)

			// 1 month ago
			reviewTime = append(reviewTime, reviewWhere)
			args = append(args, util.OneMonthAgo(time.Now())...)

			e = dbc.Model(&db.Question{}).Where("book_id=? and done=1 and slayed is null and ("+strings.Join(reviewTime, " or ")+")", args...).Update("done", nil).Error
			if e != nil {
				c.AbortWithError(500, e)
				return
			}
			log.Println("review triggered")
		}

		// today wrong count
		var todayWrongCount int64
		e = dbc.Model(&db.Question{}).Where("book_id=? and done is null and slayed is null and wrong_count>0", b.ID).Count(&todayWrongCount).Error
		if e != nil {
			c.AbortWithError(500, e)
			return
		}

		remains := util.SubAbs(b.DailyPlan, b.TodayDone)
		if remains <= uint16(todayWrongCount) {
			// today's wrong questions
			var q db.Question
			e = dbc.Model(&db.Question{}).Where("slayed is null and done is null and wrong_count>0").Order("id asc").Limit(1).First(&q).Error
			if e != nil && e != gorm.ErrRecordNotFound {
				c.AbortWithError(500, e)
				return
			}
			if e == nil {
				c.Redirect(http.StatusSeeOther, "/books/"+c.Param("bid")+"/questions/"+strconv.FormatUint(uint64(q.ID), 10))
				return
			}
		}

		// find the next not done
		var q db.Question
		e = dbc.Model(&db.Question{}).Where("slayed is null and done is null and wrong_count is null").Order("id asc").Limit(1).First(&q).Error
		if e != nil {
			if e == gorm.ErrRecordNotFound {
				// next round
				e = dbc.Model(&b).Update("round", b.Round+1).Update("today_done", 0).Error
				if e != nil {
					c.AbortWithError(500, e)
					return
				}
				e = dbc.Model(&db.Question{}).Where("book_id=?", b.ID).Update("done", nil).Update("first_review", nil).Update("wrong_count", nil).Error
				if e != nil {
					c.AbortWithError(500, e)
					return
				}

				c.HTML(200, "done.html", b)
				return
			}
			c.AbortWithError(500, e)
			return
		}

		if b.TodayDone >= b.DailyPlan {
			c.HTML(200, "done-today.html", b)
			return
		}
		c.Redirect(http.StatusSeeOther, "/books/"+c.Param("bid")+"/questions/"+strconv.FormatUint(uint64(q.ID), 10))
		return
	case "reset":
		//reset progress
		e = dbc.Model(b).Update("round", 0).Update("daily_plan", 0).Update("today_done", 0).Error
		if e != nil {
			log.Println(e)
			c.String(500, e.Error())
			return
		}

		e = dbc.Model(&db.Question{}).Where("book_id=?", b.ID).Update("slayed", nil).Update("done", nil).Update("wrong_count", nil).Update("first_review", nil).Error
		if e != nil {
			c.AbortWithError(500, e)
			return
		}
		c.Redirect(http.StatusSeeOther, "/books/"+c.Param("bid"))
		return
	}

	if time.Now().Day() != b.UpdatedAt.Day() {
		b.TodayDone = 0
		e = dbc.Model(&b).Update("today_done", 0).Error
		if e != nil {
			c.AbortWithError(500, e)
			return
		}
	}
	var total int64
	e = dbc.Model(&db.Question{}).Where("book_id=?", b.ID).Count(&total).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}
	var slayed int64
	e = dbc.Model(&db.Question{}).Where("book_id=? and slayed =1", b.ID).Count(&slayed).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	var remains int64
	e = dbc.Model(&db.Question{}).Where("book_id=? and slayed is null and done is null", b.ID).Count(&remains).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}

	c.HTML(200, "book.html", D{
		"Book":    b,
		"Remains": remains,
		"Total":   total,
		"Slayed":  slayed,
	})
}

func deleteBook(c *gin.Context) {
	e := dbc.Delete(&db.Book{}, c.Param("bid")).Error
	if e != nil {
		c.AbortWithError(500, e)
		return
	}
	os.RemoveAll(resourceDir + "/books/" + c.Param("bid"))
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
