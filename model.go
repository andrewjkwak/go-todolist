package main

import "database/sql"

type todo struct {
	ID        int    `json:"id"`
	Todo      string `json:"todo"`
	Completed bool   `json:"completed"`
}

func (t *todo) getTodo(db *sql.DB) error {
	return db.QueryRow("SELECT todo, completed FROM todos WHERE id = $1", t.ID).Scan(&t.Todo, &t.Completed)
}

func (t *todo) updateTodo(db *sql.DB) error {
	_, err := db.Exec("UPDATE todos SET todo = $1, completed = $2 WHERE id = $3", t.Todo, t.Completed, t.ID)
	return err
}

func (t *todo) deleteTodo(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM todos WHERE id = $1", t.ID)
	return err
}

func (t *todo) addTodo(db *sql.DB) error {
	err := db.QueryRow("INSERT INTO todos(todo) VALUES($1) RETURNING id", t.Todo).Scan(&t.ID)
	if err != nil {
		return err
	}
	return nil
}

func getTodos(db *sql.DB) ([]todo, error) {
	rows, err := db.Query("SELECT id, todo, completed FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []todo{}
	for rows.Next() {
		var t todo
		if err := rows.Scan(&t.ID, &t.Todo, &t.Completed); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, nil
}
