package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"todoListBE/config"
	"todoListBE/todo/entity"
	"todoListBE/todo/handler"
)
var app App
func TestMain(m *testing.M) {
	app.Initialize()
	app.NewDBConnection(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.MYSQL_USERNAME, config.MYSQL_PASSWORD, config.MYSQL_HOST, config.MYSQL_PORT, config.MYSQL_DATABASE))
	ensureTableExists()
	handler.NewTodoHandler(app.Router,app.Db)
	code := m.Run()
	clearTable()
	os.Exit(code)
}


func ensureTableExists() {
	if _, err := app.Db.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}
func clearTable() {
	app.Db.Exec("DELETE FROM todo")
	app.Db.Exec("ALTER TABLE todo AUTO_INCREMENT = 1")
}
const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS todo
(
    id int(10) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name varchar(255) NOT NULL,
   completed tinyint(1) unsigned NOT NULL DEFAULT '0'
)`

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/todo", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "{\"data\":null,\"success\":1}" {
		t.Errorf("Expected a data is null. Got %s", body)
	}
}
func TestGetNonExistPath(t *testing.T){
	clearTable()
	req, _ := http.NewRequest("GET", "/todo/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}
func TestAddTodo(t *testing.T){
	clearTable()
	var jsonStr = []byte(`{"name":"test todo"}`)
	req, _ := http.NewRequest("POST", "/todo", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)
	var todo entity.Todo
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(),&m)
	jsonString, _ := json.Marshal(m["data"])
	json.Unmarshal(jsonString,&todo)
	if todo.ID != 1 {
		t.Errorf("Expected todo id to be 'test todo'. Got '%v'", todo.ID)
	}
	if todo.Name != "test todo" {
		t.Errorf("Expected todo name to be 'test todo'. Got '%v'", todo.Name)
	}

	if todo.Completed != 0 {
		t.Errorf("Expected todo completed to be 'test todo'. Got '%v'", todo.Completed)
	}
}

func TestUpdateTodo(t *testing.T)  {
	clearTable()
	addTodo()

	req, _ := http.NewRequest("GET", "/todo/1", nil)
	response := executeRequest(req)

	var orginalTodo entity.Todo
	var orginalTodoResponse map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &orginalTodoResponse)
	todoString, _ := json.Marshal(orginalTodoResponse["data"])
	json.Unmarshal(todoString, &orginalTodo)

	var jsonStr = []byte(`{"completed":1}`)
	req, _ = http.NewRequest("PUT", "/todo/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var updatedTodo entity.Todo
	var updatedTodoResponse map[string]interface{}
	json.Unmarshal(response.Body.Bytes(),&updatedTodoResponse)
	updatedTodoString, _ := json.Marshal(updatedTodoResponse["data"])
	json.Unmarshal(updatedTodoString, &updatedTodo)

	if updatedTodo.ID != orginalTodo.ID {
		t.Errorf("Expected todo id (%v). Got '%v'", orginalTodo.ID, updatedTodo.ID)
	}

	if updatedTodo.Completed == orginalTodo.Completed {
		t.Errorf("Expected todo completed (%v). Got '%v'", orginalTodo.Completed, updatedTodo.Completed)
	}
}

func TestDeleteTodo(t *testing.T)  {
	clearTable()
	addTodo()
	req, _ := http.NewRequest("DELETE", "/todo/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/todo/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	return rr
}
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func addTodo()  {
	app.Db.Exec("INSERT INTO todo (name,completed) VALUES(?,?)","test todo",0)
}