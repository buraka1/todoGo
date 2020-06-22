package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"todoListBE/todo"
	"todoListBE/todo/entity"
	"todoListBE/todo/repository"
	"todoListBE/todo/service"
)

type TodoHandler struct {
	TodoService todo.Service
}

func (h TodoHandler) fetchAll(context *gin.Context) {
	data, _ := h.TodoService.GetAll()
	context.JSON(http.StatusOK, gin.H{"success": 1, "data": data})
	return
}

func (h TodoHandler) deleteTodo(context *gin.Context) {
	todoEntity := entity.NewTodo()
	todoEntity.ID, _ = strconv.Atoi(context.Param("id"))
	if err := h.TodoService.Delete(todoEntity); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"success": 0, "message": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"success": 1})
	return
}

func (h TodoHandler) createTodo(context *gin.Context) {
	todoEntity := entity.NewTodo()
	_ = context.BindJSON(&todoEntity)
	if err := h.TodoService.Create(todoEntity); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"success": 0, "message": err.Error()})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"success": 1, "data": todoEntity})
	return
}

func (h TodoHandler) updateTodo(context *gin.Context) {
	todoEntity := entity.NewTodo()
	_ = context.BindJSON(&todoEntity)
	todoEntity.ID, _ = strconv.Atoi(context.Param("id"))
	if err := h.TodoService.Update(todoEntity); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"success": 0, "message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"success": 1, "data": todoEntity})
	return
}

func (h TodoHandler) getTodo(context *gin.Context) {
	todoEntity := entity.NewTodo()
	todoEntity.ID, _ = strconv.Atoi(context.Param("id"))

	if err := h.TodoService.GetById(todoEntity); err != nil {
		switch err {
		case sql.ErrNoRows:
			context.JSON(http.StatusNotFound, gin.H{"success": 0})
			return
		default:
			context.JSON(http.StatusInternalServerError, gin.H{"success": 0})
			return
		}
	}
	context.JSON(http.StatusOK, gin.H{"success": 1, "data": todoEntity})
	return
}

func NewTodoHandler(r *gin.Engine, db *sql.DB) {

	handler := TodoHandler{
		TodoService: service.NewTodoService(repository.NewTodoMysqlRepository(db)),
	}

	r.GET("/todo", handler.fetchAll)
	r.GET("/todo/:id", handler.getTodo)
	r.POST("/todo", handler.createTodo)
	r.PUT("/todo/:id", handler.updateTodo)
	r.DELETE("/todo/:id", handler.deleteTodo)
}
