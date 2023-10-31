package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/abhishek0chauhan/golang-jwt-project/database"
	"github.com/abhishek0chauhan/golang-jwt-project/helpers"
	"github.com/abhishek0chauhan/golang-jwt-project/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword() {

}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprint("Password is incorrect")
	}

	return check, msg
}

func Signup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		eamilCount, eamilCountErr := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if eamilCountErr != nil {
			log.Panic(eamilCountErr)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the eamil count"})
		}

		phoneCount, phoneCountErr := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if phoneCountErr != nil {
			log.Panic(phoneCountErr)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone count"})
		}

		if eamilCount > 0 && phoneCount > 0 {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone already exist"})
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = string(user.ID.Hex())
		token, refreshToken, _ := helpers.GenerateAllTokens(*&user.Email, *&user.First_name, *&user.Last_name, *&user.User_type, *&user.User_id)
		user.Token = token
		user.Refresh_token = refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(c, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()

		ctx.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(c, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "eamil or password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*&user.Password, *&foundUser.Password)
		defer cancel()

		if msg != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}
	}
}

func GetUsers() {

}

func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("user_id")

		if err := helpers.MatchUserTypeToUid(ctx, userId); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := userCollection.FindOne(c, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, user)
	}
}
