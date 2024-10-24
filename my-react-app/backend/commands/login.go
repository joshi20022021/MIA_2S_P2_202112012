package commands

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Definición de la estructura Montaje
type Montaje struct {
	ID   string
	Ruta string
}

type ParametrosLogin struct {
	User string
	Pass string
	ID   string
}

var (
	sesionIniciada   bool   = false
	usuarioLogueado  string = ""
	particionMontada string = ""
	montajes                = make(map[string]Montaje) // Mapa de montajes compartido
)

// Analiza los parámetros del comando login
func AnalizarParametrosLogin(comando string) (ParametrosLogin, error) {
	parametros := ParametrosLogin{}
	partes := strings.Split(comando, " ")

	if len(partes) == 0 || partes[0] != "login" {
		return parametros, fmt.Errorf("comando no reconocido")
	}

	for _, parte := range partes[1:] {
		if strings.HasPrefix(parte, "-user=") {
			parametros.User = strings.Trim(strings.TrimPrefix(parte, "-user="), "\"")
		} else if strings.HasPrefix(parte, "-pass=") {
			parametros.Pass = strings.Trim(strings.TrimPrefix(parte, "-pass="), "\"")
		} else if strings.HasPrefix(parte, "-id=") {
			parametros.ID = strings.Trim(strings.TrimPrefix(parte, "-id="), "\"")
		}
	}

	if parametros.User == "" || parametros.Pass == "" || parametros.ID == "" {
		return parametros, fmt.Errorf("los parámetros -user, -pass, e -id son obligatorios")
	}

	return parametros, nil
}

// Ejecutar el comando login
func EjecutarLogin(parametros ParametrosLogin) string {
	// Verificar si ya hay una sesión activa
	if sesionIniciada {
		return "Error: ya hay una sesión iniciada. Por favor, cierre sesión antes de iniciar una nueva."
	}

	// Obtener la ruta del archivo .mia montado
	rutaUsers := obtenerRutaUsersTxt(parametros.ID)
	if rutaUsers == "" {
		return fmt.Sprintf("Error: no se encontró la partición con el ID %s", parametros.ID)
	}

	// Verificar que el archivo users.txt exista
	if _, err := os.Stat(rutaUsers); os.IsNotExist(err) {
		return fmt.Sprintf("Error: el archivo users.txt no existe en la partición %s", parametros.ID)
	}

	// Leer el archivo users.txt desde el archivo .mia
	contenidoUsers, err := leerUsersTxtDesdeDisco(rutaUsers)
	if err != nil {
		return fmt.Sprintf("Error al leer el archivo users.txt: %v", err)
	}

	// Autenticación
	lineas := strings.Split(contenidoUsers, "\n")
	for _, linea := range lineas {
		if strings.HasPrefix(linea, "1,U,") {
			campos := strings.Split(linea, ",")
			if len(campos) >= 5 && campos[3] == parametros.User && campos[4] == parametros.Pass {
				// Iniciar sesión si usuario y contraseña coinciden
				sesionIniciada = true
				usuarioLogueado = parametros.User
				particionMontada = parametros.ID
				return fmt.Sprintf("Sesión iniciada con éxito como %s", parametros.User)
			}
		}
	}

	return "Error: autenticación fallida. Usuario o contraseña incorrectos."
}

// Cierra la sesión actual
func EjecutarLogout() string {
	if !sesionIniciada {
		return "Error: no hay ninguna sesión iniciada"
	}

	sesionIniciada = false
	usuarioLogueado = ""
	particionMontada = ""
	return "Sesión cerrada con éxito"
}

// Verifica si hay un usuario logueado
func VerificarSesionActiva() bool {
	return sesionIniciada && usuarioLogueado != "" && particionMontada != ""
}

// Retorna el usuario actualmente logueado
func UsuarioLogueado() string {
	return usuarioLogueado
}

// Obtiene la ruta del archivo users.txt dado el ID de la partición
func obtenerRutaUsersTxt(id string) string {
	montaje, encontrado := montajes[id]
	if encontrado {
		return montaje.Ruta
	}
	return ""
}

// Lee el archivo users.txt desde el disco montado
func leerUsersTxtDesdeDisco(ruta string) (string, error) {
	file, err := os.Open(ruta)
	if err != nil {
		return "", fmt.Errorf("error al abrir el disco %s: %v", ruta, err)
	}
	defer file.Close()

	// Leer el archivo users.txt desde el offset predefinido
	offset := int64(2048)        // Ajusta el valor de acuerdo con tu estructura
	buffer := make([]byte, 1024) // Tamaño del buffer para leer el archivo
	_, err = file.ReadAt(buffer, offset)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("error al leer el archivo users.txt: %v", err)
	}

	return string(buffer), nil
}
