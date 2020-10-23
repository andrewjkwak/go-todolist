# Go Todolist API

Brushing up on what I learned, decided to create an REST API for a Todolist.
Implemented tests for every request as well.

Used gorilla/mux for Router, PostgreSQL for database, and godotenv to add environmental variables.

getTodos - retrieves all the todos from the database

getTodo - retrieves a single todo from the database given an id

addTodo - adds a todo into the database

updateTodo - updates a todo from the database given an id

deleteTodo - deletes a todo from the database given an id

