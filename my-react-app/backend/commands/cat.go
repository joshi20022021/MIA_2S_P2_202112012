package commands

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type ParametrosCat struct {
	Archivos []string
}

// Analiza los parámetros del comando cat
func AnalizarParametrosCat(comando string) (ParametrosCat, error) {
	parametros := ParametrosCat{}
	partes := strings.Fields(comando)

	for _, parte := range partes {
		if strings.HasPrefix(strings.ToLower(parte), "-file") {
			archivo := strings.Trim(strings.TrimPrefix(parte, "-file"), "=")
			archivo = strings.Trim(archivo, "\"")
			parametros.Archivos = append(parametros.Archivos, archivo)
		}
	}

	if len(parametros.Archivos) == 0 {
		return parametros, fmt.Errorf("se debe especificar al menos un archivo con -file")
	}

	return parametros, nil
}

// Verifica si el usuario actual tiene permisos de lectura sobre el archivo (simulado)
func VerificarPermisosLectura(archivo string) bool {
	// Simulación: siempre devolver true por simplicidad.
	// En una implementación real, debes verificar permisos basados en el sistema de archivos y el usuario logueado.
	return true
}

// Ejecuta el comando cat para mostrar el contenido de uno o varios archivos
func EjecutarCat(parametros ParametrosCat) (string, error) {
	var contenidoFinal strings.Builder

	for _, archivo := range parametros.Archivos {
		if !VerificarPermisosLectura(archivo) {
			return "", fmt.Errorf("el usuario no tiene permiso de lectura sobre el archivo: %s", archivo)
		}

		contenido, err := ioutil.ReadFile(archivo)
		if err != nil {
			return "", fmt.Errorf("error al leer el archivo %s: %v", archivo, err)
		}

		// Concatenar contenido al resultado final
		contenidoFinal.WriteString(fmt.Sprintf("Contenido del archivo %s:\n", archivo))
		contenidoFinal.WriteString(string(contenido))
		contenidoFinal.WriteString("\n\n") // Añadir una línea en blanco entre archivos
	}

	return contenidoFinal.String(), nil
}
