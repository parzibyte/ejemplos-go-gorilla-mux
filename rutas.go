package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Usuario struct {
	Id     int    `json:"id"`
	Correo string `json:"correo"`
}

// Nota: estos usuarios podrían venir de una base de datos que podrían obtenerse dentro
// de cada ruta

var usuarios []Usuario = []Usuario{
	Usuario{
		Id:     1,
		Correo: "contacto@parzibyte.me",
	},
	Usuario{
		Id:     2,
		Correo: "john.galt@atlas.com",
	},
}

func definirRutas(enrutador *mux.Router) {
	// Rutas para ejemplificar variables
	enrutador.HandleFunc("/ventas/{tipo}/{anio}", func(respuesta http.ResponseWriter, peticion *http.Request) {
		variablesDePeticion := mux.Vars(peticion)
		tipo := variablesDePeticion["tipo"]
		anio := variablesDePeticion["anio"]
		respuesta.Write([]byte("El tipo de venta es " + tipo + " y el año es " + anio))
	}).Methods("GET")

	enrutador.HandleFunc("/pedidos/{anio:[0-9]{4}}", func(respuesta http.ResponseWriter, peticion *http.Request) {
		variablesDePeticion := mux.Vars(peticion)
		anio := variablesDePeticion["anio"]
		respuesta.Write([]byte("Pedidos del año: " + anio))
	}).Methods("GET")

	enrutador.HandleFunc("/libros/{orden:(?:ascendente|descendente)}", func(respuesta http.ResponseWriter, peticion *http.Request) {
		variablesDePeticion := mux.Vars(peticion)
		orden := variablesDePeticion["orden"]
		respuesta.Write([]byte("Libros ordenados: " + orden))
	}).Methods("GET")

	enrutador.HandleFunc("/cursos", func(respuesta http.ResponseWriter, peticion *http.Request) {
		variablesGet := peticion.URL.Query()
		// Cada variable es un arreglo
		orden := variablesGet["orden"]
		// Si mide más que 0 entonces el orden sí está definido
		if len(orden) > 0 {
			fmt.Fprintf(respuesta, "El orden: %s.", orden[0]) // Acceder al primer elemento
		}
		fmt.Fprintf(respuesta, "Parámetros de consulta: %v", variablesGet)
	}).Methods("GET")

	enrutador.HandleFunc("/usuarios", obtenerUsuarios).Methods("GET")
	enrutador.HandleFunc("/usuario/{id}", obtenerUsuarioPorId).Methods("GET")
	enrutador.HandleFunc("/usuario", agregarUsuario).Methods("POST")
	enrutador.HandleFunc("/usuario", actualizarUsuario).Methods("PUT")
	enrutador.HandleFunc("/usuario/{id}", eliminarUsuario).Methods("DELETE")
}

func middlewareLog(siguienteManejador http.Handler) http.Handler {
	return http.HandlerFunc(
		func(respuesta http.ResponseWriter, peticion *http.Request) {
			log.Printf("Nueva petición. Método: %s. IP: %s. URL solicitada: %s\n",
				peticion.Method, peticion.RemoteAddr, peticion.URL)
			siguienteManejador.ServeHTTP(respuesta, peticion)
		})
}
func middlewareEliminar(siguienteManejador http.Handler) http.Handler {
	return http.HandlerFunc(
		func(respuesta http.ResponseWriter, peticion *http.Request) {
			// Si no llamamos a siguienteManejador, se detiene
			// así que podemos aquí comprobar algo y detener determinada acción
			// Por ejemplo, permitir solo si no son de tipo DELETE
			if peticion.Method == http.MethodDelete {
				http.Error(respuesta, "Permiso denegado", http.StatusForbidden)
			} else {
				// En caso de que sea permitida llamamos a siguienteManejador
				// y le pasamos la respuesta con la petición
				siguienteManejador.ServeHTTP(respuesta, peticion)
			}
		})
}

func agregarUsuario(respuesta http.ResponseWriter, peticion *http.Request) {
	var usuarioNuevo Usuario
	// Intenta decodificar el cuerpo de la petición (peticion.Body) dentro de usuario (&usuario)
	err := json.NewDecoder(peticion.Body).Decode(&usuarioNuevo)
	if err != nil {
		json.NewEncoder(respuesta).Encode("Cuerpo de petición no válido")
		return
	}
	// Si el usuario era válido lo agregamos al arreglo
	usuarios = append(usuarios, usuarioNuevo)
	json.NewEncoder(respuesta).Encode(usuarioNuevo)
}

func eliminarUsuario(respuesta http.ResponseWriter, peticion *http.Request) {
	// Nota: en realidad no se elimina, porque la implementación para eliminar
	// de un arreglo es un poco difícil de explicar y no es necesario para los propósitos
	// del código. Se hace lo mismo que en la función obtener, pero quería demostrar el método DELETE
	variablesDePeticion := mux.Vars(peticion)
	// El id viene como cadena, hay que convertirlo a entero de 32 bits
	// Aquí "id" es la variable que indicamos en la ruta
	idUsuarioBuscado, err := strconv.Atoi(variablesDePeticion["id"])
	// Si no es un entero válido:
	if err != nil {
		json.NewEncoder(respuesta).Encode("Error: id inválido")
		return
	}

	// Nota: el id puedes usarlo para filtrar en una base de datos o algo así,
	// aquí simplemente lo buscamos dentro del arreglo

	// Buscamos...
	for _, usuario := range usuarios {
		// Si lo encontramos lo devolvemos y terminamos la función
		if usuario.Id == idUsuarioBuscado {
			json.NewEncoder(respuesta).Encode(usuario)
			return
		}
	}
	// Si no lo encontramos, indicamos un error
	json.NewEncoder(respuesta).Encode("No existe un usuario con el id proporcionado")
}

func actualizarUsuario(respuesta http.ResponseWriter, peticion *http.Request) {
	var usuarioActualizado Usuario
	// Intenta decodificar el cuerpo de la petición (peticion.Body) dentro de usuario (&usuario)
	err := json.NewDecoder(peticion.Body).Decode(&usuarioActualizado)
	if err != nil {
		json.NewEncoder(respuesta).Encode("Cuerpo de petición no válido")
		return
	}
	// Ya tenemos al usuario, ahora lo actualizamos en el arreglo.
	// Para ello debemos obtener el índice
	for indice, usuarioExistente := range usuarios {
		// Si lo encontramos lo devolvemos y terminamos la función
		if usuarioExistente.Id == usuarioActualizado.Id {
			usuarios[indice] = usuarioActualizado
			json.NewEncoder(respuesta).Encode(usuarioActualizado)
			return
		}
	}
	// Si no lo encontramos por Id entonces ese usuario no existía
	json.NewEncoder(respuesta).Encode("Usuario no encontrado")
}

func obtenerUsuarios(respuesta http.ResponseWriter, peticion *http.Request) {
	// También podrías codificar otro tipo de datos como un arreglo plano
	// o una simple variable, todo lo soportado por JSON:
	// https://parzibyte.me/blog/2019/05/16/codificar-decodificar-json-go-golang/
	json.NewEncoder(respuesta).Encode(usuarios)
}

func obtenerUsuarioPorId(respuesta http.ResponseWriter, peticion *http.Request) {
	variablesDePeticion := mux.Vars(peticion)
	// El id viene como cadena, hay que convertirlo a entero de 32 bits
	// Aquí "id" es la variable que indicamos en la ruta
	idUsuarioBuscado, err := strconv.Atoi(variablesDePeticion["id"])
	// Si no es un entero válido:
	if err != nil {
		json.NewEncoder(respuesta).Encode("Error: id inválido")
		return
	}

	// Nota: el id puedes usarlo para filtrar en una base de datos o algo así,
	// aquí simplemente lo buscamos dentro del arreglo

	// Buscamos...
	for _, usuario := range usuarios {
		// Si lo encontramos lo devolvemos y terminamos la función
		if usuario.Id == idUsuarioBuscado {
			json.NewEncoder(respuesta).Encode(usuario)
			return
		}
	}
	// Si no lo encontramos, indicamos un error
	json.NewEncoder(respuesta).Encode("No existe un usuario con el id proporcionado")
}
