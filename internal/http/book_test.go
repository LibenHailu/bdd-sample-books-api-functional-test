package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"net/http/httptest"

	"github.com/DATA-DOG/go-txdb"
	bookHandler "github.com/LibenHailu/sample-books/internal/http/rest"
	bookModule "github.com/LibenHailu/sample-books/internal/module/book"
	"github.com/LibenHailu/sample-books/internal/storage/persistence"
	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	// we register an sql driver named "txdb"
	datasourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", "root", "admin", "127.0.0.1:3306", "book?parseTime=true")
	txdb.Register("txdb", "mysql", datasourceName)
}

type bookTest struct {
	bookHandler bookHandler.BookHandler
	resp        *httptest.ResponseRecorder
	dbConn      *gorm.DB
	db          *sql.DB
	reqJson     string
	// mock        sqlmock.Sqlmock
}

func (b *bookTest) resetResponse(*godog.Scenario) {
	b.resp = httptest.NewRecorder()

	db, err := sql.Open("txdb", "api")
	if err != nil {
		panic(err)
	}

	dbConn, err := gorm.Open(mysql.New(
		mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})

	if err != nil {
		log.Println(err)
		panic(err)
	}

	bookPersistence := persistence.BookInit(dbConn)
	bookUsecase := bookModule.Initialize(bookPersistence)
	bookHandler := bookHandler.BookInit(bookUsecase)
	b.bookHandler = bookHandler
	b.db = db
}

func (b *bookTest) iSendRequestTo(method, endpoint string) error {
	// req, err := http.NewRequest(method, endpoint, nil)
	var err error

	if err != nil {
		return err
	}
	// handle panic
	defer func() {
		switch t := recover().(type) {
		case string:
			err = fmt.Errorf(t)
		case error:
			err = t
		}
	}()

	switch endpoint {
	case "/v1/books":
		if method == "GET" {
			req, err := http.NewRequest(method, endpoint, nil)

			if err != nil {
				return err
			}
			c, _ := gin.CreateTestContext(b.resp)
			c.Request = req
			b.bookHandler.GetBooks(c)
		} else if method == "POST" {
			req, err := http.NewRequest(method, endpoint, bytes.NewBuffer([]byte(b.reqJson)))
			req.Header.Set("Content-Type", "application/json")

			if err != nil {
				return err
			}
			c, _ := gin.CreateTestContext(b.resp)
			c.Request = req
			b.bookHandler.InsertBook(c)
		}

	case "/v1/books/2":

		if method == "PATCH" {

			req, err := http.NewRequest(method, endpoint, bytes.NewBuffer([]byte(b.reqJson)))
			req.Header.Set("Content-Type", "application/json")
			if err != nil {
				return err
			}
			c, _ := gin.CreateTestContext(b.resp)
			c.Params = []gin.Param{
				{
					Key:   "id",
					Value: "2",
				},
			}
			c.Request = req
			b.bookHandler.UpdateBook(c)
		} else if method == "DELETE" {

			req, err := http.NewRequest(method, endpoint, bytes.NewBuffer([]byte(b.reqJson)))
			req.Header.Set("Content-Type", "application/json")

			if err != nil {
				return err
			}
			c, _ := gin.CreateTestContext(b.resp)
			c.Params = []gin.Param{
				{
					Key:   "id",
					Value: "2",
				},
			}
			c.Request = req
			b.bookHandler.DeleteBook(c)
		}

	default:
		err = fmt.Errorf("unknown endpoint: %s", endpoint)
	}
	return nil
}

func (b *bookTest) theResponseCodeShouldBe(code int) error {
	if code != b.resp.Code {
		if b.resp.Code >= 400 {
			return fmt.Errorf("expected response code to be: %d, but actual is: %d, response message: %s", code, b.resp.Code, string(b.resp.Body.Bytes()))
		}
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", code, b.resp.Code)
	}
	return nil
}

func (b *bookTest) theResponseShouldMatchJson(body *godog.DocString) error {

	var expected, actual interface{}

	// re-encode expected response
	if err := json.Unmarshal([]byte(body.Content), &expected); err != nil {

		return err
	}

	// re-encode actual response too
	if err := json.Unmarshal(b.resp.Body.Bytes(), &actual); err != nil {
		return err
	}

	// the matching may be adapted per different requirements.
	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected JSON does not match actual, %v vs. %v", expected, actual)
	}
	return nil
}

func (b *bookTest) thereAreBooks(books *godog.Table) error {
	var fields []string
	var marks []string
	head := books.Rows[0].Cells
	for _, cell := range head {
		fields = append(fields, cell.Value)
		marks = append(marks, "?")
	}

	stmt, err := b.db.Prepare("INSERT INTO books (" + strings.Join(fields, ", ") + ") VALUES(" + strings.Join(marks, ", ") + ")")
	if err != nil {
		return err
	}
	for i := 1; i < len(books.Rows); i++ {
		var vals []interface{}
		for n, cell := range books.Rows[i].Cells {
			switch head[n].Value {
			case "id":
				vals = append(vals, cell.Value)
			case "isbn":
				vals = append(vals, cell.Value)
			case "title":
				vals = append(vals, cell.Value)
			case "author":
				vals = append(vals, cell.Value)
			case "created_at":
				vals = append(vals, cell.Value)
			case "updated_at":
				vals = append(vals, cell.Value)
			default:
				return fmt.Errorf("unexpected column name: %s", head[n].Value)
			}
		}
		if _, err = stmt.Exec(vals...); err != nil {
			return err
		}
	}
	return nil
}

func (b *bookTest) iHaveRequestJson(req *godog.DocString) error {
	b.reqJson = req.Content
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	book := &bookTest{}
	ctx.BeforeScenario(book.resetResponse)

	ctx.Step(`^I send "([^"]*)" request to "([^"]*)"$`, book.iSendRequestTo)
	ctx.Step(`^the response code should be (\d+)$`, book.theResponseCodeShouldBe)
	ctx.Step(`^the response should match json:$`, book.theResponseShouldMatchJson)
	ctx.Step(`^there are books:$`, book.thereAreBooks)
	ctx.Step(`^I have request json:$`, book.iHaveRequestJson)

	ctx.AfterScenario(func(sc *godog.Scenario, err error) {
		book.db.Close()
		book.reqJson = ""
	})
}
