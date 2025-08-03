package database

import (
	"database/sql"
	"fmt"
	"log"
	"todo/server/config"

	_ "github.com/go-sql-driver/mysql"
)

func createTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS todos (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		completed BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP

	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
}

//why pass a pointer not a value
//duplicating and resetting internal state for every function call.
//*sql.DB lets you share the same connection pool across your whole app.

func Connect(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
        cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	//dsn(Data Source Name) It’s a formatted connection string used to describe how to connect to your MySQL database.

	log.Println(dsn)
	db, err := sql.Open("mysql", dsn) //creates a sql.DB struct(a wrapper)
	//holds driver info, configs like max opens/idle connections, a pool of connections(lazy conncection managemenet)
	if err != nil{
		return nil,err
	}

	fmt.Println(db)
	createTable(db)

	if err := db.Ping(); err != nil{
		return nil, err
	}

	return db, nil
	//returns a pointer of the type *sql.DB which is safe for concurrent use
	//You can share it across goroutines.
  	//Under the hood, it uses mutexes and channels to manage access.
}
