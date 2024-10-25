package commands

// Ejecuta el comando logout para cerrar la sesión
func EjecutarLogout() string {
	if !VerificarSesionActiva() {
		return "Error: No hay ninguna sesión activa"
	}

	// Cerrar sesión y reiniciar las variables de sesión
	CerrarSesion()
	return "Sesión cerrada correctamente"
}
