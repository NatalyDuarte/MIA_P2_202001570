package comandos

import (
	"encoding/binary"
	"io"
	"mimodulo/estructuras"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func Mkfs(arre_coman []string) {
	Salid_comando += "====================MKFS====================" + "\n"

	val_id := ""
	band_id := false

	val_type := "full"
	band_type := false

	val_fs := "2fs"
	band_fs := false

	band_error := false

	regex := regexp.MustCompile(`^[a-zA-Z]+`)
	val_path := ""
	var number string

	for i := 1; i < len(arre_coman); i++ {
		aux_data := strings.SplitAfter(arre_coman[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		switch {
		case strings.Contains(data, "id="):
			if band_id {
				Salid_comando += "Error: El parametro -id ya fue ingresado." + "\n"
				band_error = true
				break
			}
			val_id = val_data
			band_id = true
			letters := regex.FindString(val_id)
			Salid_comando += "Disco: " + letters + ".dsk" + "\n"
			val_path = "/home/nataly/Documentos/Mia lab/Proyecto2/MIA_P2_202001570/Backend/Discos/MIA/P2/" + letters + ".dsk"
			splitted := regexp.MustCompile(`^[a-zA-Z](\d+)`)
			match := splitted.FindStringSubmatch(val_id)

			if len(match) > 1 {
				numbe := match[1][0:1]
				number = numbe
			}
			Salid_comando += "Particion: " + number + "\n"
		case strings.Contains(data, "type="):
			if band_type {
				Salid_comando += "Error: El parametro -type ya fue ingresado." + "\n"
				band_error = true
				break
			}
			val_type = val_data
			band_type = true
			if val_type != "full" {
				Salid_comando += "Error: El type debe ser full" + "\n"
				band_error = true
				break
			}
		case strings.Contains(data, "fs="):
			if band_fs {
				Salid_comando += "Error: El parametro -fs ya fue ingresado." + "\n"
				band_error = true
				break
			}
			if val_data == "2fs" {
				val_fs = "2fs"
				band_fs = true
			} else if val_data == "3fs" {
				val_fs = "3fs"
				band_fs = true
			} else {
				Salid_comando += "Error: El valor de fs no es el indicado" + "\n"
				band_error = true
				break
			}

		default:
			Salid_comando += "Error: Parametro no valido." + "\n"
		}
	}
	if !band_error {
		if band_id {
			if val_fs == "2fs" {
				Salid_comando += "Formateando con 2fs" + "\n"
				Formatearext2(val_path, val_id)

			} else if val_fs == "3fs" {
				Salid_comando += "Formateando con 3fs" + "\n"
			}

		} else {
			Salid_comando += "Error: No se ingreso el -id y es obligatorio" + "\n"
		}

	}

}

func Formatearext2(val_path string, val_id string) {
	band_crea := false
	var empty [100]byte
	mbr := estructuras.Mbr{}
	disco, err := os.OpenFile(val_path, os.O_RDWR, 0660)

	if err != nil {
		Mens_error(err)
	}
	defer disco.Close()
	disco.Seek(0, 0)
	err = binary.Read(disco, binary.BigEndian, &mbr)
	if mbr.Mbr_tamano != empty {
		for i := 0; i < 4; i++ {
			id := string(mbr.Mbr_partition[i].Part_id[:])
			id = strings.Trim(id, "\x00")
			if id == val_id {
				posicion = i
				band_crea = true
			}
		}

	}
	if band_crea {
		TamanioParti := string(mbr.Mbr_partition[posicion].Part_size[:])
		TamanioParti = strings.Trim(TamanioParti, "\x00")
		tamanioparti, err := strconv.Atoi(TamanioParti)
		if err != nil {
			Salid_comando += "Error: No se pudo convertir el tamaño de la particion" + "\n"
		}
		InicioParti := string(mbr.Mbr_partition[posicion].Part_start[:])
		InicioParti = strings.Trim(InicioParti, "\x00")
		inicio, err := strconv.Atoi(InicioParti)
		if err != nil {
			Salid_comando += "Error: No se pudo convertir el inicio de la particion" + "\n"
		}
		sb_empty := estructuras.Super_bloque{}

		in_empty := estructuras.Inodo{}

		ba_empty := estructuras.Bloque_archivo{}

		n := (int64(tamanioparti) - int64(unsafe.Sizeof(sb_empty))) / (4 + int64(unsafe.Sizeof(in_empty)) + 3*int64(unsafe.Sizeof(ba_empty)))
		num_estructuras := n
		num_bloques := 3 * num_estructuras

		Super_bloque := estructuras.Super_bloque{}

		fecha := time.Now()
		str_fecha := fecha.Format("02/01/2006 15:04:05")

		copy(Super_bloque.S_filesystem_type[:], "2")
		copy(Super_bloque.S_inodes_count[:], strconv.Itoa(int(num_estructuras)))
		copy(Super_bloque.S_blocks_count[:], strconv.Itoa(int(num_bloques)))
		copy(Super_bloque.S_free_blocks_count[:], strconv.Itoa(int(num_bloques-2)))
		copy(Super_bloque.S_free_inodes_count[:], strconv.Itoa(int(num_estructuras-2)))
		copy(Super_bloque.S_mtime[:], str_fecha)
		copy(Super_bloque.S_mnt_count[:], "0")
		copy(Super_bloque.S_magic[:], "0xEF53")
		copy(Super_bloque.S_inode_size[:], strconv.Itoa(int(unsafe.Sizeof(in_empty))))
		copy(Super_bloque.S_block_size[:], strconv.Itoa(int(unsafe.Sizeof(ba_empty))))
		copy(Super_bloque.S_firts_ino[:], "2")
		copy(Super_bloque.S_first_blo[:], "2")
		copy(Super_bloque.S_bm_inode_start[:], strconv.Itoa(inicio+int(unsafe.Sizeof(sb_empty))))
		copy(Super_bloque.S_bm_block_start[:], strconv.Itoa(inicio+int(unsafe.Sizeof(sb_empty))+int(num_estructuras)))
		copy(Super_bloque.S_inode_start[:], strconv.Itoa(inicio+int(unsafe.Sizeof(sb_empty))+int(num_estructuras)+int(num_bloques)))
		copy(Super_bloque.S_block_start[:], strconv.Itoa(inicio+int(unsafe.Sizeof(sb_empty))+int(num_estructuras)+int(num_bloques)+int(unsafe.Sizeof(in_empty))+int(num_estructuras)))

		inodo := estructuras.Inodo{}
		bloque := estructuras.Bloque_carpeta{}

		// Libre
		buffer := "0"
		// Usado o archivo
		buffer2 := "1"
		// Carpeta
		buffer3 := "2"

		if err == nil {
			// Super Bloque
			disco.Seek(int64(inicio), io.SeekStart)
			err = binary.Write(disco, binary.BigEndian, &Super_bloque)

			s_bm_inode_start := string(Super_bloque.S_bm_inode_start[:])
			s_bm_inode_start = strings.Trim(s_bm_inode_start, "\x00")
			i_bm_inode_start, _ := strconv.Atoi(s_bm_inode_start)

			for i := 0; i < int(num_estructuras); i++ {
				disco.Seek(int64(i_bm_inode_start+i), io.SeekStart)
				disco.Write([]byte(buffer))
			}

			disco.Seek(int64(i_bm_inode_start), io.SeekStart)
			disco.Write([]byte(buffer2))
			disco.Write([]byte(buffer2))

			s_bm_block_start := string(Super_bloque.S_bm_block_start[:])
			s_bm_block_start = strings.Trim(s_bm_block_start, "\x00")
			i_bm_block_start, _ := strconv.Atoi(s_bm_block_start)

			for i := 0; i < int(num_bloques); i++ {
				disco.Seek(int64(i_bm_block_start+i), io.SeekStart)
				disco.Write([]byte(buffer))
			}

			// Marcando el Bitmap de Inodos para la carpeta "/" y el archivo "users.txt"
			disco.Seek(int64(i_bm_block_start), io.SeekStart)
			disco.Write([]byte(buffer2))
			disco.Write([]byte(buffer3))

			/* Inodo para la carpeta "/" */
			copy(inodo.I_uid[:], "1")
			copy(inodo.I_gid[:], "1")
			copy(inodo.I_size[:], "0")
			copy(inodo.I_atime[:], str_fecha)
			copy(inodo.I_ctime[:], str_fecha)
			copy(inodo.I_mtime[:], str_fecha)
			copy(inodo.I_block[0:1], "0")

			// Nodos libres
			for i := 1; i < 15; i++ {
				copy(inodo.I_block[i:i+1], "$")
			}

			// 1 = archivo o 0 = Carpeta
			copy(inodo.I_type[:], "0")
			copy(inodo.I_perm[:], "664")

			s_inode_start := string(Super_bloque.S_inode_start[:])
			s_inode_start = strings.Trim(s_inode_start, "\x00")
			i_inode_start, _ := strconv.Atoi(s_inode_start)

			disco.Seek(int64(i_inode_start), io.SeekStart)
			err = binary.Write(disco, binary.BigEndian, &inodo)

			/* Bloque Para Carpeta "/" */
			copy(bloque.B_content[0].B_name[:], ".")
			copy(bloque.B_content[0].B_inodo[:], "0")

			// Directorio Padre
			copy(bloque.B_content[1].B_name[:], "..")
			copy(bloque.B_content[1].B_inodo[:], "0")

			// Nombre de la carpeta o archivo
			copy(bloque.B_content[2].B_name[:], "users.txt")
			copy(bloque.B_content[2].B_inodo[:], "1")
			copy(bloque.B_content[3].B_name[:], ".")
			copy(bloque.B_content[3].B_inodo[:], "-1")

			s_block_start := string(Super_bloque.S_block_start[:])
			s_block_start = strings.Trim(s_block_start, "\x00")
			i_block_start, _ := strconv.Atoi(s_block_start)

			disco.Seek(int64(i_block_start), io.SeekStart)
			err = binary.Write(disco, binary.BigEndian, &bloque)

			/* Inodo Para "users.txt" */
			copy(inodo.I_uid[:], "1")
			copy(inodo.I_gid[:], "1")
			copy(inodo.I_size[:], "29")
			copy(inodo.I_atime[:], str_fecha)
			copy(inodo.I_ctime[:], str_fecha)
			copy(inodo.I_mtime[:], str_fecha)
			copy(inodo.I_block[0:1], "1")

			// Nodos libres
			for i := 1; i < 15; i++ {
				copy(inodo.I_block[i:i+1], "$")
			}

			copy(inodo.I_type[:], "1")
			copy(inodo.I_perm[:], "755")

			disco.Seek(int64(i_inode_start+int(unsafe.Sizeof(in_empty))), io.SeekStart)
			err = binary.Write(disco, binary.BigEndian, &inodo)

			/* Bloque Para "users.txt" */
			archivo := estructuras.Bloque_archivo{}

			for i := 0; i < 100; i++ {
				copy(archivo.B_content[i:i+1], "0")
			}

			bc_empty := estructuras.Bloque_carpeta{}

			// GID, TIPO, GRUPO
			// UID, TIPO, GRUPO, USUARIO, CONTRASEÑA
			copy(archivo.B_content[:], "1,G,root\n1,U,root,root,123\n")
			disco.Seek(int64(i_block_start+int(unsafe.Sizeof(bc_empty))), io.SeekStart)
			err = binary.Write(disco, binary.BigEndian, &archivo)

			Salid_comando += "El Disco se formateo en el sistema EXT2 con exitosamente" + "\n"
			disco.Close()
		}
	}
}

func ReturnValueArr(value string) [12]byte {
	var tmp [12]byte
	for i := 0; i < 12; i++ {
		if i >= len(value) {
			break
		}
		tmp[i] = value[i]
	}
	return tmp
}
func ReturnValueArr64(value string) [64]byte {
	var tmp [64]byte
	for i := 0; i < 64; i++ {
		if i >= len(value) {
			break
		}
		tmp[i] = value[i]
	}
	return tmp
}
