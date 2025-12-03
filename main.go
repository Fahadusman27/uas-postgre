package main

import (
	. "GOLANG/Domain/config"
	"GOLANG/Domain/route"
	"log"
)

func main() {
	LoadEnv()

	// Connect PostgreSQL
	db := ConnectDB()
	if err := db.Ping(); err != nil {
		log.Fatal("Koneksi PostgreSQL gagal: ", err)
	}

	// Connect MongoDB
	ConnectMongoDB()

	// Jalankan migrations (opsional - bisa di-comment jika sudah dijalankan manual)
	// if err := RunMigrations(); err != nil {
	// 	log.Fatal("Migration gagal: ", err)
	// }

	app := route.NewApp(db)

	// Register routes
	route.AuthRoute(app)        // 5.1 Authentication
	route.UserRoute(app)         // 5.2 Users (Admin)
	route.AchievementRoute(app)  // 5.4 Achievements

	port := "4000"
	log.Printf("Server running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
