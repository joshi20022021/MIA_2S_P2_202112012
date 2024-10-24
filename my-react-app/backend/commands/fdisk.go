package commands

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

type ParametrosFDisk struct {
	Tamano int64
	Ruta   string
	Nombre string
	Tipo   string
	Ajuste string
	Unidad string
}

func AnalizarParametrosFDisk(comando string) (ParametrosFDisk, error) {
	parametros := ParametrosFDisk{Ajuste: "FF", Tipo: "P", Unidad: "K"}

	partes := strings.Split(comando, " ")
	if len(partes) == 0 || partes[0] != "fdisk" {
		return parametros, fmt.Errorf("comando no reconocido")
	}

	for i := 1; i < len(partes); i++ {
		parte := partes[i]
		if strings.HasPrefix(parte, "-size=") {
			fmt.Sscanf(strings.TrimPrefix(parte, "-size="), "%d", &parametros.Tamano)
		} else if strings.HasPrefix(parte, "-path=") {
			parametros.Ruta = strings.Trim(strings.TrimPrefix(parte, "-path="), "\"")
		} else if strings.HasPrefix(parte, "-name=") {
			parametros.Nombre = strings.TrimPrefix(parte, "-name=")
		} else if strings.HasPrefix(parte, "-type=") {
			parametros.Tipo = strings.ToUpper(strings.TrimPrefix(parte, "-type="))
		} else if strings.HasPrefix(parte, "-fit=") {
			parametros.Ajuste = strings.ToUpper(strings.TrimPrefix(parte, "-fit="))
		} else if strings.HasPrefix(parte, "-unit=") {
			parametros.Unidad = strings.ToUpper(strings.TrimPrefix(parte, "-unit="))
		}
	}

	if parametros.Tamano == 0 {
		return parametros, fmt.Errorf("el parámetro size es obligatorio")
	}
	if parametros.Ruta == "" {
		return parametros, fmt.Errorf("el parámetro path es obligatorio")
	}
	if parametros.Nombre == "" {
		return parametros, fmt.Errorf("el parámetro name es obligatorio")
	}

	return parametros, nil
}

func EjecutarFDisk(parametros ParametrosFDisk) string {
	file, err := os.OpenFile(parametros.Ruta, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Sprintf("Error al abrir el archivo: %v", err)
	}
	defer file.Close()

	var mbr MBR
	err = binary.Read(file, binary.LittleEndian, &mbr)
	if err != nil {
		return fmt.Sprintf("Error al leer el MBR: %v", err)
	}

	// Convertir el tamaño basado en la unidad especificada
	tamanoEnBytes := parametros.Tamano
	if parametros.Unidad == "K" {
		tamanoEnBytes *= 1024
	} else if parametros.Unidad == "M" {
		tamanoEnBytes *= 1024 * 1024
	}

	// Verificar si el tamaño de la partición es mayor que el disco
	if tamanoEnBytes > int64(mbr.Mbr_tamano) {
		return "Error: El tamaño de la partición es mayor que el tamaño del disco"
	}

	// Calcular espacio disponible
	espacioDisponible := int32(mbr.Mbr_tamano) - calcularEspacioUsado(mbr.Mbr_partition)
	if espacioDisponible < int32(tamanoEnBytes) {
		return "Error: No hay suficiente espacio disponible para la nueva partición"
	}

	// Buscar un espacio vacío para la nueva partición
	for i := range mbr.Mbr_partition {
		if mbr.Mbr_partition[i].Part_size == 0 {
			mbr.Mbr_partition[i].Part_size = int32(tamanoEnBytes)

			nombreParticion := strings.TrimSpace(parametros.Nombre)
			if len(nombreParticion) > 16 {
				nombreParticion = nombreParticion[:16]
			}
			copy(mbr.Mbr_partition[i].Part_name[:], nombreParticion)

			mbr.Mbr_partition[i].Part_type[0] = parametros.Tipo[0]
			copy(mbr.Mbr_partition[i].Part_fit[:], parametros.Ajuste[:2])

			// Calcular el inicio de la partición
			mbr.Mbr_partition[i].Part_start = calcularInicioParticion(mbr.Mbr_partition)

			// Si es extendida, inicializa EBR
			if parametros.Tipo == "E" {
				var ebr EBR
				ebr.Part_status[0] = '0'
				ebr.Part_start = mbr.Mbr_partition[i].Part_start + 1
				ebr.Part_next = -1
				ebr.Part_size = 0

				// Escribir el EBR en la partición extendida
				file.Seek(int64(mbr.Mbr_partition[i].Part_start), 0)
				err := binary.Write(file, binary.LittleEndian, &ebr)
				if err != nil {
					return fmt.Sprintf("Error al escribir el EBR: %v", err)
				}
			}

			// Escribir el MBR actualizado
			file.Seek(0, 0)
			err = binary.Write(file, binary.LittleEndian, &mbr)
			if err != nil {
				return fmt.Sprintf("Error al escribir el MBR: %v", err)
			}

			return fmt.Sprintf("Partición %s creada en %s", parametros.Nombre, parametros.Ruta)
		}
	}

	return "Error: No se pudo crear la partición, no hay espacio disponible"
}

// Función auxiliar para calcular el espacio usado por las particiones
func calcularEspacioUsado(particiones [4]Partition) int32 {
	espacioUsado := int32(0)
	for _, particion := range particiones {
		espacioUsado += particion.Part_size
	}
	return espacioUsado
}

// Función auxiliar para calcular el punto de inicio de una nueva partición
func calcularInicioParticion(particiones [4]Partition) int32 {
	inicio := int32(512) // Inicia justo después del MBR
	for _, particion := range particiones {
		if particion.Part_size > 0 {
			finParticion := particion.Part_start + particion.Part_size
			if finParticion > inicio {
				inicio = finParticion
			}
		}
	}
	return inicio
}
