// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/tylerb/graceful"

	"basicAPI/api/swagger/restapi/operations"
	"basicAPI/api/swagger/restapi/operations/users"
	"fmt"
	"sync"
	"basicAPI/api/swagger/models"
	"github.com/go-openapi/swag"
	_ "database/sql"
	_ "fmt"
	_ "github.com/lib/pq"
	_ "time"
	"database/sql"
)

//go:generate swagger generate server --target .. --name user-list --spec ../swagger.yml


// Code to add and delete users

var usersList = make(map[int64]*models.User)
var lastID int64

var userLock = &sync.Mutex{}

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


/* For Non DB- insert
func newUserId() int64 {
	return atomic.AddInt64(&lastID, 1)
}


func addUser(user *models.User) error {
	if user == nil {
		return errors.New(500, "user must be present")
	}

	userLock.Lock()
	defer userLock.Unlock()

	newID := newUserId()
	user.ID = newID
	usersList[newID] = user

	return nil
}
*/

/*
func updateItem(id int64, user *models.User) error {
	if user == nil {
		return errors.New(500, "item must be present")
	}
	userLock.Lock()
	defer userLock.Unlock()

	_, exists := usersList[id]
	if !exists {
		return errors.NotFound("not found: item %d", id)
	}

	user.ID = id
	usersList[id] = user
	return nil
} */

func updateItem(id int64, user *models.User) error {
	if user == nil {
		return errors.New(500, "item must be present")
	}

	updateID := id
	updateName := user.Name

	dbConn := getDbConn()

	stmt,err := dbConn.Prepare("update userinfo set name=$1 where uid=$2")
	res,err := stmt.Exec(updateName,updateID)
	checkErr(err)
	affect,err := res.RowsAffected()
	fmt.Println(affect,"rows affected")
	checkErr(err)

	return nil
}


/* func deleteItem(id int64) error {
	userLock.Lock()
	defer userLock.Unlock()

	fmt.Println(id)
	_, exists := usersList[id]
	if !exists {
		return errors.NotFound("not found: item %d", id)
	}

	delete(usersList, id)
	return nil
} */

func deleteItem(id int64) error {

	dbConn := getDbConn()

	stmt,err := dbConn.Prepare("delete from userinfo where uid=$1")
	res,err := stmt.Exec(id)
	checkErr(err)
	affect,err := res.RowsAffected()
	fmt.Println(affect,"rows affected")
	checkErr(err)

	return nil
}

func allUsers() (result []*models.User) {

	usersList = allUsersDB()

	result = make([]*models.User, 0)

	for item := range usersList {
		result = append(result, usersList[item])
	}
	return
}

func allUsersDB() (result map[int64]*models.User) {
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

func getSingleUser(id int64) (result *models.User, err error) {
	return getSingleUserDB(id)
}

/*
func getSingleUser(id int64) (result *models.User, err error) {

	_, exists := usersList[id]
	if !exists {
		return nil, errors.NotFound("User not found %d",id);
	}
	result = usersList[id]
	return result,nil
}*/

// end of code to add and delete users


func configureFlags(api *operations.UserListAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.UserListAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// curl -i localhost:35577 -d "{\"name\":\"message $RANDOM\"}" -H 'Content-Type: application/json'
	// Not needed - Using post as add one user to the list
	/*
	api.UsersAddOneHandler = users.AddOneHandlerFunc(func(params users.AddOneParams) middleware.Responder{
		if err := addUser(params.Body); err != nil {
			return users.NewAddOneDefault(500).WithPayload(&models.Error{Code:500, Message: swag.String(err.Error())})
		}
		return users.NewAddOneCreated().WithPayload(params.Body)
	})
	*/


	// CURL Commands: curl -i localhost:35577 - GET
	api.UsersFindUserHandler = users.FindUserHandlerFunc(func(params users.FindUserParams) middleware.Responder {
		return users.NewFindUserOK().WithPayload(allUsers())
	})


	// getall users - POST
	api.UsersAddOneHandler = users.AddOneHandlerFunc(func(params users.AddOneParams) middleware.Responder{
		return users.NewFindUserOK().WithPayload(allUsers())
	})

	// Get details of single user - GET {user}
	api.UsersGetSingleUserHandler = users.GetSingleUserHandlerFunc(func(params users.GetSingleUserParams) middleware.Responder{
		var result *models.User
		result, error := getSingleUser(params.ID)
		if error != nil{
			return users.NewGetSingleUserDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(error.Error())})
		}
		return users.NewGetSingleUserOK().WithPayload(result)
	})

	// DELETE {user}
	api.UsersDeleteUserHandler = users.DeleteUserHandlerFunc(func(params users.DeleteUserParams) middleware.Responder{
		if  err := deleteItem(params.ID); err != nil{
			return users.NewDeleteUserDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		}
		return users.NewDeleteUserNoContent()
	})

	// PUT {user}
	api.UsersUpdateUserHandler = users.UpdateUserHandlerFunc(func(params users.UpdateUserParams) middleware.Responder{
		if  err := updateItem(params.ID, params.Updateid); err != nil{
			return users.NewUpdateUserDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		}
		return users.NewUpdateUserOK().WithPayload(params.Updateid)
	})


	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}