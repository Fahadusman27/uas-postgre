package main

import (
	. "GOLANG/Domain/config"
	"GOLANG/Domain/route"
	"log"

	_ "GOLANG/docs" // Import generated swagger docs

	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Achievement Management System API
// @version 1.0
// @description API untuk sistem manajemen prestasi mahasiswa
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:4000
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	LoadEnv()

	// Connect PostgreSQL
	db := ConnectDB()
	if err := db.Ping(); err != nil {
		log.Fatal("Koneksi PostgreSQL gagal: ", err)
	}

	// Connect MongoDB
	ConnectMongoDB()

	app := route.NewApp(db)

	// Swagger documentation
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Register routes
	route.AuthRoute(app)
	route.UserRoute(app)
	route.AchievementRoute(app)

	port := "4000"
	log.Printf("Server running on port %s", port)
	log.Printf("Swagger documentation available at http://localhost:%s/swagger/index.html", port)
	log.Fatal(app.Listen(":" + port))
}
