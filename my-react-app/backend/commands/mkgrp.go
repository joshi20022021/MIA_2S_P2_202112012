package commands

import (
	"fmt"
	"strings"
)

// Estructura para los parámetros del comando mkgrp
type ParametrosMkgrp struct {
	Nombre string
}

// Función para analizar los parámetros de mkgrp
func AnalizarParametrosMkgrp(comando string) (ParametrosMkgrp, error) {
	parametros := ParametrosMkgrp{}
	for _, parte := range strings.Fields(comando) {
		if strings.HasPrefix(strings.ToLower(parte), "-name=") {
			parametros.Nombre = strings.TrimPrefix(parte, "-name=")
		}
	}
	if parametros.Nombre == "" {
		return parametros, fmt.Errorf("el parámetro name es obligatorio")
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

	lineas := strings.Split(contenido, "\n")
	gid := 1
	for _, linea := range lineas {
		datos := strings.Split(linea, ",")
		if len(datos) == 3 && strings.TrimSpace(datos[1]) == "G" {
			gid++
			if strings.TrimSpace(datos[2]) == parametros.Nombre {
				return fmt.Sprintf("Error: el grupo %s ya existe", parametros.Nombre)
			}
		}
	}

	nuevaLinea := fmt.Sprintf("%d,G,%s\n", gid, parametros.Nombre)
	contenidoActualizado := contenido + nuevaLinea

	if err := escribirUsersTxtEnDisco(rutaUsersTxt, contenidoActualizado); err != nil {
		return fmt.Sprintf("Error al escribir en users.txt: %v", err)
	}

	return fmt.Sprintf("Grupo %s creado exitosamente", parametros.Nombre)
}
