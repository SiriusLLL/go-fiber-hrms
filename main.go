package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	DB     *mongo.Database
}

var mg MongoInstance

const dbName = "fiber-hrms"
const mongoURI = "mongodb://localhost:27017" + dbName

type Employee struct {
	ID     string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string  `json:"name"`
	Salary float64 `json:"salary"`
	Age    float64 `json:"age"`
}

func Connect() error {

	// // NewClient() 已被弃用，使用Connect()直接创建并连接mongoDB客户端

	// client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	// if err != nil {
	// 	return err
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()
	// err1 := client.Connect(ctx)
	// if err1 != nil {
	// 	return err1
	// }

	// 创建并连接mongoDB客户端
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(dbName)
	mg = MongoInstance{
		Client: client,
		DB:     db,
	}
	return nil
}

func GetEmployee(c *fiber.Ctx) error {
	query := bson.D{{}}
	cursor, err := mg.DB.Collection("employees").Find(c.Context(), query)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	var employees []Employee = make([]Employee, 0)

	if err := cursor.All(c.Context(), &employees); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(employees)
}

func CreateEmployee(c *fiber.Ctx) error {
	collection := mg.DB.Collection("employee")
	employee := new(Employee)

	if err := c.BodyParser(employee); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	employee.ID = ""
	insertionResult, err := collection.InsertOne(c.Context(), employee)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	filter := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
	createRecord := collection.FindOne(c.Context(), filter)
	createEmployee := &Employee{}
	createRecord.Decode(createEmployee)

	return c.Status(201).JSON(createEmployee)
}

func main() {

	if err := Connect(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/employee", GetEmployee)

	app.Post("/employee", CreateEmployee)
	app.Put("/employee/:id")
	app.Delete("/employee/:id")
}