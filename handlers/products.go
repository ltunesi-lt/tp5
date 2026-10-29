package handlers

import (
	"net/http"
	"strconv"
	db "tp5/DB/sqlc"
	"tp5/views"
)

func RegisterRoutes(mux *http.ServeMux, queries *db.Queries) {
	// 1. Ruta GET / (Listar productos)
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		products, err := queries.ListProducts(ctx)
		if err != nil {
			println("ERROR REAL DE SQLC:", err.Error())
			http.Error(w, "Error al obtener los productos de la base de datos", http.StatusInternalServerError)
			return
		}

		err = views.IndexPage("Listado de Productos", products).Render(ctx, w)
		if err != nil {
			http.Error(w, "Error al renderizar la vista", http.StatusInternalServerError)
			return
		}
	})

	// 2. NUEVA RUTA: POST /products (Crear un producto)
	mux.HandleFunc("POST /products", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Parsear los datos del formulario
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
			return
		}

		// Extraer los valores enviados por los inputs ("name" y "price")
		name := r.FormValue("name")
		priceStr := r.FormValue("price")

		// Convertir el precio de string a float64
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			http.Error(w, "El precio proporcionado no es válido", http.StatusBadRequest)
			return
		}

		// Usar CreateProduct de sqlc para guardar el nuevo producto
		// (Asegúrate de que los campos coincidan con los generados por tu schema)
		// Usar CreateProduct de sqlc (solo retorna error)
		err = queries.CreateProduct(ctx, db.CreateProductParams{
			Name:  name,
			Price: price,
		})
		if err != nil {
			http.Error(w, "Error al guardar el producto en la base de datos", http.StatusInternalServerError)
			return
		}

		// Redirigir al usuario de vuelta a la página principal para ver la lista actualizada
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}