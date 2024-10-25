package commands

import (
	"fmt"
	"strings"
)

type ParametrosRmgrp struct {
	Nombre string
}

// Analizar los parámetros del comando rmgrp
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
