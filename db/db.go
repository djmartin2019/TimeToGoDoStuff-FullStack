package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID          int
	Description string
	Completed   bool
}

func InitDB(filePath string) *sql.DB {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		panic(err)
	}

	if db == nil {
		panic("db nil")
	}
	return db
}

func CreateTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS tasks (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
    "description" TEXT,
    "completed" BOOLEAN
  );`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		panic(err)
	}
}

func CreateTask(db *sql.DB, description string) {
	sqlAddTask := `INSERT INTO tasks(description, completed) VALUES (?, FALSE)`
	_, err := db.Exec(sqlAddTask, description)
	if err != nil {
		panic(err)
	}
}

func GetAllTasks(db *sql.DB) []Task {
	sqlGetAll := `SELECT * FROM tasks`
	rows, err := db.Query(sqlGetAll)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.ID, &task.Description, &task.Completed)
		if err != nil {
			panic(err)
		}
		tasks = append(tasks, task)
	}

	return tasks
}

func UpdateTask(db *sql.DB, id int, description string, completed bool) {
	_, err := db.Exec("UPDATE tasks SET description = ?, completed = ? WHERE id = ?", description, completed, id)
	if err != nil {
		panic(err)
	}
}

func DeleteTask(db *sql.DB, id string) {
	_, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		panic(err)
	}
}
