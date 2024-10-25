package commands

import (
	"fmt"
	"strings"
)

type ParametrosMkgrp struct {
	Nombre string
}

// Analiza los parámetros del comando mkgrp
func AnalizarParametrosMkgrp(comando string) (ParametrosMkgrp, error) {
	parametros := ParametrosMkgrp{}

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

func EjecutarMkgrp(parametros ParametrosMkgrp) string {
	if !VerificarSesionActiva() || UsuarioLogueado() != "root" {
		return "Error: solo el usuario root puede crear grupos o no hay sesión activa"
	}

	rutaUsersTxt := obtenerRutaUsersTxt(ParticionActiva())

	contenido, err := leerUsersTxtDesdeDisco(rutaUsersTxt)
	if err != nil {
		return fmt.Sprintf("Error al leer users.txt: %v", err)
	}

	// Verificar si el grupo ya existe y agregarlo si no existe
	lineas := strings.Split(contenido, "\n")
	for _, linea := range lineas {
		datos := strings.Split(linea, ",")
		if len(datos) == 3 && strings.TrimSpace(datos[1]) == "G" && strings.TrimSpace(datos[2]) == parametros.Nombre {
			return fmt.Sprintf("Error: el grupo %s ya existe", parametros.Nombre)
		}
	}

	// Agregar el nuevo grupo
	contenido += fmt.Sprintf("0,G,%s\n", parametros.Nombre)

	// Escribir el contenido actualizado
	if err := escribirUsersTxtEnDisco(rutaUsersTxt, contenido); err != nil {
		return fmt.Sprintf("Error al escribir en users.txt: %v", err)
	}

	return fmt.Sprintf("Grupo %s creado exitosamente", parametros.Nombre)
}
