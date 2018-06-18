package main

import (
	_ "fmt"
	_ "github.com/lib/pq"
	_ "time"
	"fmt"
	"database/sql"
	"testAPI/api/swagger/models"
)

const(
	DB_HOST = "localhost"
	DB_PORT = 5432
	DB_USER = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME = "test"
)

func getDbConn() (db *sql.DB){
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DB_HOST,DB_PORT,DB_USER, DB_PASSWORD,DB_NAME)
	db,err := sql.Open("postgres",dbinfo)
	checkErr(err)
	return db
}

func allUsers() (result map[int64]*models.User) {
	result = make(map[int64]*models.User, 0)

	// result = append(result, usersList[item])
	dbConn := getDbConn()

	rows, err := dbConn.Query("select uid,name from userinfo")
	checkErr(err)

	var userList = make(map[int64]*models.User)

	for rows.Next(){
		var uid int64
		var name string
		err = rows.Scan(&uid, &name)
		checkErr(err)
		userList[uid] = &models.User{ID:uid,Name:&name}
	}
	return userList
}

func getSingleUserDB(id int64) (result *models.User, err error) {

	// result = append(result, usersList[item])
	dbConn := getDbConn()

	var name string

	err = dbConn.QueryRow("select name from userinfo where uid=$1", id).Scan(&name)

	checkErr(err)
	return &models.User{ID:id,Name:&name},err
}

func main(){
	ans,err := getSingleUserDB(5)
	fmt.Println(ans)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}


