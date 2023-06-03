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

	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (id INT PRIMARY KEY AUTO_INCREMENT NOT NULL, name VARCHAR(255) UNIQUE NOT NULL, color VARCHAR(255) NOT NULL, status VARCHAR(255) NOT NULL, conn JSON)")
	if err != nil {
		pkg.LogError(err)
		panic(err.Error())
	}

	_, err = db.Exec("INSERT INTO users(name, color, status) VALUES(?, ?, ?)", user.Name, user.UserColor, "offline")
	if err != nil {
		pkg.LogError(err)
		panic(err.Error())
	}
}

// Deletes user from database
func DeleteDBUser(db *sql.DB, username string) {
	_, err := db.Exec("DELETE FROM users WHERE name = ?", username)
	if err != nil {
		pkg.LogError(err)
		panic(err.Error())
	}
	UpdateDBUserInfo(db, username, "status", "Offline")
}

// Updates user info
func UpdateDBUserInfo(db *sql.DB, username string, arg string, newOption any) {
	var err error
	switch arg {
	case "username":
		_, err = db.Exec("UPDATE users SET name = ? WHERE name = ?", newOption, username)
	case "color":
		_, err = db.Exec("UPDATE users SET color = ? WHERE name = ?", newOption, username)
	case "status":
		_, err = db.Exec("UPDATE users SET status = ? WHERE name = ?", newOption, username)
	case "conn":
		conn := pkg.ParseToJson(newOption, "Couldnt parse to json")
		fmt.Println(conn)
		_, err = db.Exec("UPDATE users SET conn = ? WHERE name = ?", conn, username)
		UpdateDBUserInfo(db, username, "status", "Online")
	case "password":
		_, err = db.Exec("UPDATE users SET password = ? WHERE name = ?", newOption, username)
	default:
		return
	}

	if err != nil {
		pkg.LogError(err)
		panic(err.Error())
	}
}

func FindDBUserInfo(db *sql.DB, info string, arg string, option any) string {
	var data string
	var err error

	if arg == "conn" {
		option = pkg.ParseToJson(option, "Couldnt parse to json")
		fmt.Println(option)
	}

	err = db.QueryRow("SELECT "+info+" FROM users WHERE "+arg+" = ?", option).Scan(&data)

	if err != nil {
		pkg.LogError(err)
		fmt.Println(err.Error())
		return ""
	}
	return data
}

func CheckDBUserConn(db *sql.DB, username string) bool {
	var conn any

	result := db.QueryRow("SELECT conn FROM users WHERE name = ? LIMIT 1", username)
	err := result.Scan(&conn)
	if err != nil {
		pkg.LogError(err)
		fmt.Println(err.Error())
		return false
	}
	if conn == nil {
		return false
	}
	return true
}

// Fetches user password
func GetDBUserpass(db *sql.DB, username string) string {
	var password string

	row := db.QueryRow("SELECT name FROM users WHERE name = ? LIMIT 1", username)
	err := row.Scan(&password)
	if err != nil {
		pkg.LogError(err)
		return ""
	}
	return password
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
