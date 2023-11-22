package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/BansosPlus/bansosplus-backend.git/database"
	"github.com/BansosPlus/bansosplus-backend.git/handler"
	authMw "github.com/BansosPlus/bansosplus-backend.git/middleware"
	"github.com/BansosPlus/bansosplus-backend.git/model"
	"github.com/BansosPlus/bansosplus-backend.git/repository"
	"github.com/BansosPlus/bansosplus-backend.git/utility"
)

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the token from the request headers
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"code":    http.StatusUnauthorized,
				"status":  "error",
				"message": "Unauthorized",
			})
		}

		// Verify the token
		const bearerPrefix = "Bearer "
		tokenString := authHeader[len(bearerPrefix):]
		token, err := authMw.VerifyToken(tokenString)
		log.WithFields(log.Fields{"token": token}).Info("Token")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"code":    http.StatusUnauthorized,
				"status":  "error",
				"message": "Invalid or expired token",
			})
		}

		// Set the token claims as a request attribute
		c.Set("token", token.Claims)

		// Call the next middleware or handler
		return next(c)
	}
}

func main() {
	utility.PrintConsole("API started", "info")
	utility.PrintConsole("Loading application configuration", "info")
	configuration, errConfig := utility.LoadApplicationConfiguration("")
	if errConfig != nil {
		log.WithFields(log.Fields{"error": errConfig}).Fatal("Failed to load app configuration")
	}
	utility.PrintConsole("Application configuration loaded successfully", "info")

	db, gormDB, err := database.Open(configuration)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("Failed to open database")
	}
	defer db.Close()

	err = gormDB.AutoMigrate(&model.User{}, &model.Bansos{}, &model.Feedback{}, &model.Grocery{}, &model.QRCode{}, &model.BansosRegistration{})
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("Error to migrate database")
		return
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))


	// Repository
	authRepository := repository.NewAuthRepository(gormDB)
	userRepository := repository.NewUserRepository(gormDB)

	// Handler
	authHandler := handler.NewAuthHandler(authRepository)
	userHandler := handler.NewUserHandler(userRepository)

	// Router
	api := e.Group("/api")
	api.POST("/register", authHandler.Register)
	api.POST("/login", authHandler.Login)

	// Middleware Router
	apiAuth := api.Group("/", authMiddleware)
	// NOTE: This route is an example middleware route
	apiAuth.GET("users", userHandler.GetUserHandler)

	e.Logger.Fatal(e.Start(":" + configuration.Http.HttpPort))
}