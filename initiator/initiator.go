package initiator

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/LibenHailu/sample-books/internal/constant/model"
	"github.com/LibenHailu/sample-books/internal/glue/routing"
	"github.com/LibenHailu/sample-books/internal/http/rest"
	"github.com/LibenHailu/sample-books/internal/module/book"
	"github.com/LibenHailu/sample-books/internal/storage/persistence"
	ginRouter "github.com/LibenHailu/sample-books/platform/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Initiate() {

	// DSN := "postgres://postgres:admin@localhost:5432/book?sslmode=disable"

	// conn, err := gorm.Open(postgres.Open(DSN), &gorm.Config{
	// 	SkipDefaultTransaction: true,
	// })

	// mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", "root", "admin", "127.0.0.1:3306", "book")

	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// conn.AutoMigrare automigrates gorm models
	conn.AutoMigrate(&model.Book{})

	if err != nil {
		log.Printf("Error when Opening database connection: %v", err)
		os.Exit(1)
	}

	// Get generic database object sql.DB to use its functions
	sqlDB, err := conn.DB()

	if err != nil {
		log.Printf("Error when connecting to gorm pool: %v", err)
		os.Exit(1)
	}

	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetConnMaxLifetime(15 * time.Minute)

	defer sqlDB.Close()

	bookPersistence := persistence.BookInit(conn)
	bookUsecase := book.Initialize(bookPersistence)
	bookHandler := rest.BookInit(bookUsecase)
	bookRouting := routing.BookRouting(bookHandler)

	var routersList []ginRouter.Router
	routersList = append(routersList, bookRouting...)

	host := "localhost"
	port := "8080"
	srv := ginRouter.NewRouting(host, port, routersList)

	srv.Serve()
}
