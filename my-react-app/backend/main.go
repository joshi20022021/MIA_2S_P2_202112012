package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"my-react-app/backend/commands"
	"net/http"
	"strings"
)

type SolicitudComando struct {
	Comando string `json:"command"`
}

type RespuestaComando struct {
	Salida string `json:"output"`
	Error  string `json:"error,omitempty"` // Incluye campo de error opcional
}

// Maneja las solicitudes de comandos
func manejarComando(w http.ResponseWriter, r *http.Request) {
	// Configuración de CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	var solicitud SolicitudComando
	if err := json.NewDecoder(r.Body).Decode(&solicitud); err != nil {
		http.Error(w, "Solicitud incorrecta", http.StatusBadRequest)
		return
	}

	// Procesar el comando línea por línea
	var salidaFinal bytes.Buffer
	numeroDisco := 1 // Iniciar el contador de discos

	// Procesar el contenido del comando como un flujo continuo
	comandos := strings.Split(solicitud.Comando, "\n")
	for _, linea := range comandos {
		linea = strings.TrimSpace(linea)

		// Ignorar comentarios y líneas vacías
		if strings.HasPrefix(linea, "#") || linea == "" {
			continue
		}

		// Ejecutar los comandos según la línea
		salida := procesarLinea(linea, &numeroDisco)
		salidaFinal.WriteString(salida + "\n")
	}

	// Preparar y enviar la respuesta final
	respuesta := RespuestaComando{Salida: salidaFinal.String()}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		fmt.Println("Error al codificar la respuesta:", err)
	}
}

// procesarLinea ejecuta el comando según el contenido de la línea y devuelve la salida
func procesarLinea(linea string, numeroDisco *int) string {
	// Convertir la línea a minúsculas para hacerla insensible a mayúsculas/minúsculas
	comando := strings.ToLower(strings.Fields(linea)[0]) // Solo toma el primer término de la línea
	var salida string

	switch comando {
	case "mkdisk":
		parametros, err := commands.AnalizarParametrosMkDisk(linea)
		if err != nil {
			salida = "Error: " + err.Error()
		} else {
			salida = commands.EjecutarMkDisk(parametros)
		}

	case "rmdisk":
		parametros, err := commands.AnalizarParametrosRmDisk(linea)
		if err != nil {
			salida = "Error: " + err.Error()
		} else {
			salida = commands.EjecutarRmDisk(parametros)
		}

	case "fdisk":
		parametros, err := commands.AnalizarParametrosFDisk(linea)
		if err != nil {
			salida = "Error: " + err.Error()
		} else {
			salida = commands.EjecutarFDisk(parametros)
		}

	case "mount":
		parametros, err := commands.AnalizarParametrosMontarNew(linea)
		if err != nil {
			salida = "Error: " + err.Error()
		} else {
			salida = commands.EjecutarMontar(parametros)
			*numeroDisco++ // Incrementar el número de disco después de cada montaje
		}

	case "mkfs":
		parametros, err := commands.AnalizarParametrosMkfs(linea)
		if err != nil {
			salida = "Error: " + err.Error()
		} else {
			salida = commands.EjecutarMkfs(parametros)
		}

	case "login":
		parametros, err := commands.AnalizarParametrosLogin(linea)
		if err != nil {
			salida = "Error: " + err.Error()
		} else {
			salida = commands.EjecutarLogin(parametros)
		}

	case "logout":
		salida = commands.EjecutarLogout()

	case "cat":
		if commands.VerificarSesionActiva() {
			parametros, err := commands.AnalizarParametrosCat(linea)
			if err != nil {
				salida = "Error: " + err.Error()
			} else {
				salida = commands.EjecutarCat(parametros, commands.UsuarioLogueado())
			}
		} else {
			salida = "Error: Debes iniciar sesión para ejecutar este comando."
		}

	default:
		salida = "Comando no reconocido"
	}

	return salida
}

func main() {
	http.HandleFunc("/api/command", manejarComando)

	fmt.Println("Servidor corriendo en http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}
