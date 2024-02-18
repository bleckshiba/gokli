package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type TaskRepo struct {
	DB *sql.DB
}

var db, user, pass, host string
var connStr string

func createInstance() TaskRepo {
	return TaskRepo{}
}

func (repo *TaskRepo) Init() {
	initVariables()
	connectDB()
}

func initVariables() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln(err)
	}

	if db = os.Getenv("DBNAME"); db == "" {
		log.Fatalln("DBNAME not found")
	}

	if user = os.Getenv("DBUSER"); user == "" {
		log.Fatalln("DBUSER not found")
	}

	if pass = os.Getenv("DBPASS"); pass == "" {
		log.Fatalln("DBPASS not found")
	}

	if host = os.Getenv("DBHOST"); host == "" {
		log.Fatalln("DBHOST not found")
	}

	connStr = fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		user,
		pass,
		host,
		db,
	)

}

func connectDB() {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	// defer db.Close()

	taskRepo.DB = db
}

func (repo *TaskRepo) FetchTasks() []Task {
	rows, err := repo.DB.Query("SELECT * FROM tasks")
	if err != nil {
		log.Fatalf("err: %s\n", err)
		log.Fatalln(err)
		os.Exit(1)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.title, &task.Desc, &task.Status)
		if err != nil {
			log.Fatalf("err: %s\n", err)
		}
		tasks = append(tasks, task)
	}
	return tasks
}

func (repo *TaskRepo) InsertTask(task Task) {
	query := "INSERT INTO tasks (id, title, description, status) VALUES ($1, $2, $3, $4)"
	_, err := repo.DB.Exec(query, task.ID, task.title, task.Desc, task.Status)
	if err != nil {
		log.Fatalln(err)
	}
}

func (repo *TaskRepo) DeleteTask(task Task) {
	query := "DELETE FROM tasks WHERE id = $1"
	_, err := repo.DB.Exec(query, task.ID)
	if err != nil {
		log.Fatalln(err)
	}
}

func (repo *TaskRepo) UpdateTask(task Task) {
	query := "UPDATE tasks SET title = $1, description = $2, status = $3 WHERE id = $4"
	_, err := repo.DB.Exec(query, task.title, task.Desc, task.Status, task.ID)
	if err != nil {
		log.Fatalln(err)
	}
}
