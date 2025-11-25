package main

import (
	. "POSTGRE/Domain/config"
	"POSTGRE/Domain/route"
	"log"
)

func main() {
	LoadEnv()
    db := ConnectDB()
    if err := db.Ping(); err != nil {
        log.Fatal("Koneksi database gagal: ", err)
    }


	// userRepo := repository.NewUserRepository(db)
    app := route.NewApp(db)
	// routes.AuthRoutes(app, userRepo)
	// routes.Alumni(app, &userRepo)
	// routes.PekerjaanAlumni(app, &userRepo)
	// routes.UserRoutes(app)

	port := "4000"

    log.Fatal(app.Listen(":" + port))
}