package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	// Crear el enrutador y definir las rutas en la función definirRutas
	enrutador := mux.NewRouter()
	definirRutas(enrutador)
	enrutador.Use(middlewareLog)
	enrutador.Use(middlewareEliminar)
	// Dirección del servidor. En este caso solo indicamos el puerto
	// pero podría ser algo como "127.0.0.1:8000"
	direccion := ":8000"

	servidor := &http.Server{
		Handler: enrutador,
		Addr:    direccion,
		// Timeouts para evitar que el servidor se quede "colgado" por siempre
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Printf("Escuchando en %s. Presiona CTRL + C para salir", direccion)
	log.Fatal(servidor.ListenAndServe())
}
