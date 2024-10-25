package commands

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"
)

// Parámetros del comando mkfs
type ParametrosMkfs struct {
	Id   string
	Type string
}

// Analiza los parámetros del comando mkfs
func AnalizarParametrosMkfs(comando string) (ParametrosMkfs, error) {
	parametros := ParametrosMkfs{Type: "full"} // Tipo por defecto es "full"

	for _, parte := range strings.Fields(comando) {
		if strings.HasPrefix(strings.ToLower(parte), "-id=") {
			parametros.Id = strings.TrimPrefix(parte, "-id=")
		} else if strings.HasPrefix(strings.ToLower(parte), "-type=") {
			parametros.Type = strings.ToLower(strings.TrimPrefix(parte, "-type="))
		}
	}

	if parametros.Id == "" {
		return parametros, fmt.Errorf("el parámetro id es obligatorio")
	}

	return parametros, nil
}

// Ejecuta el comando mkfs para formatear una partición en EXT2
func EjecutarMkfs(parametros ParametrosMkfs) string {
	particion, existe := particionesMontadas[parametros.Id]
	if !existe {
		return fmt.Sprintf("Error: No se encontró la partición con el id %s", parametros.Id)
	}

	// Abre el archivo de disco
	archivo, err := os.OpenFile(particion.Ruta, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Sprintf("Error al abrir el archivo de disco: %v", err)
	}
	defer archivo.Close()

	// Formatear la partición en EXT2 y crear el archivo `users.txt`
	if err := FormatearEXT2(archivo, particion); err != nil {
		return fmt.Sprintf("Error al formatear la partición: %v", err)
	}
	if err := CrearArchivoUsersTxt(particion); err != nil {
		return fmt.Sprintf("Error al crear users.txt: %v", err)
	}

	return fmt.Sprintf("Partición %s formateada exitosamente en EXT2", parametros.Id)
}

// Inicializa la estructura del sistema de archivos EXT2 en la partición
func FormatearEXT2(archivo *os.File, particion ParticionMontada) error {
	superBlock := SuperBlock{
		S_filesystem_type:   0xEF53,
		S_inodes_count:      1024,
		S_blocks_count:      3072,
		S_free_blocks_count: 3072,
		S_free_inodes_count: 1024,
		S_mnt_count:         0,
		S_magic:             0xEF53,
		S_inode_size:        int32(binary.Size(Inode{})),
		S_block_size:        64,
		S_first_ino:         2,
		S_first_blo:         1,
	}
	copy(superBlock.S_mtime[:], time.Now().Format("20060102150405"))

	archivo.Seek(int64(particion.Inicio), 0)
	if err := binary.Write(archivo, binary.LittleEndian, &superBlock); err != nil {
		return fmt.Errorf("error al escribir el súper bloque: %v", err)
	}

	return nil
}

// CrearArchivoUsersTxt crea el archivo users.txt con los datos iniciales en la partición
func CrearArchivoUsersTxt(particion ParticionMontada) error {
	rutaUsersTxt := obtenerRutaUsersTxt(particion)

	// Crear o abrir el archivo users.txt en la ruta simulada
	file, err := os.OpenFile(rutaUsersTxt, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error al crear el archivo users.txt: %v", err)
	}
	defer file.Close()

	// Escribir contenido inicial en users.txt
	contenido := "1,G,root\n1,U,root,root,123\n"
	if _, err := file.WriteString(contenido); err != nil {
		return fmt.Errorf("error al escribir en users.txt: %v", err)
	}

	return nil
}
