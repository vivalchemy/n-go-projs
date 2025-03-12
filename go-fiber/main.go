package main

import (
	"log"
	"vivalchemy/go-fiber/database"
	"vivalchemy/go-fiber/lead"

	"github.com/gofiber/fiber"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRoutes(app *fiber.App) {
	app.Get("/api/v1/lead", lead.GetLeads)
	app.Get("/api/v1/lead/:id", lead.GetLead)
	app.Post("/api/v1/lead", lead.PostLeads)
	app.Delete("/api/v1/lead/:id", lead.DeleteLeads)
}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open(sqlite.Open("leads.db"))
	if err != nil {
		log.Fatal("Failed to connect to db", err)
	}
	log.Println("Connected to db leads.db")
	database.DBConn.AutoMigrate(&lead.Lead{})
	log.Println("Database migrated")
}

func main() {
	app := fiber.New()

	initDatabase()
	setupRoutes(app)

	log.Fatal(app.Listen(":8080"))
}
