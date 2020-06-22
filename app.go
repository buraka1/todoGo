package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type App struct {
	Db *sql.DB
	Router *gin.Engine
}

func (a *App) Initialize()  {
	a.Router = gin.Default()
	a.Router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin","Content-Type"},
		ExposeHeaders:    []string{"Content-Length","Content-Type"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

}

func (a *App) NewDBConnection(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	err = db.Ping()
	if err != nil {
		return nil
	}
	a.Db = db
	return db
}

func (a *App) Run(addr string) {
	a.Router.Run(addr)
}