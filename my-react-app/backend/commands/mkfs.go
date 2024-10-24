package commands

import (
	"fmt"
	"os"
	"strings"
)

type ParametrosMkfs struct {
	Tipo string
	ID   string
}

// Analizar los parámetros del comando mkfs
func AnalizarParametrosMkfs(comando string) (ParametrosMkfs, error) {
	parametros := ParametrosMkfs{Tipo: "full"} // "full" es el valor por defecto

	partes := strings.Split(comando, " ")
	if len(partes) == 0 || partes[0] != "mkfs" {
		return parametros, fmt.Errorf("comando no reconocido")
	}

	for i := 1; i < len(partes); i++ {
		parte := partes[i]
		if strings.HasPrefix(parte, "-type=") {
			parametros.Tipo = strings.ToLower(strings.TrimPrefix(parte, "-type="))
		} else if strings.HasPrefix(parte, "-id=") {
			parametros.ID = strings.TrimPrefix(parte, "-id=")
		}
	}

	if parametros.ID == "" {
		return parametros, fmt.Errorf("el id de la partición es obligatorio")
	}

	return parametros, nil
}

// Ejecutar mkfs para formatear una partición
func EjecutarMkfs(parametros ParametrosMkfs) string {
	// Verificar si la partición está montada
	montaje, encontrado := montajes[parametros.ID]
	if !encontrado {
		return fmt.Sprintf("error: la partición con ID %s no está montada", parametros.ID)
	}

	// Simulación del formateo de la partición
	fmt.Printf("Formateando la partición %s como %s\n", parametros.ID, parametros.Tipo)

	// Crear el archivo users.txt en la raíz de la partición
	err := crearUsersTxt(montaje.Ruta)
	if err != nil {
		return fmt.Sprintf("error al crear el archivo users.txt: %v", err)
	}

	return fmt.Sprintf("Partición %s formateada como %s y users.txt creado", parametros.ID, parametros.Tipo)
}

// Crear el archivo users.txt en el disco
func crearUsersTxt(ruta string) error {
	// Simular la creación del archivo users.txt dentro de la partición
	usersContent := "1,G,root\n1,U,root,root,123\n"
	return crearArchivoDentroDelDisco(ruta, usersContent, "users.txt")
}

func crearArchivoDentroDelDisco(ruta string, contenido string, nombreArchivo string) error {
	// Abrir el archivo .mia para escribir en él
	file, err := os.OpenFile(ruta, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error al abrir el disco %s: %v", ruta, err)
	}
	defer file.Close()

	// Escribir el contenido en una posición específica
	offset := int64(2048)
	_, err = file.WriteAt([]byte(contenido), offset)
	if err != nil {
		return fmt.Errorf("error al escribir en el disco %s: %v", ruta, err)
	}

	fmt.Printf("Contenido del archivo %s escrito correctamente en el disco %s\n", nombreArchivo, ruta)
	return nil
}
