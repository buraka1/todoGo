package entity

type Todo struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Completed int `json:"completed"`
}
type Todos []*Todo

func NewTodo() *Todo  {
	return &Todo{Completed: 0}
}
