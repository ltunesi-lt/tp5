package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
	db "tp5/DB/sqlc"
	"tp5/handlers"
)

func main() {
	// 1. Conectar a la base de datos
	database, err := sql.Open("sqlite", "./tu_base.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// 2. CREAR LA TABLA AUTOMÁTICAMENTE (Esto soluciona el error)
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		price REAL NOT NULL
	);
	`
	_, err = database.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Error al crear la tabla: ", err)
	}

	// 3. Inicializar sqlc y rutas
	queries := db.New(database)
	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux, queries)

	// 4. Iniciar el servidor
	log.Println("Servidor corriendo en http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}