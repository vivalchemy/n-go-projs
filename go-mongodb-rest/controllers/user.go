package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vivalchemy/n-go-projs/go-mongodb-rest/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserController struct {
	session *mongo.Client
}

func NewUserController(session *mongo.Client) *UserController {
	return &UserController{session}
}

// The the user with the id mentionend in the path paramenter
func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("oid", oid)
	u := models.User{}
	fmt.Println("u", u)

	err = uc.session.Database("mongo-golang").Collection("users").FindOne(context.Background(), bson.M{"_id": oid}).Decode(&u)
	fmt.Println("u", u)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	userJSON, err := json.Marshal(u)
	fmt.Println("userJSON", userJSON)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(userJSON) since json.Marshall is done above no need to encode it again
	w.Write(userJSON)
}

// create a new user and store in db
func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	u := models.User{}
	json.NewDecoder(r.Body).Decode(&u)
	fmt.Println(u)
	// u.Id = bson.ObjectId(bson.NewObjectIDFromTimestamp(time.Now()))
	u.Id = bson.NewObjectID()
	fmt.Println(u)
	insertResult, err := uc.session.Database("mongo-golang").Collection("users").InsertOne(context.Background(), u)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(insertResult)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"userId":"` + u.Id.Hex() + `","message":"Created the user with id: ` + u.Id.Hex() + `"}`))
}

// delete the first user with the matching id in the db
func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = uc.session.Database("mongo-golang").Collection("users").DeleteOne(r.Context(), bson.M{"_id": oid})
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"message":"Delete the user with id: `+id+`", "userId":{`+id+`}`)
}
