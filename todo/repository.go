package todo

import "todoListBE/todo/entity"

type Repository interface {
	GetAll() (todos entity.Todos, err error)
	Delete(todo *entity.Todo) error
	Update(todo *entity.Todo) error
	Create(todo *entity.Todo) error
	GetById(todo *entity.Todo) error
}