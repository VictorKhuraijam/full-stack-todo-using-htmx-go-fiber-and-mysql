package services

import (
	"database/sql"
	// "fmt"
	"todo/server/models"
	// "log"
)

type TodoService struct {
	db *sql.DB
}

//constructor function - in go it’s idiomatic to name constructor functions as NewTypeName, like NewTodoService, NewUserRepo, NewAppConfig, etc.
//to create and initialize a new instance of a struct.
func NewTodoService(db *sql.DB) *TodoService  {
	return &TodoService{db: db} // returns a pointer to the struct todoService
}

func (s *TodoService) GetAll() ([]models.Todo, error)  {
	query := "SELECT * FROM todos ORDER BY created_at DESC"
    rows, err := s.db.Query(query)
	// using db.Query(), Go opens a database cursor and keeps a connection open to allow you to read results.
	//opens multiple rows
    if err != nil {
        return nil, err
    }

    defer rows.Close() //Releases DB connection + memory

	var todos []models.Todo //Declares a slice to hold all the fetched todos.
	for rows.Next(){ //moves the cursor to the next row in the result set.
		var todo models.Todo //Creates a new Todo struct to store the current row.
		 err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)//maps column values into the struct fields.
		 //& means we are passing a pointer to each field so Scan can populate it.
        if err != nil {
            return nil, err
        }
		// log.Printf("Fetched Todo:%+v\n", todo)
		todos = append(todos, todo)
	}
	return todos, nil
}

func (s *TodoService) GetByID(id int) (*models.Todo, error) {
	query := "SELECT id, title, completed, created_at, updated_at FROM todos WHERE id = ?"
    row := s.db.QueryRow(query, id)//used to fetch single row

	var todo models.Todo
	err := row.Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
    if err != nil {
        return nil, err
    }

	return &todo, nil
}

func (s *TodoService) Create(title string) (*models.Todo, error){
	query := "INSERT INTO todos (title) VALUES (?)"
	result, err := s.db.Exec(query, title)

	if err != nil{
		return nil,err
	}

	id, err := result.LastInsertId()
	if err != nil{
		return nil, err
	}
	return s.GetByID(int(id))
}

func (s *TodoService) ToggleComplete(id int) (*models.Todo, error) {
	query := "UPDATE todos SET completed = NOT completed WHERE id = ?"
	_, err := s.db.Exec(query, id)
	if err != nil{
		return nil, err
	}
	return s.GetByID(id)
}

func (s *TodoService) UpdateTodo(id int, title string, completed bool) (*models.Todo, error){
	query := "UPDATE todos SET title = ?, completed = ?, updated_at = CURRENT_TIMESTAMP WHERE ID = ?"
	_, err := s.db.Exec(query, title, completed,id)
	if err != nil{
		return nil, err
	}
	return s.GetByID(id)
}

func (s *TodoService) Delete(id int) error  {
	query := "DELETE FROM todos WHERE id =?"
	_,err := s.db.Exec(query, id)

	return err
}
