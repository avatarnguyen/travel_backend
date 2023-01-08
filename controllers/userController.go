package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/avatarnguyen/travel_backend/database"
	"github.com/avatarnguyen/travel_backend/helpers"
	"github.com/avatarnguyen/travel_backend/models"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
)

var (
	userCollection = database.OpenConnection(database.Client, "user")
	validate       = validator.New()
)

func SignUp() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		user := new(models.User)
		if err := c.BodyParser(user); err != nil {
			cancel()
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		validateErr := validate.Struct(user)
		if validateErr != nil {
			cancel()
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": validateErr.Error(),
			})
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()

		if err != nil {
			log.Panic(err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "error occured while checking email",
			})
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "error occured while checking for phone",
			})
		}

		if count > 0 {
			log.Print("Documents Count: $count")
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "this email or phone number exists",
			})
		}
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		token, refreshToken, _ := helpers.GenerateAllToken(
			*user.Email,
			*user.First_name,
			*user.Last_name,
			*user.User_type,
			user.User_id,
		)
		user.Token = &token
		user.Refresh_token = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintln("User item was not created: " + insertErr.Error())
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": msg,
			})
		}

		defer cancel()
		return c.Status(http.StatusOK).JSON(resultInsertionNumber)
	}
}

func Login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var foundUser models.User

		user := new(models.User)

		if err := c.BodyParser(user); err != nil {
			return c.Status(http.StatusBadRequest).JSONP(fiber.Map{
				"error": err.Error(),
			})
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSONP(fiber.Map{
				"error": "email or password is incorrect",
			})
		}

		passwordIsValid, msg := VerifyPassword(
			*user.Password,
			*foundUser.Password,
		)
		defer cancel()

		if !passwordIsValid {
			return c.Status(http.StatusInternalServerError).JSONP(fiber.Map{
				"error": msg,
			})
		}

		if foundUser.Email == nil {
			return c.Status(http.StatusInternalServerError).JSONP(fiber.Map{
				"error": "user not found",
			})
		}

		token, refreshToken, _ := helpers.GenerateAllToken(
			*foundUser.Email,
			*foundUser.First_name,
			*foundUser.Last_name,
			*foundUser.User_type,
			foundUser.User_id,
		)
		helpers.UpdateAllToken(token, refreshToken, foundUser.User_id)

		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSONP(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(200).JSON(foundUser)
	}
}

func GetUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Params("user_id")

		fmt.Println("User Id : " + userId)
		if err := helpers.MatchUserTypeToUuid(c, userId); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)

		defer cancel()

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(200).JSON(user)
	}
}

func GetUsers() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		groupStage := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
			{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
		}}}
		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}},
			}},
		}

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage,
		})

		fmt.Println("Result: ", result)
		defer cancel()

		if err != nil {
			fmt.Println("Error: ", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "error occurred while listening user items",
			})
		}

		var allUsers []bson.M

		if err = result.All(ctx, &allUsers); err != nil {
			log.Fatal(err)
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"result": allUsers[0],
		})
	}
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic()
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintln("email of password is incorrect")
		check = false
	}

	return check, msg
}
