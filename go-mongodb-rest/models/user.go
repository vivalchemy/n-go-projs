package models

import "go.mongodb.org/mongo-driver/v2/bson"

type User struct {
	Id     bson.ObjectID `json:"id" bson:"_id"`
	Name   string        `json:"name" bson:"name"`
	Gender string        `json:"gender" bson:"gender"`
	Age    int           `json:"age" bson:"age"`
}
