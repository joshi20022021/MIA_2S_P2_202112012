// mkdir.go
package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ParametrosMkdir struct {
	Path string
	P    bool
}

// Analizar los parámetros para el comando mkdir
func AnalizarParametrosMkdir(comando string) (ParametrosMkdir, error) {
	parametros := ParametrosMkdir{}

	for _, parte := range strings.Fields(comando) {
		if strings.HasPrefix(parte, "-path=") {
			parametros.Path = strings.TrimPrefix(parte, "-path=")
			parametros.Path = strings.Trim(parametros.Path, "\"") // Quitar comillas si existen
		} else if parte == "-p" {
			parametros.P = true
		} else if strings.HasPrefix(parte, "-p=") {
			return parametros, fmt.Errorf("error: el parámetro -p no debe recibir un valor")
		}
	}

	if parametros.Path == "" {
		return parametros, fmt.Errorf("error: el parámetro -path es obligatorio")
	}

	return parametros, nil
}

func EjecutarMkdir(parametros ParametrosMkdir) string {
	if !VerificarSesionActiva() {
		return "Error: no hay sesión activa"
	}

	rutaCompleta := obtenerRutaAbsoluta(parametros.Path)
	parentDir := filepath.Dir(rutaCompleta)

	// Verificar si el usuario tiene permisos de escritura en la carpeta padre
	if _, err := os.Stat(parentDir); os.IsNotExist(err) && !parametros.P {
		return fmt.Sprintf("Error: la carpeta padre %s no existe y no se especificó -p para crearla", parentDir)
	} else if err == nil && !tienePermisoEscritura(parentDir) {
		return fmt.Sprintf("Error: no tienes permiso de escritura en la carpeta padre %s", parentDir)
	}

	// Crear las carpetas según la opción -p
	if parametros.P {
		if err := os.MkdirAll(rutaCompleta, 0755); err != nil {
			return fmt.Sprintf("Error al crear las carpetas en el path %s: %v", rutaCompleta, err)
		}
	} else {
		// Crear solo la carpeta especificada
		if err := os.Mkdir(rutaCompleta, 0755); err != nil {
			return fmt.Sprintf("Error al crear la carpeta en el path %s: %v", rutaCompleta, err)
		}
	}

	return fmt.Sprintf("Carpeta %s creada exitosamente", parametros.Path)
}

// Función auxiliar para verificar si el usuario actual tiene permisos de escritura en una carpeta
func tienePermisoEscritura(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	perm := fileInfo.Mode().Perm()
	return (perm&0200 != 0) || UsuarioLogueado() == "root" // Permitir si root está logueado
}

// Función auxiliar para obtener ruta absoluta de la carpeta a crear
func obtenerRutaAbsoluta(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	// Asumimos que la ruta es relativa al directorio actual
	absPath, _ := filepath.Abs(path)
	return absPath
}
