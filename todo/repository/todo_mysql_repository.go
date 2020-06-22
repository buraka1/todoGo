package repository

import (
	"database/sql"
	"log"

	"todoListBE/todo"
	"todoListBE/todo/entity"
)

type todoMysqlRepository struct {
	Conn *sql.DB
}

func (tmr *todoMysqlRepository) GetById(todo *entity.Todo) error {
	return tmr.Conn.QueryRow(`SELECT * FROM todo WHERE id = ?`,todo.ID).Scan(&todo.ID,&todo.Name,&todo.Completed)
}

func (tmr *todoMysqlRepository) GetAll() (todos entity.Todos, err error) {
	rows, err := tmr.Conn.Query(`SELECT * FROM todo ORDER BY id ASC`)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var todoEntity entity.Todo
		err = rows.Scan(&todoEntity.ID, &todoEntity.Name, &todoEntity.Completed)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		todos = append(todos, &todoEntity)
	}
	return todos, nil
}

func (tmr *todoMysqlRepository) Delete(t *entity.Todo) error {
	_, err := tmr.Conn.Exec(`DELETE FROM todo WHERE id = ?`,t.ID)
	return err
}

func (tmr *todoMysqlRepository) Update(t *entity.Todo) error {
	res, err := tmr.Conn.Exec(`UPDATE todo SET completed = ? WHERE id = ?`,t.Completed,t.ID)
	if err != nil {
		return  err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return  err
	}
	return nil
}

func (tmr *todoMysqlRepository) Create(t *entity.Todo) error {
	res, err := tmr.Conn.Exec(`INSERT INTO todo (name, completed)  VALUES (?,?)`,t.Name, t.Completed)
	if err != nil {
		return  err
	}
	id,err := res.LastInsertId()
	if err != nil {
		return  err
	}
	t.ID = int(id)
	return nil
}

func NewTodoMysqlRepository(Conn *sql.DB) todo.Repository {
	return &todoMysqlRepository{Conn}
}
