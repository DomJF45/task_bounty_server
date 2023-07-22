package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"task_bounty_server/configs"
	"task_bounty_server/models"
	"task_bounty_server/responses"
)

var projectCollection *mongo.Collection = configs.GetCollection(configs.DB, "projects")

func CreateProject(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var project models.Project
	defer cancel()

	if err := c.BodyParser(&project); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ProjectResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if validationErr := validate.Struct(&project); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ProjectResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	userId, ok := c.Locals("UserID").(string)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(responses.ProjectResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": "invalid user id"},
		})
	}

	fmt.Printf("\nuserId: %v", userId)

	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ProjectResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	project.Tasks = make([]models.Task, 0)

	newProject := models.Project{
		ID:             primitive.NewObjectID(),
		Name:           project.Name,
		ProjectManager: objectID,
		Rank:           project.Rank,
		Exp:            project.Exp,
		Tasks:          project.Tasks,
	}

	result, err := projectCollection.InsertOne(ctx, newProject)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ProjectResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.ProjectResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetUserProjects(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId, ok := c.Locals("UserID").(string)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": "invalid user id"},
		})
	}

	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": "invalid user id"},
		})
	}

	filter := bson.M{"projectManager": objectID}

	cursor, err := projectCollection.Find(ctx, filter)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ProjectResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	defer cursor.Close(ctx)

	var projects []models.Project
	for cursor.Next(ctx) {
		var project models.Project
		if err := cursor.Decode(&project); err != nil {
			return c.Status(http.StatusBadRequest).JSON(responses.ProjectResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
		projects = append(projects, project)
	}

	if err := cursor.Err(); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ProjectResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.ProjectResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": projects}})
}

func GetProjectById(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var project models.Project
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(c.Params("projectID"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ProjectResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	err = projectCollection.FindOne(ctx, bson.D{{Key: "_id", Value: objectID}}).Decode(&project)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusBadRequest).JSON(responses.ProjectResponse{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data:    &fiber.Map{"data": err.Error()},
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(responses.ProjectResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	return c.Status(http.StatusOK).JSON(responses.ProjectResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &fiber.Map{"data": project},
	})
}

func TakeProject(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(c.Params("projectID"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ProjectResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	filter := bson.D{{Key: "_id", Value: objectID}}
	userID, err := primitive.ObjectIDFromHex(c.Locals("UserID").(string))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ProjectResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}
	var user models.User
	if err := userCollection.FindOne(ctx, bson.D{{Key: "_id", Value: userID}}).Decode(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ProjectResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	fmt.Printf("\nUpdated Task User: %v", user.FirstName)
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "takenBy", Value: user.FirstName}}}}

	result, err := projectCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ProjectResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	return c.Status(http.StatusOK).JSON(responses.ProjectResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &fiber.Map{"data": result},
	})
}

func AddTask(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(c.Params("projectID"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ProjectResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ProjectResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	if validationErr := validate.Struct(&task); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ProjectResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	update := bson.D{{Key: "$push", Value: bson.D{{Key: "tasks", Value: task}}}}

	result, err := projectCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ProjectResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	return c.Status(http.StatusOK).JSON(responses.ProjectResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &fiber.Map{"data": result},
	})
}

func GetTasks(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(c.Params("projectID"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ProjectResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}
	fmt.Printf("\n\nTasks: %v", objectID)
	// Find the project with the specified ID
	var project models.Project
	err = projectCollection.FindOne(ctx, bson.D{{Key: "_id", Value: objectID}}).Decode(&project)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(responses.ProjectResponse{
				Status:  http.StatusNotFound,
				Message: "error",
				Data:    &fiber.Map{"data": "project not found"},
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(responses.ProjectResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	// Return the tasks from the project
	return c.Status(http.StatusOK).JSON(responses.ProjectResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &fiber.Map{"data": project.Tasks},
	})
}
