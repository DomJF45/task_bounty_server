package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"task_bounty_server/configs"
	"task_bounty_server/models"
	"task_bounty_server/responses"
)

func GenerateToken(userID string) (string, error) {
	claims := models.Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(configs.EnvJWTSecret())
}

func VerifyToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected sigining method: %v", token.Header["alg"])
		}

		return []byte(configs.EnvJWTSecret()), nil
	})
	if err != nil {
		return "", fmt.Errorf("invalid token: %v", err.Error())
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		fmt.Printf("\nClaims: %v", claims)
		return claims.UserID, nil
	}

	return "", fmt.Errorf("invalid token")
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func comparePasswords(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

var (
	userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
	validate                         = validator.New()
)

func RegisterUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User // want either ProjectManager or IndividualContributor depending on params
	defer cancel()

	// validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	userHasEmail := userCollection.FindOne(ctx, bson.D{{Key: "email", Value: user.Email}})

	var existingUser models.User
	if err := userHasEmail.Decode(&existingUser); err == nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": "user emial already in use"}})
	}

	newUser := &models.User{
		ID:        primitive.NewObjectID(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  hashedPassword,
		Title:     user.Title,
	}

	if _, err := userCollection.InsertOne(ctx, newUser); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	userID := primitive.ObjectID.Hex(newUser.ID)

	token, err := GenerateToken(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	fmt.Printf("NewUserID: %v", newUser.ID)

	return c.Status(http.StatusCreated).JSON(responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"token": token}})
}

func LoginUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.LoginUser
	defer cancel()

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	var requestedUser models.User
	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&requestedUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if !comparePasswords(user.Password, requestedUser.Password) {
		return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &fiber.Map{"token": "Incorrect Password"}})
	}

	userID := primitive.ObjectID.Hex(requestedUser.ID)

	token, err := GenerateToken(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	fmt.Printf("Found User 🔥")
	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"token": token}})
}

func GetUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	userID, ok := c.Locals("UserID").(string)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": "invalid user id"},
		})
	}

	fmt.Printf("\nuserId: %v", userID)

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	if err := userCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &fiber.Map{"data": user},
	})
}
