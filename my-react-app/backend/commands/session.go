package commands

// Variables de sesión globales
var sesionActiva bool
var usuarioLogueado string
var particionActiva ParticionMontada

// Función para iniciar sesión
func IniciarSesion(usuario string, particion ParticionMontada) {
	sesionActiva = true
	usuarioLogueado = usuario
	particionActiva = particion
}

// Verifica si hay una sesión activa
func VerificarSesionActiva() bool {
	return sesionActiva
}

// Cierra la sesión
func CerrarSesion() {
	sesionActiva = false
	usuarioLogueado = ""
	particionActiva = ParticionMontada{}
}

// Retorna el nombre del usuario logueado
func UsuarioLogueado() string {
	return usuarioLogueado
}

// Retorna la partición activa
func ParticionActiva() ParticionMontada {
	return particionActiva
}
