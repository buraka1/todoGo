package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"todoListBE/config"
	"todoListBE/todo/handler"
)

func main() {
	app := App{}
	app.Initialize()
	app.NewDBConnection(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.MYSQL_USERNAME, config.MYSQL_PASSWORD, config.MYSQL_HOST, config.MYSQL_PORT, config.MYSQL_DATABASE))

	handler.NewTodoHandler(app.Router,app.Db)

	app.Run(fmt.Sprintf(":%s", config.API_PORT))
}


