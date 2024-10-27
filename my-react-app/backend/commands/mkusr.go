// mkusr.go
package commands

import (
	"fmt"
	"strings"
)

// Estructura de parámetros para el comando mkusr
type ParametrosMkusr struct {
	User string
	Pass string
	Grp  string
}

// Analizar el comando mkusr
func AnalizarParametrosMkusr(comando string) (ParametrosMkusr, error) {
	parametros := ParametrosMkusr{}
	for _, parte := range strings.Fields(comando) {
		if strings.HasPrefix(parte, "-user=") {
			parametros.User = strings.TrimPrefix(parte, "-user=")
		} else if strings.HasPrefix(parte, "-pass=") {
			parametros.Pass = strings.TrimPrefix(parte, "-pass=")
		} else if strings.HasPrefix(parte, "-grp=") {
			parametros.Grp = strings.TrimPrefix(parte, "-grp=")
		}
	}

	if parametros.User == "" || parametros.Pass == "" || parametros.Grp == "" {
		return parametros, fmt.Errorf("los parámetros -user, -pass y -grp son obligatorios")
	}
	if len(parametros.User) > 10 || len(parametros.Pass) > 10 || len(parametros.Grp) > 10 {
		return parametros, fmt.Errorf("los parámetros -user, -pass y -grp deben tener un máximo de 10 caracteres")
	}
	return parametros, nil
}

func EjecutarMkusr(parametros ParametrosMkusr) string {
	if !VerificarSesionActiva() || UsuarioLogueado() != "root" {
		return "Error: solo el usuario root puede crear usuarios o no hay sesión activa"
	}

	rutaUsersTxt := obtenerRutaUsersTxt(ParticionActiva())
	contenido, err := leerUsersTxtDesdeDisco(rutaUsersTxt)
	if err != nil {
		return fmt.Sprintf("Error al leer users.txt: %v", err)
	}

	if !grupoExisteActivo(parametros.Grp, contenido) {
		return fmt.Sprintf("Error: el grupo %s no existe", parametros.Grp)
	}

	if usuarioExisteActivo(parametros.User, contenido) {
		return fmt.Sprintf("Error: el usuario \"%s\" ya existe", parametros.User)
	}

	nuevoUID := 1
	lineas := strings.Split(contenido, "\n")

	for _, linea := range lineas {
		datos := strings.Split(linea, ",")
		if len(datos) == 5 && strings.TrimSpace(datos[1]) == "U" {
			nuevoUID++
		}
	}

	nuevaLinea := fmt.Sprintf("%d,U,%s,%s,%s\n", nuevoUID, parametros.Grp, parametros.User, parametros.Pass)
	contenido += nuevaLinea

	if err := escribirUsersTxtEnDisco(rutaUsersTxt, contenido); err != nil {
		return fmt.Sprintf("Error al escribir en users.txt: %v", err)
	}

	return fmt.Sprintf("Usuario \"%s\" creado exitosamente en el grupo %s", parametros.User, parametros.Grp)
}
