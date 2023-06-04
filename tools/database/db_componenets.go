package database

import (
	"database/sql"
	"fmt"
	"netBlast/pkg"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Opens a new database connection
func OpenDB() *sql.DB {
	if _, ok := os.LookupEnv("DBUSER"); ok == false {
		return nil
	}
	if _, ok := os.LookupEnv("DBPASS"); ok == false {
		return nil
	}

	// dataSourceName := mysql.Config{
	// 	User:   os.Getenv("DBUSER"),
	// 	Passwd: os.Getenv("DBPASS"),
	// 	Net:    net,
	// 	Addr:   ip + ":" + port,
	// 	DBName: dBName,
	// }

	db, err := sql.Open(driver, dataSourceName)
	if err != nil {
		pkg.LogError(err)
		fmt.Println("Couldnt open database.")
		return nil
	} else if !PingDB(db) {
		pkg.LogError(err)
		fmt.Println("Couldnt ping database.")
		return nil
	}

	return db
}

// Checks database connection
func PingDB(db *sql.DB) bool {
	err := db.Ping()
	if err != nil {
		pkg.LogError(err)
		fmt.Println("bad ping to DB")
		return false
	}
	return true
}

// Inserts new user into the database
func InsertDBUser(db *sql.DB, user pkg.User) {
	if !PingDB(db) {
		return
	}

	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (id INT PRIMARY KEY AUTO_INCREMENT NOT NULL, name VARCHAR(255) UNIQUE NOT NULL, color VARCHAR(255) NOT NULL, status VARCHAR(255) NOT NULL)")
	if err != nil {
		pkg.LogError(err)
		panic(err.Error())
	}

	_, err = db.Exec("INSERT INTO users(name, color, status) VALUES(?, ?, ?)", user.Name, user.UserColor, "Offline")
	if err != nil {
		pkg.LogError(err)
		panic(err.Error())
	}
}

// Deletes user from database
func DeleteDBUser(db *sql.DB, username string) {
	if !PingDB(db) {
		return
	}

	_, err := db.Exec("DELETE FROM users WHERE name = ?", username)
	if err != nil {
		pkg.LogError(err)
		panic(err.Error())
	}
	UpdateDBUserInfo(db, "status", "name", "Offline", username)
}

// Updates user info
func UpdateDBUserInfo(db *sql.DB, info string, arg string, infoOption any, argOption any) {
	if !PingDB(db) {
		return
	}

	var err error

	_, err = db.Exec("UPDATE users SET "+info+" = ? WHERE "+arg+" = ?", infoOption, argOption)

	if err != nil {
		pkg.LogError(err)
		panic(err.Error())
	}
}

// FInds user info
func FindDBUserInfo(db *sql.DB, info string, arg string, option any) string {
	if !PingDB(db) {
		return ""
	}
	var data string
	var id int
	var err error

	result := db.QueryRow("SELECT "+info+" FROM users WHERE "+arg+" = ? LIMIT 1", option)

	if info == "id" {
		err = result.Scan(&id)
	} else {
		err = result.Scan(&data)
	}

	if err != nil {
		return ""
	}

	return fmt.Sprintf("%+v", data)
}

// Performs specific queries to the database
func QueryDB(db *sql.DB, query string) string {
	if !PingDB(db) {
		return ""
	}
	var data any
	err := db.QueryRow(query).Scan(&data)

	if err != nil || data == nil {
		return ""
	}
	var returnData any

	pkg.ParseFromJson(data.([]byte), &returnData, "baddddd")

	return fmt.Sprintf("%+v", returnData)
}

// Check if table is available
func CheckDBTable(table string, db *sql.DB) bool {
	_, err := db.Query("SELECT * FROM " + table + ";")
	if err != nil {
		pkg.LogError(err)
		return false
	}
	return true
}

// Creates new table
func CreateDBTable(table string, db *sql.DB) {
	_, err := db.Query("CREATE TABLE " + table + ";")
	if err != nil {
		pkg.LogError(err)
		panic(err.Error())
	}
}
