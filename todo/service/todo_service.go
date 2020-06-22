package service

import (
	"todoListBE/todo"
	"todoListBE/todo/entity"
)

type todoService struct {
	todoRepo todo.Repository
}



func NewTodoService(todoRepo todo.Repository) todo.Service {
	return &todoService{todoRepo: todoRepo}
}

func (service *todoService) GetAll() (ts entity.Todos, err error) {
	return service.todoRepo.GetAll()
}
func (service *todoService) GetById(t *entity.Todo)  error {
	return service.todoRepo.GetById(t)
}
func (service *todoService) Delete(t *entity.Todo) error {
	return service.todoRepo.Delete(t)
}

func (service *todoService) Update(t *entity.Todo) error {
	return service.todoRepo.Update(t)
}

func (service *todoService) Create(t *entity.Todo)  error {
	return service.todoRepo.Create(t)
}
