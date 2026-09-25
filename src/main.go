package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/dotenv-org/godotenvvault"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/visionJAMx/whoknows/src/internal/delivery"
	"github.com/visionJAMx/whoknows/src/internal/repository"
)

// @title WhoKnows
// @version 0.1.0
// @description API til WhoKnows-projektet.

func main() {
	// Indlæs projektets miljøvariabler.
	if err := godotenvvault.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	databasePath := os.Getenv("DATABASE_PATH")
	if databasePath == "" {
		databasePath = "../data/whoknows.db"
	}

	// Åbn databaseforbindelsen, og luk den, når main afsluttes normalt.
	db, err := repository.Open(databasePath)
	if err != nil {
		log.Fatalf("could not connect to repository: %v", err)
	}
	defer db.Close()

	// Opret tabellerne, hvis de ikke allerede findes.
	if err := repository.Initialize(context.Background(), db); err != nil {
		log.Fatalf("could not initialize repository: %v", err)
	}

	// Sessionsnøglen skal komme fra miljøet — aldrig hardcodes.
	sessionSecret := os.Getenv("SESSION_SECRET")
	if len(sessionSecret) < 32 {
		log.Fatal("SESSION_SECRET must contain at least 32 characters")
	}

	// Sessionen gemmes i en signeret cookie.
	store := cookie.NewStore([]byte(sessionSecret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   gin.Mode() == gin.ReleaseMode,
		SameSite: http.SameSiteLaxMode,
	})

	router := gin.Default()
	router.HTMLRender = createRenderer()
	router.Static("/static", "./static")

	// Tilslut sessions før de routes, der bruger dem.
	router.Use(sessions.Sessions("whoknows_session", store))

	//Routes
	router.GET("/login", delivery.LoginPage)
	router.GET("/register", delivery.RegisterPage)
	router.GET("/about", delivery.AboutPage)
	router.GET("/", delivery.SearchPage(db))
	// Logout rydder den aktuelle brugers session.
	router.GET("/api/logout", delivery.LogoutHandler)

	loginHandler := delivery.NewLoginHandler(db)
	router.POST("/api/login", loginHandler.APILogin)

	registerHandler := delivery.NewRegisterHandler(db)
	router.POST("/api/register", registerHandler.APIRegister)

	router.GET("/api/search", delivery.SearchAPI(db))

	router.GET("/health", delivery.HealthHandler)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}

// Hver side bruger det fælles layout sammen med sin egen template.
func createRenderer() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("login", "templates/layout.html", "templates/login.html")
	r.AddFromFiles("search", "templates/layout.html", "templates/search.html")
	r.AddFromFiles("register", "templates/layout.html", "templates/register.html")
	r.AddFromFiles("about", "templates/layout.html", "templates/about.html")
	return r
}
