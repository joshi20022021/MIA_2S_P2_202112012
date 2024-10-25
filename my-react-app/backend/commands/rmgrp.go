package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// obtenerRutaUsersTxt devuelve la ruta virtual de users.txt dentro de la partición montada
func obtenerRutaUsersTxt(particion ParticionMontada) string {
	// Usa un directorio simulado en la partición
	return fmt.Sprintf("%s_users.txt", particion.Ruta)
}

// Estructura de parámetros para el comando rmgrp
type ParametrosRmgrp struct {
	Nombre string
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
	file, err := os.OpenFile(ruta, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el archivo %s para escritura: %v", ruta, err)
	}
	defer file.Close()

	_, err = file.WriteString(contenido)
	if err != nil {
		return fmt.Errorf("error al escribir en %s: %v", ruta, err)
	}

	return nil
}

// Analizar el comando rmgrp
func AnalizarParametrosRmgrp(comando string) (ParametrosRmgrp, error) {
	parametros := ParametrosRmgrp{}

	for _, parte := range strings.Fields(comando) {
		if strings.HasPrefix(strings.ToLower(parte), "-name=") {
			parametros.Nombre = strings.TrimPrefix(parte, "-name=")
		}
	}

	if parametros.Nombre == "" {
		return parametros, fmt.Errorf("el parámetro -name es obligatorio")
	}

	return parametros, nil
}
func EjecutarRmgrp(parametros ParametrosRmgrp) string {
	if !VerificarSesionActiva() || UsuarioLogueado() != "root" {
		return "Error: solo el usuario root puede eliminar grupos o no hay sesión activa"
	}

	rutaUsersTxt := obtenerRutaUsersTxt(ParticionActiva())

	contenido, err := leerUsersTxtDesdeDisco(rutaUsersTxt)
	if err != nil {
		return fmt.Sprintf("Error al leer users.txt: %v", err)
	}

	lineas := strings.Split(contenido, "\n")
	for i, linea := range lineas {
		datos := strings.Split(linea, ",")
		if len(datos) > 2 && strings.TrimSpace(datos[1]) == "G" && strings.TrimSpace(datos[2]) == parametros.Nombre {
			if datos[0] == "0" {
				return fmt.Sprintf("Error: el grupo %s ya fue eliminado", parametros.Nombre)
			}
			datos[0] = "0"
			lineas[i] = strings.Join(datos, ",")
		}
	}

	contenidoActualizado := strings.Join(lineas, "\n")
	if err := escribirUsersTxtEnDisco(rutaUsersTxt, contenidoActualizado); err != nil {
		return fmt.Sprintf("Error al escribir en users.txt: %v", err)
	}

	return fmt.Sprintf("Grupo %s eliminado exitosamente", parametros.Nombre)
}
