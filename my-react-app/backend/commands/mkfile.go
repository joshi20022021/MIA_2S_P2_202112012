// mkfile.go
package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type ParametrosMkfile struct {
	Path string
	R    bool
	Size int
	Cont string
}

// Analizar los parámetros para el comando mkfile
func AnalizarParametrosMkfile(comando string) (ParametrosMkfile, error) {
	parametros := ParametrosMkfile{Size: 0}

	for _, parte := range strings.Fields(comando) {
		if strings.HasPrefix(parte, "-path=") {
			parametros.Path = strings.TrimPrefix(parte, "-path=")
			parametros.Path = strings.Trim(parametros.Path, "\"")
		} else if parte == "-r" {
			parametros.R = true
		} else if strings.HasPrefix(parte, "-size=") {
			sizeStr := strings.TrimPrefix(parte, "-size=")
			fmt.Sscanf(sizeStr, "%d", &parametros.Size)
			if parametros.Size < 0 {
				return parametros, fmt.Errorf("error: el tamaño del archivo no puede ser negativo")
			}
		} else if strings.HasPrefix(parte, "-cont=") {
			parametros.Cont = strings.TrimPrefix(parte, "-cont=")
			parametros.Cont = strings.Trim(parametros.Cont, "\"")
		} else if strings.HasPrefix(parte, "-r=") || strings.HasPrefix(parte, "-size=") || strings.HasPrefix(parte, "-cont=") {
			return parametros, fmt.Errorf("error: el parámetro %s no debe recibir un valor", parte)
		}
	}

	if parametros.Path == "" {
		return parametros, fmt.Errorf("error: el parámetro -path es obligatorio")
	}

	return parametros, nil
}

func EjecutarMkfile(parametros ParametrosMkfile) string {
	if !VerificarSesionActiva() {
		return "Error: no hay sesión activa"
	}

	rutaCompleta := obtenerRutaAbsoluta(parametros.Path)

	// Si el archivo ya existe, preguntar si se desea sobrescribir
	if _, err := os.Stat(rutaCompleta); err == nil {
		return fmt.Sprintf("El archivo %s ya existe. ¿Desea sobrescribir?", parametros.Path)
	}

	// Si `-r` está presente, crear directorios faltantes en el path
	dir := filepath.Dir(rutaCompleta)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if parametros.R {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Sprintf("Error al crear directorios en el path %s: %v", dir, err)
			}
		} else {
			return fmt.Sprintf("Error: las carpetas en el path %s no existen y no se especificó -r para crearlas", dir)
		}
	}

	// Determinar contenido del archivo
	var contenido []byte
	if parametros.Cont != "" {
		// Leer contenido desde archivo de disco
		content, err := ioutil.ReadFile(parametros.Cont)
		if err != nil {
			return fmt.Sprintf("Error al leer el archivo de contenido: %v", err)
		}
		contenido = content
	} else {
		// Generar contenido según el tamaño especificado
		contenido = generarContenido(parametros.Size)
	}

	// Crear archivo y escribir contenido
	if err := ioutil.WriteFile(rutaCompleta, contenido, 0664); err != nil {
		return fmt.Sprintf("Error al crear el archivo: %v", err)
	}

	return fmt.Sprintf("Archivo %s creado exitosamente con tamaño %d bytes", parametros.Path, len(contenido))
}

// Función auxiliar para generar contenido del archivo en función del tamaño especificado
func generarContenido(size int) []byte {
	contenido := make([]byte, size)
	for i := 0; i < size; i++ {
		contenido[i] = byte('0' + (i % 10)) // Rellenar con números del 0 al 9
	}
	return contenido
}

// Función auxiliar para obtener ruta absoluta del archivo a crear
func obtenerRutaAbsoluta(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	// Asumimos que la ruta es relativa al directorio actual
	absPath, _ := filepath.Abs(path)
	return absPath
}
