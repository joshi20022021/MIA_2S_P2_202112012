package commands

// Variables de sesión globales
var sesionActiva bool                // Indica si hay una sesión activa
var usuarioLogueado string           // Nombre del usuario actualmente logueado
var particionActiva ParticionMontada // Partición actualmente activa

// Función para iniciar sesión
func IniciarSesion(usuario string, particion ParticionMontada) {
	sesionActiva = true
	usuarioLogueado = usuario
	particionActiva = particion
}

// Función para verificar si hay una sesión activa
func VerificarSesionActiva() bool {
	return sesionActiva
}

// Función para cerrar sesión
func CerrarSesion() {
	sesionActiva = false
	usuarioLogueado = ""
	particionActiva = ParticionMontada{}
}

// Retorna el nombre del usuario logueado (útil para otros comandos)
func UsuarioLogueado() string {
	return usuarioLogueado
}

// Retorna la partición activa (útil para otros comandos)
func ParticionActiva() ParticionMontada {
	return particionActiva
}
