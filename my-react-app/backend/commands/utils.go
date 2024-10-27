// utils.go

package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// obtenerRutaUsersTxt devuelve la ruta completa de users.txt en el sistema de archivos simulado
func obtenerRutaUsersTxt(particion ParticionMontada) string {
	// Asegúrate de construir la ruta de manera compatible
	return filepath.Join(filepath.Dir(particion.Ruta), "users.txt")
}

// leerUsersTxtDesdeDisco lee el contenido del archivo users.txt de una partición específica
func leerUsersTxtDesdeDisco(ruta string) (string, error) {
	file, err := os.Open(ruta)
	if err != nil {
		return "", fmt.Errorf("no se pudo abrir el archivo %s: %v", ruta, err)
	}
	defer file.Close()

	content := new(strings.Builder)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error al leer %s: %v", ruta, err)
	}
	return content.String(), nil
}

// escribirUsersTxtEnDisco escribe el contenido en users.txt dentro de la partición montada
func escribirUsersTxtEnDisco(ruta string, contenido string) error {
	// Crea el directorio si no existe
	if err := os.MkdirAll(filepath.Dir(ruta), 0755); err != nil {
		return fmt.Errorf("no se pudo crear el directorio para users.txt: %v", err)
	}

	// Abre o crea el archivo `users.txt`
	file, err := os.OpenFile(ruta, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error al crear el archivo users.txt: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(contenido)
	if err != nil {
		return fmt.Errorf("error al escribir en users.txt: %v", err)
	}

	return nil
}

// Verifica si el grupo existe y está activo en users.txt
func grupoExisteActivo(grupo string, contenido string) bool {
	lineas := strings.Split(contenido, "\n")
	for _, linea := range lineas {
		datos := strings.Split(linea, ",")
		if len(datos) >= 3 && strings.TrimSpace(datos[1]) == "G" && strings.TrimSpace(datos[2]) == grupo && strings.TrimSpace(datos[0]) != "0" {
			return true
		}
	}
	return false
}

// Verifica si el usuario existe y está activo en users.txt, sin importar el grupo
func usuarioExisteActivo(usuario string, contenido string) bool {
	lineas := strings.Split(contenido, "\n")
	for _, linea := range lineas {
		datos := strings.Split(linea, ",")
		if len(datos) >= 5 && strings.TrimSpace(datos[1]) == "U" && strings.TrimSpace(datos[3]) == usuario && strings.TrimSpace(datos[0]) != "0" {
			return true
		}
	}
	return false
}
