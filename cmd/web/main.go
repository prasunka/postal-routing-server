package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/prasunka/postal-routing-server/pkg/models/mysql"
	"github.com/rs/cors"
)

func dsn(envs map[string]string) string {
	username := envs["DB_USERNAME"]
	password := envs["DB_PASSWORD"]
	host := envs["DB_HOST"]
	dbName := envs["DB_NAME"]
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, host, dbName)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

type application struct {
	errorLog  *log.Logger
	infoLog   *log.Logger
	routes    *mysql.RouteModel
	endpoints *mysql.EndpointModel
}

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	envs, err := godotenv.Read(".env")
	if err != nil {
		errorLog.Fatal("Error loading .env file")
	}

	dsn := dsn(envs)
	DB, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		infoLog:   infoLog,
		errorLog:  errorLog,
		routes:    &mysql.RouteModel{DB: DB},
		endpoints: &mysql.EndpointModel{DB: DB},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/create", app.createRoute)

	handler := cors.Default().Handler(mux)

	log.Println("Starting server on :4000")
	err = http.ListenAndServe(":4000", handler)

	errorLog.Fatal(err)
}
