// rmusr.go
package commands

import (
	"fmt"
	"strings"
)

type ParametrosRmusr struct {
	User string
}

// Analizar los parámetros del comando rmusr
func AnalizarParametrosRmusr(comando string) (ParametrosRmusr, error) {
	parametros := ParametrosRmusr{}
	for _, parte := range strings.Fields(comando) {
		if strings.HasPrefix(strings.ToLower(parte), "-user=") {
			parametros.User = strings.TrimPrefix(parte, "-user=")
		}
	}
	if parametros.User == "" {
		return parametros, fmt.Errorf("el parámetro -user es obligatorio")
	}
	return parametros, nil
}

func EjecutarRmusr(parametros ParametrosRmusr) string {
	if !VerificarSesionActiva() || UsuarioLogueado() != "root" {
		return "Error: solo el usuario root puede eliminar usuarios o no hay sesión activa"
	}

	rutaUsersTxt := obtenerRutaUsersTxt(ParticionActiva())
	contenido, err := leerUsersTxtDesdeDisco(rutaUsersTxt)
	if err != nil {
		return fmt.Sprintf("Error al leer users.txt: %v", err)
	}

	usuarioEncontrado := false
	lineas := strings.Split(contenido, "\n")

	for i, linea := range lineas {
		datos := strings.Split(linea, ",")
		if len(datos) == 5 && strings.TrimSpace(datos[1]) == "U" && strings.TrimSpace(datos[3]) == parametros.User && datos[0] != "0" {
			usuarioEncontrado = true
			datos[0] = "0" // Marcar usuario como eliminado
			lineas[i] = strings.Join(datos, ",")
			break
		}
	}

	if !usuarioEncontrado {
		return fmt.Sprintf("Error: el usuario %s no existe", parametros.User)
	}

	contenidoActualizado := strings.Join(lineas, "\n")
	if err := escribirUsersTxtEnDisco(rutaUsersTxt, contenidoActualizado); err != nil {
		return fmt.Sprintf("Error al escribir en users.txt: %v", err)
	}

	return fmt.Sprintf("Usuario %s eliminado exitosamente", parametros.User)
}
