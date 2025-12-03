package main

import (
	. "GOLANG/Domain/config"
	"GOLANG/Domain/route"
	"log"
)

func main() {
	LoadEnv()
	db := ConnectDB()
	if err := db.Ping(); err != nil {
		log.Fatal("Koneksi database gagal: ", err)
	}

	// Jalankan migrations (opsional - bisa di-comment jika sudah dijalankan manual)
	// if err := RunMigrations(); err != nil {
	// 	log.Fatal("Migration gagal: ", err)
	// }

	app := route.NewApp(db)
	route.AuthRoute(app)

	port := "4000"
	log.Printf("Server running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
