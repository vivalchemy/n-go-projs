package lead

import (
	"vivalchemy/go-fiber/database"

	"github.com/gofiber/fiber"
	"gorm.io/gorm"
)

type Lead struct {
	gorm.Model
	Name    string `json:"name"`
	Company string `json:"company"`
	Email   string `json:"email"`
	Phone   int    `json:"phone"`
}

func GetLeads(c *fiber.Ctx) {
	db := database.DBConn
	var leads []Lead
	db.Find(&leads)
	c.JSON(leads)
}

func GetLead(c *fiber.Ctx) {
	id := c.Params("id")
	db := database.DBConn
	var lead Lead
	db.First(&lead, id)
	c.JSON(lead)
}

func PostLeads(c *fiber.Ctx) {
	db := database.DBConn
	lead := new(Lead)
	if err := c.BodyParser(lead); err != nil {
		c.Status(400).JSON(fiber.Map{"error": err.Error()})
		return
	}
	db.Create(&lead)
	c.JSON(lead)
}

func DeleteLeads(c *fiber.Ctx) {
	id := c.Params("id")
	db := database.DBConn

	var lead Lead
	db.First(&lead, id)
	if lead.Name == "" {
		c.Status(404).JSON(fiber.Map{"message": "Lead not found"})
		return
	}
	db.Delete(&lead)
	c.Status(200).JSON(fiber.Map{"message": "Lead deleted successfully"})
}
