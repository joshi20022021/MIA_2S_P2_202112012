package commands

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Definición de los parámetros para mkdisk
type ParametrosMkDisk struct {
	Tamano int64
	Ajuste string
	Unidad string
	Ruta   string
}

// Analiza los parámetros del comando mkdisk
func AnalizarParametrosMkDisk(comando string) (ParametrosMkDisk, error) {
	parametros := ParametrosMkDisk{Ajuste: "FF", Unidad: "M"}

	var parte string
	comando = strings.TrimSpace(comando)
	for comando != "" {
		index := strings.Index(comando, " ")
		if index != -1 {
			parte = comando[:index]
			comando = strings.TrimSpace(comando[index+1:])
		} else {
			parte = comando
			comando = ""
		}

		if strings.HasPrefix(parte, "-size=") {
			tamanoStr := strings.TrimPrefix(parte, "-size=")
			tamano, err := strconv.ParseInt(tamanoStr, 10, 64)
			if err != nil || tamano <= 0 {
				return parametros, fmt.Errorf("tamaño inválido")
			}
			parametros.Tamano = tamano
		} else if strings.HasPrefix(parte, "-fit=") {
			parametros.Ajuste = strings.ToUpper(strings.TrimPrefix(parte, "-fit="))
			if parametros.Ajuste != "BF" && parametros.Ajuste != "FF" && parametros.Ajuste != "WF" {
				return parametros, fmt.Errorf("ajuste inválido")
			}
		} else if strings.HasPrefix(parte, "-unit=") {
			parametros.Unidad = strings.ToUpper(strings.TrimPrefix(parte, "-unit="))
			if parametros.Unidad != "K" && parametros.Unidad != "M" {
				return parametros, fmt.Errorf("unidad inválida")
			}
		} else if strings.HasPrefix(parte, "-path=") {
			parametros.Ruta = strings.Trim(strings.TrimPrefix(parte, "-path="), "\"")
		}
	}

	if parametros.Tamano == 0 {
		return parametros, fmt.Errorf("el parámetro size es obligatorio")
	}
	if parametros.Ruta == "" {
		return parametros, fmt.Errorf("el parámetro path es obligatorio")
	}

	return parametros, nil
}

// Función que inicializa el MBR
func inicializarMBR(tamano int64, ajuste string) MBR {
	mbr := MBR{
		Mbr_tamano:        int32(tamano),
		Mbr_dsk_signature: generarNumeroAleatorio(),
	}

	// Asignar la fecha de creación del MBR
	timestamp := time.Now().Format("02-01-2006 15:04:05")
	copy(mbr.Mbr_fecha_creacion[:], timestamp)

	// Asignar el tipo de ajuste
	copy(mbr.Dsk_fit[:], ajuste)

	// Inicializar las particiones con valores por defecto
	for i := 0; i < 4; i++ {
		mbr.Mbr_partition[i] = Partition{
			Part_status: [1]byte{'0'}, // Partición no montada
			Part_type:   [1]byte{'0'}, // Sin tipo
			Part_fit:    [2]byte{'F'}, // Ajuste por defecto (FF - First Fit)
		}
	}

	return mbr
}

// Función que calcula el tamaño en bytes basado en la unidad
func calcularTamanoEnBytes(tamano int64, unidad string) int64 {
	if unidad == "K" {
		return tamano * 1024 // Convertir a kilobytes
	} else if unidad == "M" {
		return tamano * 1024 * 1024 // Convertir a megabytes
	}
	return tamano // Si no es K o M, se asume que ya está en bytes
}

// Después de escribir el MBR en el archivo, leerlo nuevamente para validar
func EjecutarMkDisk(parametros ParametrosMkDisk) string {
	archivo, err := os.Create(parametros.Ruta)
	if err != nil {
		return fmt.Sprintf("Error al crear el archivo: %v", err)
	}
	defer archivo.Close()

	tamanoEnBytes := calcularTamanoEnBytes(parametros.Tamano, parametros.Unidad)
	mbr := inicializarMBR(tamanoEnBytes, parametros.Ajuste)

	// Escribir el MBR
	err = binary.Write(archivo, binary.LittleEndian, &mbr)
	if err != nil {
		return fmt.Sprintf("Error al escribir el MBR en el archivo: %v", err)
	}

	// Leer el MBR nuevamente para validar la escritura
	archivo.Seek(0, 0) // Volver al inicio del archivo
	var mbrLeido MBR
	err = binary.Read(archivo, binary.LittleEndian, &mbrLeido)
	if err != nil {
		return fmt.Sprintf("Error al leer el MBR después de la escritura: %v", err)
	}

	// Validar si los datos coinciden
	if mbrLeido.Mbr_tamano != mbr.Mbr_tamano {
		return "Error: El MBR no fue escrito correctamente"
	}

	// Continuar con el resto de la lógica...
	return fmt.Sprintf("Disco creado exitosamente en %s", parametros.Ruta)
}

// Función para generar un número aleatorio que servirá como dsk_signature
func generarNumeroAleatorio() int32 {
	return int32(time.Now().UnixNano() & 0xFFFFFFFF)
}
