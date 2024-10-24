package commands

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

var contadorLetra = make(map[string]rune) // Controla la letra para cada disco montado
var montajesDiscos = map[string]int{}     // Controla el número de particiones montadas por disco
const sufijoCarnet = "12"                 // Últimos dos dígitos del carnet

type ParametrosMontar struct {
	Ruta   string
	Nombre string
}

// Analizar el comando mount sin usar listas
func AnalizarParametrosMontarNew(comando string) (ParametrosMontar, error) {
	parametros := ParametrosMontar{}

	// Separar el comando en sus partes
	for _, parte := range strings.Fields(comando) {
		if strings.HasPrefix(strings.ToLower(parte), "-path=") {
			parametros.Ruta = strings.Trim(strings.TrimPrefix(parte, "-path="), "\"")
		} else if strings.HasPrefix(strings.ToLower(parte), "-name=") {
			parametros.Nombre = strings.TrimPrefix(parte, "-name=")
		}
	}

	if parametros.Ruta == "" || parametros.Nombre == "" {
		return parametros, fmt.Errorf("los parámetros path y name son obligatorios")
	}

	return parametros, nil
}

func EjecutarMontar(parametros ParametrosMontar) string {
	archivo, err := os.OpenFile(parametros.Ruta, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Sprintf("Error al abrir el archivo: %v", err)
	}
	defer archivo.Close()

	var mbr MBR
	err = binary.Read(archivo, binary.LittleEndian, &mbr)
	if err != nil {
		return fmt.Sprintf("Error al leer el MBR: %v", err)
	}

	// Buscar la partición por nombre
	for i := 0; i < len(mbr.Mbr_partition); i++ {
		particion := mbr.Mbr_partition[i]
		nombreParticion := strings.TrimSpace(strings.TrimRight(string(particion.Part_name[:]), "\x00"))
		if nombreParticion == strings.TrimSpace(parametros.Nombre) {
			if particion.Part_type[0] == 'P' { // Verificar si es una partición primaria
				return montarPrimaria(&mbr, archivo, i, parametros)
			} else if particion.Part_type[0] == 'E' { // Si es extendida, buscar lógicas
				return montarLogica(archivo, particion.Part_start, parametros)
			}
		}
	}

	return fmt.Sprintf("Error: No se encontró la partición con el nombre %s", parametros.Nombre)
}

// Montar una partición primaria y generar un ID único
// Montar una partición primaria y generar un ID único
func montarPrimaria(mbr *MBR, archivo *os.File, index int, parametros ParametrosMontar) string {
	particion := mbr.Mbr_partition[index]

	// Verificar si la partición ya está montada
	if particion.Part_status[0] == '1' {
		return fmt.Sprintf("Error: La partición %s ya está montada", parametros.Nombre)
	}

	// Generar el ID del montaje
	id, err := agregarMontaje(parametros)
	if err != nil {
		return fmt.Sprintf("Error al generar el ID de montaje: %v", err)
	}

	// Asegúrate de que el campo Part_id se actualice correctamente
	copy(particion.Part_id[:], id)

	// Marcar la partición como montada
	particion.Part_status[0] = '1'
	archivo.Seek(0, 0)
	err = binary.Write(archivo, binary.LittleEndian, mbr)
	if err != nil {
		return fmt.Sprintf("Error al escribir el MBR actualizado: %v", err)
	}

	// Leer el MBR nuevamente para validar que el ID se escribió correctamente
	var mbrVerificado MBR
	archivo.Seek(0, 0)
	err = binary.Read(archivo, binary.LittleEndian, &mbrVerificado)
	if err != nil {
		return fmt.Sprintf("Error al verificar el MBR después del montaje: %v", err)
	}

	// Validar si el ID se escribió correctamente
	if string(mbrVerificado.Mbr_partition[index].Part_id[:]) != string(particion.Part_id[:]) {
		return "Error: El ID de la partición no se escribió correctamente"
	}

	return fmt.Sprintf("Partición %s montada con éxito. ID generado: %s", parametros.Nombre, id)
}

// Montar una partición lógica y generar un ID único
func montarLogica(archivo *os.File, start int32, parametros ParametrosMontar) string {
	var ebr EBR
	archivo.Seek(int64(start), 0)
	err := binary.Read(archivo, binary.LittleEndian, &ebr)
	if err != nil {
		return fmt.Sprintf("Error al leer el EBR: %v", err)
	}

	nombreParticion := strings.TrimRight(string(ebr.Part_name[:]), "\x00")
	if strings.EqualFold(nombreParticion, parametros.Nombre) {
		if ebr.Part_status[0] == '1' {
			return fmt.Sprintf("Error: La partición %s ya está montada", parametros.Nombre)
		}

		// Generar el ID del montaje
		id, err := agregarMontaje(parametros)
		if err != nil {
			return fmt.Sprintf("Error al generar el ID de montaje: %v", err)
		}
		// Validar que el ID tenga exactamente 4 caracteres
		if len(id) != 4 {
			return fmt.Sprintf("Error: ID generado tiene longitud incorrecta: %s", id)
		}

		// Asignar el ID al campo Part_id de la partición lógica
		copy(ebr.Part_id[:], id)

		// Marcar la partición lógica como montada
		ebr.Part_status[0] = '1'
		archivo.Seek(int64(start), 0)
		err = binary.Write(archivo, binary.LittleEndian, &ebr)
		if err != nil {
			return fmt.Sprintf("Error al escribir el EBR actualizado: %v", err)
		}

		// Registrar el montaje en el mapeo persistente
		//err = RegistrarMontaje(id, parametros.Ruta)
		if err != nil {
			return fmt.Sprintf("Error al registrar el montaje: %v", err)
		}

		return fmt.Sprintf("Partición lógica %s montada con éxito. ID generado: %s", parametros.Nombre, id)

	}

	return fmt.Sprintf("Error: No se encontró la partición lógica con el nombre %s", parametros.Nombre)
}

// Generar el ID de montaje
func agregarMontaje(parametros ParametrosMontar) (string, error) {
	// Asignar la letra al disco si es un nuevo disco
	if _, existe := contadorLetra[parametros.Ruta]; !existe {
		// Asignar la siguiente letra disponible
		if len(contadorLetra) == 0 {
			contadorLetra[parametros.Ruta] = 'A'
		} else {
			var ultimaLetra rune = 'A'
			for _, letra := range contadorLetra {
				if letra > ultimaLetra {
					ultimaLetra = letra
				}
			}
			if ultimaLetra >= 'Z' {
				return "", fmt.Errorf("se han alcanzado el límite de letras para discos")
			}
			contadorLetra[parametros.Ruta] = ultimaLetra + 1
		}
		montajesDiscos[parametros.Ruta] = 1 // Iniciar contador de particiones
	}

	letra := contadorLetra[parametros.Ruta]

	// Obtener el número actual de partición
	numeroParticion := montajesDiscos[parametros.Ruta]

	// Validar que el número de partición esté entre 1 y 9
	if numeroParticion < 1 || numeroParticion > 9 {
		return "", fmt.Errorf("el número de partición debe estar entre 1 y 9")
	}

	// Generar el ID: 12 (fijo) + número de partición (1-9) + letra del disco
	id := fmt.Sprintf("%s%d%c", sufijoCarnet, numeroParticion, letra)

	// Incrementar el contador para la próxima partición
	montajesDiscos[parametros.Ruta]++

	return id, nil
}
