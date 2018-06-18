package db

import (
	_ "database/sql"
	_ "fmt"
	_ "github.com/lib/pq"
	_ "time"
	"fmt"
	"database/sql"
	"testAPI/api/swagger/models"
)

const(
	DB_USER = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME = "test"
)

func getDbConn() (db *sql.DB){
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD,DB_NAME)
	db,err := sql.Open("postgres",dbinfo)
	checkErr(err)
	defer db.Close()
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

func main(){
	ans := allUsers()
	for item := range ans {
		fmt.Println("%3v | %8v ",item,ans[item].Name)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}


