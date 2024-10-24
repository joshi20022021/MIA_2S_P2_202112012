package commands

import (
	"fmt"
	"os"
	"strings"
)

type ParametrosCat struct {
	Archivos []string
}

// Analizar el comando CAT para obtener los archivos
func AnalizarParametrosCat(comando string) (ParametrosCat, error) {
	parametros := ParametrosCat{}
	partes := strings.Split(comando, " ")

	if len(partes) == 0 || partes[0] != "cat" {
		return parametros, fmt.Errorf("comando no reconocido")
	}

	// Buscar los archivos pasados como parámetros
	for _, parte := range partes[1:] {
		if strings.HasPrefix(parte, "-file") {
			file := strings.Trim(strings.Split(parte, "=")[1], "\"")
			parametros.Archivos = append(parametros.Archivos, file)
		}
	}

	if len(parametros.Archivos) == 0 {
		return parametros, fmt.Errorf("al menos un archivo es obligatorio")
	}

	return parametros, nil
}

// Ejecutar el comando CAT para concatenar el contenido de los archivos
func EjecutarCat(parametros ParametrosCat, usuario string) string {
	if !VerificarSesionActiva() {
		return "Error: No hay una sesión activa."
	}

	rutaParticion := obtenerRutaUsersTxt(particionMontada)
	if rutaParticion == "" {
		return fmt.Sprintf("Error: No se encontró la partición montada para el ID %s", particionMontada)
	}

	var resultado strings.Builder
	for _, archivo := range parametros.Archivos {
		contenido, err := leerArchivoEnParticion(rutaParticion, archivo)
		if err != nil {
			return fmt.Sprintf("Error al leer el archivo %s: %v", archivo, err)
		}
		// Concatenar contenido de archivos
		resultado.WriteString(contenido + "\n")
	}

	return resultado.String()
}

// Leer el archivo desde la partición montada
func leerArchivoEnParticion(rutaDisco, nombreArchivo string) (string, error) {
	file, err := os.Open(rutaDisco)
	if err != nil {
		return "", fmt.Errorf("error al abrir el disco %s: %v", rutaDisco, err)
	}
	defer file.Close()

	// Ajustar el offset si es necesario según la estructura de la partición
	offset := int64(2048)
	buffer := make([]byte, 1024)
	_, err = file.ReadAt(buffer, offset)
	if err != nil {
		return "", fmt.Errorf("error al leer el archivo %s: %v", nombreArchivo, err)
	}

	return string(buffer), nil
}
