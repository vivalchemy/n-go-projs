package main

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	dbName         = "hrms"
	collectionName = "employee"
	// you can go with either of the two authSource is required else it will look
	// for the user in $dbName auth database
	mongoURL = "mongodb://hrms_admin:HrmsSecurePass123@localhost:27017/" + dbName + "/?authSource=admin"
	// mongoURL       = "mongodb://hrms_admin:HrmsSecurePass123@localhost:27017/"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var mg MongoInstance

type Employee struct {
	ID     bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string        `json:"name"`
	Salary float64       `json:"salary"`
	Age    int           `json:"age"`
}

func Connect() error {
	client, err := mongo.Connect(options.Client().ApplyURI(mongoURL))
	if err != nil {
		return err
	}
	log.Println("connecting to mongodb")

	db := client.Database(dbName)

	mg = MongoInstance{
		Client: client,
		Db:     db,
	}

	log.Println("connected to mongodb")

	return nil
}

func getEmployees(c fiber.Ctx) error {
	query := bson.D{{}}

	cursor, err := mg.Db.Collection(collectionName).Find(c.Context(), query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var employees []Employee = make([]Employee, 0)

	if err := cursor.All(c.Context(), &employees); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(employees)
}

func getEmployee(c fiber.Ctx) error {
	collection := mg.Db.Collection(collectionName)

	id := c.Params("id")
	employeeId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	filter := bson.D{{Key: "_id", Value: employeeId}}
	findResult := collection.FindOne(c.Context(), filter)
	if findResult.Err() != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": findResult.Err().Error()})
	}

	employee := Employee{}
	findResult.Decode(&employee)

	return c.Status(http.StatusOK).JSON(employee)
}

func addEmployee(c fiber.Ctx) error {
	collection := mg.Db.Collection(collectionName)

	employee := new(Employee)
	if err := c.Bind().Body(employee); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	insertionResult, err := collection.InsertOne(c.Context(), employee)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	filter := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
	createdRecord := collection.FindOne(c.Context(), filter)

	createdEmployee := Employee{}
	createdRecord.Decode(&createdEmployee)

	return c.Status(http.StatusCreated).JSON(createdEmployee)
}

func modifyEmployee(c fiber.Ctx) error {
	collection := mg.Db.Collection(collectionName)

	id := c.Params("id")
	employeeId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	employee := new(Employee)
	if err := c.Bind().Body(employee); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	query := bson.D{{Key: "_id", Value: employeeId}}
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "name", Value: employee.Name},
				{Key: "age", Value: employee.Age},
				{Key: "salary", Value: employee.Salary},
			},
		},
	}

	err = collection.FindOneAndUpdate(c.Context(), query, update).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	employee.ID = employeeId
	return c.Status(http.StatusOK).JSON(employee)
}

func deleteEmployee(c fiber.Ctx) error {
	collection := mg.Db.Collection(collectionName)

	id := c.Params("id")
	employeeId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	filter := bson.D{{Key: "_id", Value: employeeId}}
	deleteResult, err := collection.DeleteOne(c.Context(), filter)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.SendStatus(http.StatusNotFound)
		}
		return c.SendStatus(http.StatusInternalServerError)
	}

	if deleteResult.DeletedCount < 1 {
		return c.SendStatus(http.StatusNotFound)
	}

	return c.Status(http.StatusOK).JSON("record deleted")
}

func main() {
	if err := Connect(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()
	app.Get("/employee", getEmployees)
	app.Get("/employee/:id", getEmployee)
	app.Post("/employee", addEmployee)
	app.Put("/employee/:id", modifyEmployee)
	app.Delete("/employee/:id", deleteEmployee)

	app.Listen(":8080")
}
