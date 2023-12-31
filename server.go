package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"

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
	feedbackRepository := repository.NewFeedbackRepository(gormDB)
	bansosRegistrationRepository := repository.NewBansosRegistrationRepository(gormDB)
	bansosRepository := repository.NewBansosRepository(gormDB)
	qrCodeRepository := repository.NewQRCodeRepository(gormDB)

	// Handler
	authHandler := handler.NewAuthHandler(authRepository)
	userHandler := handler.NewUserHandler(userRepository, configuration.Storage.BucketName, configuration.Storage.Credentials)
	feedbackHandler := handler.NewFeedbackHandler(feedbackRepository)
	bansosRegistrationHandler := handler.NewBansosRegistrationHandler(bansosRegistrationRepository, configuration.Api.MlUrl)
	bansosHandler := handler.NewBansosHandler(bansosRepository, configuration.Storage.BucketName, configuration.Storage.Credentials)
	qrCodeHandler := handler.NewQRCodeHandler(qrCodeRepository, configuration.Http.Protocol, configuration.Http.Host, configuration.Http.HttpPort)

	// Router
	api := e.Group("/api")
	api.POST("/register", authHandler.Register)
	api.POST("/login", authHandler.Login)
	api.GET("/qr-codes/:uuid", qrCodeHandler.GetQRCodeByUUIDHandler)
	api.PUT("/bansos-registration/accept", bansosRegistrationHandler.AcceptBansosRegisHandler)
	api.PUT("/bansos-registration/reject", bansosRegistrationHandler.RejectBansosRegisHandler)

	// Middleware Router
	apiAuth := api.Group("/", authMiddleware)

	// User
	apiAuth.GET("users", userHandler.GetUserHandler)
	apiAuth.PUT("users", userHandler.UpdateUserHandler)

	// Bansos
	apiAuth.POST("bansos", bansosHandler.AddBansosHandler)
	apiAuth.GET("bansos", bansosHandler.GetBansosHandler)
	apiAuth.GET("bansos/:id", bansosHandler.GetBansosByIDHandler)
	apiAuth.PUT("bansos/:id", bansosHandler.UpdateBansosHandler)
	apiAuth.DELETE("bansos/:id", bansosHandler.DeleteBansosHandler)

	// feedback
	apiAuth.POST("feedbacks", feedbackHandler.AddFeedbackHandler)
	apiAuth.PUT("feedbacks", feedbackHandler.UpdateFeedbackHandler)
	apiAuth.DELETE("feedbacks", feedbackHandler.DeleteFeedbackHandler)
	apiAuth.GET("users/feedbacks", feedbackHandler.GetFeedbackByUserIDHandler)
	apiAuth.GET("bansos/feedbacks", feedbackHandler.GetFeedbackByBansosIDHandler)

	// bansos registration
	apiAuth.GET("users/bansos-registration", bansosRegistrationHandler.GetBansosRegisByUserIDHandler)
	apiAuth.GET("bansos/bansos-registration", bansosRegistrationHandler.GetBansosRegisByBansosIDHandler)
	apiAuth.POST("bansos-registration", bansosRegistrationHandler.RegisterBansosHandler)
	apiAuth.GET("bansos-registration", bansosRegistrationHandler.GetBansosRegisByIDHandler)
	apiAuth.GET("bansos-registration/on-progress", bansosRegistrationHandler.GetOnProgressBansosRegisHandler)
	apiAuth.PUT("bansos-registration/validate", bansosRegistrationHandler.ValidateBansosRegisHandler)

	// QR Code
	apiAuth.POST("qr-codes", qrCodeHandler.CreateQRCodeHandler)
	apiAuth.GET("qr-codes/show", qrCodeHandler.ShowQRCodeByIDHandler)

	e.Logger.Fatal(e.Start(":" + configuration.Http.HttpPort))
}
