// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	graceful "github.com/tylerb/graceful"

	"basicAPI/api/swagger/restapi/operations"
	"basicAPI/api/swagger/restapi/operations/users"
	"fmt"
	"sync"
	"sync/atomic"
	"basicAPI/api/swagger/models"
	"github.com/go-openapi/swag"
)

//go:generate swagger generate server --target .. --name user-list --spec ../swagger.yml


// Code to add and delete users

var usersList = make(map[int64]*models.User)
var lastID int64

var userLock = &sync.Mutex{}

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
}

func deleteItem(id int64) error {
	userLock.Lock()
	defer userLock.Unlock()

	fmt.Println(id)
	_, exists := usersList[id]
	if !exists {
		return errors.NotFound("not found: item %d", id)
	}

	delete(usersList, id)
	return nil
}

func allUsers() (result []*models.User) {

	result = make([]*models.User, 0)

	for item := range usersList {
		result = append(result, usersList[item])
	}
	return
	//  return allUsers()
}

func getSingleUser(id int64) (result *models.User, err error) {

	_, exists := usersList[id]
	if !exists {
		return nil, errors.NotFound("User not found %d",id);
	}
	result = usersList[id]
	return result,nil
}

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
	api.UsersAddOneHandler = users.AddOneHandlerFunc(func(params users.AddOneParams) middleware.Responder{
		if err := addUser(params.Body); err != nil {
			return users.NewAddOneDefault(500).WithPayload(&models.Error{Code:500, Message: swag.String(err.Error())})
		}
		return users.NewAddOneCreated().WithPayload(params.Body)
	})


	// CURL Commands: curl -i localhost:35577 - GET
	api.UsersFindUserHandler = users.FindUserHandlerFunc(func(params users.FindUserParams) middleware.Responder {
		return users.NewFindUserOK().WithPayload(allUsers())
	})


	// - POST
	/*
	api.UsersAddOneHandler = users.AddOneHandlerFunc(func(params users.AddOneParams) middleware.Responder{
		return users.NewFindUserOK().WithPayload(allUsers())
	}) */

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
