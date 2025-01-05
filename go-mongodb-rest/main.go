package main

import (
	"context"
	"log"
	"net/http"

	"github.com/vivalchemy/n-go-projs/go-mongodb-rest/controllers"
	"github.com/vivalchemy/vhttp"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	r := vhttp.NewServeMux()

	// create a mongodb client
	client := getSession()
	uc := controllers.NewUserController(getSession())
	defer client.Disconnect(context.Background())

	// register route handlers
	r.FuncGet("/user/{id}", uc.GetUser)
	r.FuncPost("/user", uc.CreateUser)
	r.FuncDelete("/user/{id}", uc.DeleteUser)

	// start the server
	log.Println("Listening on port 8080")
	http.ListenAndServe(":8080", r)
}

// getSession returns a mongodb client
func getSession() *mongo.Client {
	session, err := mongo.Connect(options.Client().ApplyURI("mongodb://root:toor@localhost:27017"))
	if err != nil {
		panic(err)
	}
	return session
}
