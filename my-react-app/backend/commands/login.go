package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type ParametrosLogin struct {
	User string
	Pass string
	Id   string
}

// Analiza los parámetros para el comando login
func AnalizarParametrosLogin(comando string) (ParametrosLogin, error) {
	parametros := ParametrosLogin{}

	for _, parte := range strings.Fields(comando) {
		if strings.HasPrefix(parte, "-user=") {
			parametros.User = strings.Trim(strings.TrimPrefix(parte, "-user="), "\"")
		} else if strings.HasPrefix(parte, "-pass=") {
			parametros.Pass = strings.Trim(strings.TrimPrefix(parte, "-pass="), "\"")
		} else if strings.HasPrefix(parte, "-id=") {
			parametros.Id = strings.TrimPrefix(parte, "-id=")
		}
	}

	if parametros.User == "" || parametros.Pass == "" || parametros.Id == "" {
		return parametros, fmt.Errorf("los parámetros user, pass e id son obligatorios")
	}

	return parametros, nil
}
func EjecutarLogin(parametros ParametrosLogin) string {
	if sesionActiva {
		return "Error: Ya hay un usuario logueado. Cierra sesión antes de iniciar una nueva."
	}

	particion, existe := particionesMontadas[parametros.Id]
	if !existe {
		return fmt.Sprintf("Error: No se encontró el id %s en las particiones montadas", parametros.Id)
	}

	rutaUsersTxt := obtenerRutaUsersTxt(particion)

	file, err := os.Open(rutaUsersTxt)
	if err != nil {
		return fmt.Sprintf("Error al abrir users.txt: %v", err)
	}
	defer file.Close()

	// Leer y verificar usuario
	scanner := bufio.NewScanner(file)
	autenticado := false
	for scanner.Scan() {
		linea := scanner.Text()
		datos := strings.Split(linea, ",")
		if len(datos) == 5 && strings.TrimSpace(datos[1]) == "U" {
			usuario := strings.TrimSpace(datos[3])
			contrasena := strings.TrimSpace(datos[4])
			if usuario == parametros.User && contrasena == parametros.Pass {
				autenticado = true
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Sprintf("Error al leer users.txt: %v", err)
	}

	if autenticado {
		IniciarSesion(parametros.User, particion)
		return fmt.Sprintf("Sesión iniciada correctamente como %s", parametros.User)
	}

	return "Error: Usuario o contraseña incorrectos"
}
