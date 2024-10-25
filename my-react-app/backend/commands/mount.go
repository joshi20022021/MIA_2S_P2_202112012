package commands

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

var contadorLetra = make(map[string]rune)                   // Controla la letra para cada disco montado
var montajesDiscos = map[string]int{}                       // Controla el número de particiones montadas por disco
var particionesMontadas = make(map[string]ParticionMontada) // Mapa global de particiones montadas
const sufijoCarnet = "12"                                   // Últimos dos dígitos del carnet

// Estructura para almacenar la información de una partición montada
type ParticionMontada struct {
	Id     string
	Ruta   string
	Inicio int32 // Inicio de la partición
	Tipo   byte  // Tipo de partición (P para primaria, L para lógica)
}

type ParametrosMontar struct {
	Ruta   string
	Nombre string
}

// Analiza los parámetros para el comando mount
func AnalizarParametrosMontarNew(comando string) (ParametrosMontar, error) {
	parametros := ParametrosMontar{}

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
				return montarPrimaria(&mbr, i, parametros)
			} else if particion.Part_type[0] == 'E' { // Si es extendida, buscar lógicas
				return montarLogica(archivo, particion.Part_start, parametros)
			}
		}
	}

	return fmt.Sprintf("Error: No se encontró la partición con el nombre %s", parametros.Nombre)
}

func montarPrimaria(mbr *MBR, index int, parametros ParametrosMontar) string {
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

	// Marcar la partición como montada (simulación en memoria)
	particion.Part_status[0] = '1'
	copy(particion.Part_id[:], id)

	// Agregar la partición montada al mapa global
	particionesMontadas[id] = ParticionMontada{
		Id:     id,
		Ruta:   parametros.Ruta,
		Inicio: particion.Part_start,
		Tipo:   particion.Part_type[0],
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

		// Marcar la partición lógica como montada en memoria
		ebr.Part_status[0] = '1'
		copy(ebr.Part_id[:], id)

		// Agregar la partición montada al mapa global
		particionesMontadas[id] = ParticionMontada{
			Id:     id,
			Ruta:   parametros.Ruta,
			Inicio: ebr.Part_start,
			Tipo:   'L',
		}

		return fmt.Sprintf("Partición lógica %s montada con éxito. ID generado: %s", parametros.Nombre, id)
	}

	return fmt.Sprintf("Error: No se encontró la partición lógica con el nombre %s", parametros.Nombre)
}

// Genera el ID único para el montaje usando el último número del carnet y reglas de incrementación
func agregarMontaje(parametros ParametrosMontar) (string, error) {
	// Asignar letra al disco si es nuevo
	if _, existe := contadorLetra[parametros.Ruta]; !existe {
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
		montajesDiscos[parametros.Ruta] = 1 // Inicia el conteo de particiones
	}

	letra := contadorLetra[parametros.Ruta]
	numeroParticion := montajesDiscos[parametros.Ruta]

	if numeroParticion < 1 || numeroParticion > 9 {
		return "", fmt.Errorf("el número de partición debe estar entre 1 y 9")
	}

	// Generar ID: 12 + número de partición (1-9) + letra del disco
	id := fmt.Sprintf("%s%d%c", sufijoCarnet, numeroParticion, letra)
	montajesDiscos[parametros.Ruta]++ // Incrementa para próxima partición

	return id, nil
}
