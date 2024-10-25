package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type ParametrosCat struct {
	Archivos []string
}

// Analizar los parámetros del comando cat
func AnalizarParametrosCat(comando string) (ParametrosCat, error) {
	parametros := ParametrosCat{}

	for _, parte := range strings.Fields(comando) {
		if strings.HasPrefix(strings.ToLower(parte), "-file") {
			parametro := strings.SplitN(parte, "=", 2)
			if len(parametro) == 2 {
				parametros.Archivos = append(parametros.Archivos, strings.Trim(parametro[1], "\""))
			}
		}
	}

	if len(parametros.Archivos) == 0 {
		return parametros, fmt.Errorf("no se especificaron archivos para el comando cat")
	}

	return parametros, nil
}

// Ejecuta el comando cat para leer múltiples archivos
func EjecutarCat(parametros ParametrosCat) (string, error) {
	var contenido strings.Builder

	for _, ruta := range parametros.Archivos {
		// Comprobar si el archivo existe antes de intentar abrirlo
		if _, err := os.Stat(ruta); os.IsNotExist(err) {
			return "", fmt.Errorf("error: el archivo %s no existe o la ruta es incorrecta", ruta)
		}

		// Intentar abrir el archivo
		file, err := os.Open(ruta)
		if err != nil {
			return "", fmt.Errorf("no se pudo abrir el archivo %s: %v", ruta, err)
		}
		defer file.Close()

		// Leer el contenido del archivo línea por línea
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			contenido.WriteString(scanner.Text() + "\n")
		}

		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("error al leer el archivo %s: %v", ruta, err)
		}
		contenido.WriteString("\n") // Separador entre archivos
	}

	return contenido.String(), nil
}
