package routes

import (
	"github.com/abhishek0chauhan/golang-jwt-project/controllers"
	"github.com/abhishek0chauhan/golang-jwt-project/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(route *gin.Engine){
	route.Use(middlewares.Authenticate())
	route.GET("/users", controllers.GetUsers())
	route.GET("/users/:user_id", controllers.GetUser())
}