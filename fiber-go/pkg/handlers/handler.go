package handlers

import (
	"fiber-go/pkg/database"
	"fiber-go/pkg/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	_ "log"
)

var mg *database.MongoInstance = database.GetMongoInstance()

func GetAllEmployees(c *fiber.Ctx) error {
	query := bson.M{}
	cursor, err := mg.DB.Collection("employees").Find(c.Context(), query)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	var employees []models.Employee = make([]models.Employee, 0)
	err = cursor.All(c.Context(), &employees)

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(employees)
}

func GetEmployeeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := mg.DB.Collection("employees")
	employeeID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	query := bson.M{"_id": employeeID}
	result := collection.FindOne(c.Context(), query)

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return c.SendStatus(404)
		}
		return c.Status(500).SendString(result.Err().Error())
	}
	employee := new(models.Employee)
	result.Decode(employee)

	return c.JSON(employee)
}

func CreateNewEmployee(c *fiber.Ctx) error {
	employee := new(models.Employee)
	collection := mg.DB.Collection("employees")

	if err := c.BodyParser(employee); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	employee.ID = ""
	insertionResult, err := collection.InsertOne(c.Context(), employee)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	filter := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
	createdRecord := collection.FindOne(c.Context(), filter)
	createdEmployee := &models.Employee{}
	createdRecord.Decode(createdEmployee)

	return c.JSON(createdEmployee)
}

func UpdateEmployee(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := mg.DB.Collection("employees")
	employeeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	employee := new(models.Employee)
	if err := c.BodyParser(employee); err != nil {
		return c.Status(404).SendString(err.Error())
	}
	query := bson.M{"_id": employeeID}
	updateOperators := bson.M{
		"$set": bson.M{
			"name":   employee.Name,
			"salary": employee.Salary,
			"age":    employee.Age,
		},
	}
	err = collection.FindOneAndUpdate(c.Context(), query, updateOperators).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.SendStatus(404)
		}
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(employee)
}

func DeleteEmployee(c *fiber.Ctx) error {
	collection := mg.DB.Collection("employees")
	employeeID, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	query := bson.M{"_id": employeeID}
	result, err := collection.DeleteOne(c.Context(), query)

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	if result.DeletedCount < 1 {
		return c.SendStatus(404)
	}
	return c.SendStatus(200)
}
