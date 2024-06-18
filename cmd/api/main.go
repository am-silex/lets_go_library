package main

import (
	"context"
	"database/sql"
	"fmt"
	data "github.com/am-silex/go_library/internal/data"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

var lock = &sync.Mutex{}

type config struct {
	dbHost string
	dbPort int
	dbUser string
	dbPass string
	dbName string

	webPort int
}

type application struct {
	config config
	models data.Models
	logger *log.Logger
	wg     sync.WaitGroup
	db     *sql.DB
}

type Application interface {
	Serve() error
}

var app *application

func getApp() *application {
	if app == nil {
		lock.Lock()
		defer lock.Unlock()
		if app == nil {
			app = &application{}
		}
	}
	return app
}

func main() {

	// init app
	app = getApp()

	// config & init parameters
	configLogger(app)
	configApp(app)

	db, err := openDB(app.config)
	if err != nil {
		app.logger.Println(err, nil)
	}
	app.db = db

	defer db.Close()

	app.logger.Println("database connection pool established", nil)

	app.models = data.NewModels(db)

	// Start Http server
	err = app.Serve()
	if err != nil {
		app.logger.Fatalln(err)
	}

}

func configApp(app *application) {
	app.config.dbHost = os.Getenv("DB_HOST")
	app.config.dbUser = os.Getenv("DB_USER")
	app.config.dbPass = os.Getenv("DB_PASS")
	app.config.dbName = os.Getenv("DB_NAME")
	app.config.dbPort, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	app.config.webPort, _ = strconv.Atoi(os.Getenv("WEB_PORT"))
}

func configLogger(app *application) {
	l := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	app.logger = l
}

func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the
	// config struct.
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.dbHost, cfg.dbPort, cfg.dbUser, cfg.dbPass, cfg.dbName)
	// Debugging
	app.logger.Println(psqlInfo)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	// Return the sql.DB connection pool.
	return db, nil

}
