package main

import (
	"learn/config"
	"learn/handler"
	"learn/repository"
	"learn/service"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := config.ConnectDb()
	validate := validator.New()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, validate)

	addresRepo := repository.NewAddressRepository(db)
	addressService := service.NewAddressService(&addresRepo)
	addressHandler := handler.NewAddressHandler(addressService, validate)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)

	// USER
	// Public
	router.Post("/register", userHandler.Register)
	router.Post("/login", userHandler.Login)

	// Auth
	router.Get("/profile", handler.Auth(userHandler.Profile))
	router.Post("/change-password", handler.Auth(userHandler.ChangePassword))

	// ADDRESS
	// Public

	// Auth
	router.Post("/address", handler.Auth(addressHandler.AddAddress))
	router.Get("/address", handler.Auth(addressHandler.AddAddress))

	http.ListenAndServe(":3000", router)
}
