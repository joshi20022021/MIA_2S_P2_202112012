package commands

// Definición de la estructura Partition
type Partition struct {
	Part_status [1]byte
	Part_type   [1]byte
	Part_fit    [2]byte
	Part_start  int32
	Part_size   int32
	Part_name   [16]byte
	Part_id     [4]byte // Nuevo campo agregado
}

// Definición de la estructura MBR (Master Boot Record)
type MBR struct {
	Mbr_tamano         int32
	Mbr_fecha_creacion [19]byte
	Mbr_dsk_signature  int32
	Dsk_fit            [1]byte
	Mbr_partition      [4]Partition
}

// Definición de la estructura EBR (Extended Boot Record)
type EBR struct {
	Part_status [1]byte
	Part_fit    [2]byte
	Part_start  int32
	Part_size   int32
	Part_next   int32
	Part_name   [16]byte
	Part_id     [4]byte // Nuevo campo agregado
}

// Definición del súper bloque para EXT2
type SuperBlock struct {
	S_filesystem_type   int32    // Identifica el sistema de archivos (EXT2 = 0xEF53)
	S_inodes_count      int32    // Total de inodos
	S_blocks_count      int32    // Total de bloques
	S_free_blocks_count int32    // Bloques libres
	S_free_inodes_count int32    // Inodos libres
	S_mtime             [19]byte // Última fecha de montaje
	S_umtime            [19]byte // Última fecha de desmontaje
	S_mnt_count         int32    // Número de veces montado
	S_magic             int32    // Número mágico EXT2
	S_inode_size        int32    // Tamaño del inodo
	S_block_size        int32    // Tamaño del bloque
	S_first_ino         int32    // Primer inodo libre
	S_first_blo         int32    // Primer bloque libre
	S_bm_inode_start    int32    // Inicio del bitmap de inodos
	S_bm_block_start    int32    // Inicio del bitmap de bloques
	S_inode_start       int32    // Inicio de la tabla de inodos
	S_block_start       int32    // Inicio de la tabla de bloques
}

// Definición de un inodo en EXT2
type Inode struct {
	I_uid   int32     // UID del propietario
	I_gid   int32     // GID del grupo
	I_size  int32     // Tamaño del archivo en bytes
	I_atime [19]byte  // Fecha de último acceso
	I_ctime [19]byte  // Fecha de creación
	I_mtime [19]byte  // Fecha de última modificación
	I_block [15]int32 // Punteros a bloques (12 directos, 1 indirecto simple, 1 indirecto doble, 1 indirecto triple)
	I_type  byte      // Tipo de archivo (0 = directorio, 1 = archivo)
	I_perm  [3]byte   // Permisos UGO en octal
}

// Definición de un bloque de carpetas
type DirectoryBlock struct {
	B_content [4]DirectoryContent // Contenido de la carpeta
}

type DirectoryContent struct {
	B_name  [12]byte // Nombre del archivo/carpeta
	B_inode int32    // Apuntador al inodo correspondiente
}

// Definición de un bloque de archivos (contenido)
type FileBlock struct {
	B_content [64]byte // Contenido del archivo (64 bytes)
}

// Definición de un bloque de apuntadores indirectos
type IndirectBlock struct {
	B_pointers [16]int32 // Apuntadores a otros bloques
}
