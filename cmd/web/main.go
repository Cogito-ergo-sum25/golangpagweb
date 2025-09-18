package main

import (
	"github.com/Cogito-ergo-sum25/golangpagweb/pkg/config"
	"github.com/Cogito-ergo-sum25/golangpagweb/pkg/database"
	"github.com/Cogito-ergo-sum25/golangpagweb/pkg/handlers"
	"github.com/Cogito-ergo-sum25/golangpagweb/pkg/render"

	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main function
func main() {

	dbCfg := database.Config{
		Username: "root",
		Password: "12345",
		Host:     "127.0.0.1",
		Port:     "3306",
		DBName:   "mydb",
	}

	db, err := database.NewConnection(dbCfg)
	fmt.Println("¡Conexión exitosa a MySQL!")
	if err != nil {
		log.Fatal("Error conectando a DB:", err)
	}
	defer db.Close()

	// 2. Configuración general de la app
	app := config.AppConfig{ // Cambiado a valor (no puntero)
		DB:           db,
		InProduction: false,
	}

		// set up the session
		session = scs.New()
		session.Lifetime = 24 * time.Hour
		session.Cookie.Persist = true
		session.Cookie.SameSite = http.SameSiteLaxMode
		session.Cookie.Secure = app.InProduction
		app.Session = session

		
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Printf("Staring application on port %s\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}