package commands

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

type ParametrosRmDisk struct {
	Ruta string
}

// Analiza los parámetros del comando rmdisk
func AnalizarParametrosRmDisk(comando string) (ParametrosRmDisk, error) {
	parametros := ParametrosRmDisk{}
	partes := strings.Split(comando, " ")

	if len(partes) == 0 || partes[0] != "rmdisk" {
		return parametros, fmt.Errorf("comando no reconocido")
	}

	for _, parte := range partes {
		if strings.HasPrefix(parte, "-path=") {
			parametros.Ruta = strings.Trim(strings.TrimPrefix(parte, "-path="), "\"")
		}
	}

	if parametros.Ruta == "" {
		return parametros, fmt.Errorf("el parámetro path es obligatorio")
	}

	return parametros, nil
}

// Verifica si el archivo contiene un MBR antes de eliminarlo
func EjecutarRmDisk(parametros ParametrosRmDisk) string {
	if _, err := os.Stat(parametros.Ruta); os.IsNotExist(err) {
		return fmt.Sprintf("Error: El archivo %s no existe", parametros.Ruta)
	}

	// Intentar abrir el archivo para verificar el MBR
	archivo, err := os.Open(parametros.Ruta)
	if err != nil {
		return fmt.Sprintf("Error al abrir el archivo: %v", err)
	}
	defer archivo.Close()

	// Leer el MBR para verificar si el archivo contiene un MBR válido
	var mbr MBR
	err = binary.Read(archivo, binary.LittleEndian, &mbr)
	if err != nil {
		return fmt.Sprintf("Error: El archivo %s no contiene un MBR válido", parametros.Ruta)
	}

	// Confirmar que es un disco válido antes de eliminarlo
	if mbr.Mbr_tamano == 0 || mbr.Mbr_dsk_signature == 0 {
		return fmt.Sprintf("Error: El archivo %s no contiene un disco válido", parametros.Ruta)
	}

	// Eliminar el archivo si es un disco válido
	if err := os.Remove(parametros.Ruta); err != nil {
		return fmt.Sprintf("Error al eliminar el archivo: %v", err)
	}

	return fmt.Sprintf("Archivo %s eliminado exitosamente", parametros.Ruta)
}
