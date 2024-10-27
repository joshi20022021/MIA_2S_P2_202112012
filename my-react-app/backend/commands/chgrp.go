// chgrp.go
package commands

import (
	"fmt"
	"strings"
)

type ParametrosChgrp struct {
	User string
	Grp  string
}

// Analizar los parámetros del comando chgrp
func AnalizarParametrosChgrp(comando string) (ParametrosChgrp, error) {
	parametros := ParametrosChgrp{}
	for _, parte := range strings.Fields(comando) {
		if strings.HasPrefix(strings.ToLower(parte), "-user=") {
			parametros.User = strings.TrimPrefix(parte, "-user=")
		} else if strings.HasPrefix(strings.ToLower(parte), "-grp=") {
			parametros.Grp = strings.TrimPrefix(parte, "-grp=")
		}
	}

	if parametros.User == "" || parametros.Grp == "" {
		return parametros, fmt.Errorf("los parámetros -user y -grp son obligatorios")
	}
	return parametros, nil
}

func EjecutarChgrp(parametros ParametrosChgrp) string {
	if !VerificarSesionActiva() || UsuarioLogueado() != "root" {
		return "Error: solo el usuario root puede cambiar grupos de usuarios o no hay sesión activa"
	}

	rutaUsersTxt := obtenerRutaUsersTxt(ParticionActiva())
	contenido, err := leerUsersTxtDesdeDisco(rutaUsersTxt)
	if err != nil {
		return fmt.Sprintf("Error al leer users.txt: %v", err)
	}

	if !grupoExisteActivo(parametros.Grp, contenido) {
		return fmt.Sprintf("Error: el grupo %s no existe", parametros.Grp)
	}

	if !usuarioExisteActivo(parametros.User, contenido) {
		return fmt.Sprintf("Error: el usuario %s no existe", parametros.User)
	}

	lineas := strings.Split(contenido, "\n")
	for i, linea := range lineas {
		datos := strings.Split(linea, ",")
		if len(datos) == 5 && strings.TrimSpace(datos[1]) == "U" && strings.TrimSpace(datos[3]) == parametros.User {
			datos[2] = parametros.Grp // Actualizar el grupo del usuario
			lineas[i] = strings.Join(datos, ",")
			break
		}
	}

	contenidoActualizado := strings.Join(lineas, "\n")
	if err := escribirUsersTxtEnDisco(rutaUsersTxt, contenidoActualizado); err != nil {
		return fmt.Sprintf("Error al escribir en users.txt: %v", err)
	}

	return fmt.Sprintf("Grupo del usuario %s cambiado exitosamente a %s", parametros.User, parametros.Grp)
}
