package commands

import (
	"fmt"
	"os"
	"strings"
)

// Parámetros del comando MKGRP
type ParametrosMkgrp struct {
	Nombre string
}

// Analizar el comando mkgrp
func AnalizarParametrosMkgrp(comando string) (ParametrosMkgrp, error) {
	parametros := ParametrosMkgrp{}
	partes := strings.Split(comando, " ")

	if len(partes) == 0 || partes[0] != "mkgrp" {
		return parametros, fmt.Errorf("comando no reconocido")
	}

	for i := 1; i < len(partes); i++ {
		parte := partes[i]
		if strings.HasPrefix(parte, "-name=") {
			parametros.Nombre = strings.TrimPrefix(parte, "-name=")
		}
	}

	if parametros.Nombre == "" {
		return parametros, fmt.Errorf("el nombre del grupo es obligatorio")
	}

	return parametros, nil
}
func EjecutarMkgrp(parametros ParametrosMkgrp) string {
	// Verificar si hay sesión activa
	if !VerificarSesionActiva() {
		return "Error: no hay sesión iniciada"
	}

	// Verificar si el usuario es root
	if UsuarioLogueado() != "root" {
		return "Error: solo el usuario root puede crear grupos"
	}

	// Obtener la ruta de users.txt en la partición montada
	rutaUsers := obtenerRutaUsersTxt(particionMontada)
	if rutaUsers == "" {
		return fmt.Sprintf("Error: no se encontró la partición montada con ID %s", particionMontada)
	}

	// Leer el contenido actual del archivo users.txt
	contenidoUsers, err := leerUsersTxtDesdeDisco(rutaUsers)
	if err != nil {
		return fmt.Sprintf("Error al leer el archivo users.txt: %v", err)
	}

	// Verificar si el grupo ya existe
	lineas := strings.Split(contenidoUsers, "\n")
	for _, linea := range lineas {
		if strings.Contains(linea, fmt.Sprintf(",G,%s", parametros.Nombre)) {
			return fmt.Sprintf("Error: el grupo %s ya existe", parametros.Nombre)
		}
	}

	// Calcular el nuevo GID
	gid := len(lineas) / 2 // Asume que cada grupo y usuario ocupa dos líneas
	nuevaLinea := fmt.Sprintf("%d,G,%s\n", gid+1, parametros.Nombre)

	// Agregar la nueva línea al archivo users.txt
	contenidoUsers += nuevaLinea
	err = escribirUsersTxtEnDisco(rutaUsers, contenidoUsers)
	if err != nil {
		return fmt.Sprintf("Error al escribir el archivo users.txt: %v", err)
	}

	// Mostrar el contenido actualizado de users.txt
	fmt.Printf("Contenido actualizado de users.txt:\n%s\n", contenidoUsers)

	return fmt.Sprintf("Grupo %s creado exitosamente", parametros.Nombre)
}
func escribirUsersTxtEnDisco(ruta string, contenido string) error {
	// Abrir el archivo .mia para escribir en él
	file, err := os.OpenFile(ruta, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error al abrir el disco %s: %v", ruta, err)
	}
	defer file.Close()

	// Escribir el contenido en la posición específica donde se encuentra users.txt
	offset := int64(2048) // Asegúrate de usar el mismo offset de creación
	_, err = file.WriteAt([]byte(contenido), offset)
	if err != nil {
		return fmt.Errorf("error al escribir en el disco %s: %v", ruta, err)
	}

	fmt.Printf("Contenido del archivo users.txt actualizado correctamente en el disco %s\n", ruta)
	return nil
}
