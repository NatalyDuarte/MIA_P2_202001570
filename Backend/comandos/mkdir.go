package comandos

import (
	"encoding/binary"
	"fmt"
	"mimodulo/estructuras"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

func Mkdir(arre_coman []string) {
	Salid_comando += "===================MKDIR====================" + "\n"

	val_path := ""
	band_path := false

	band_r := false

	band_error := false
	band_enc := false

	for i := 1; i < len(arre_coman); i++ {
		aux_data := strings.SplitAfter(arre_coman[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		switch {
		case strings.Contains(data, "path="):
			if band_path {
				Salid_comando += "Error: El parametro -path ya fue ingresado." + "\n"
				band_error = true
				break
			}
			val_path = val_data
			band_path = true
		case strings.Contains(data, "r"):
			if band_r {
				Salid_comando += "Error: El parametro -r ya fue ingresado." + "\n"
				band_error = true
				break
			}
			band_r = true

		default:
			Salid_comando += "Error: Parametro no valido." + "\n"
		}
	}

	if !band_error {
		if band_path {
			mbr := estructuras.Mbr{}
			sb := estructuras.Super_bloque{}
			disco, err := os.OpenFile(discos, os.O_RDWR, 0660)
			var empty [100]byte
			if err != nil {
				Mens_error(err)
			}
			defer func() {
				disco.Close()
			}()
			disco.Seek(0, 0)
			err = binary.Read(disco, binary.BigEndian, &mbr)
			if mbr.Mbr_tamano != empty {

			}
			numero_parti := 0
			for i := 0; i < 4; i++ {
				s_part_id := string(mbr.Mbr_partition[i].Part_id[:])
				s_part_id = strings.Trim(s_part_id, "\x00")

				if s_part_id == particion {
					numero_parti = i
					band_enc = true
					break
				}

			}
			if band_enc {
				s_part_startas := string(mbr.Mbr_partition[numero_parti].Part_start[:])
				s_part_startas = strings.Trim(s_part_startas, "\x00")
				part_starta, err := strconv.Atoi(s_part_startas)
				if err != nil {
					Mens_error(err)
				}
				disco.Seek(int64(part_starta), 0)
				err = binary.Read(disco, binary.BigEndian, &sb)
				if sb.S_filesystem_type != empty {
					inodeA := estructuras.Inodo{}
					s_inodo_start := string(sb.S_inode_start[:])
					s_inodo_start = strings.Trim(s_inodo_start, "\x00")
					inodo_start, err := strconv.Atoi(s_inodo_start)
					if err != nil {
						Mens_error(err)
					}
					disco.Seek(int64(int32(inodo_start)+int32(unsafe.Sizeof(estructuras.Inodo{}))), 0)
					err = binary.Read(disco, binary.LittleEndian, &inodeA)
					if err != nil {
						Salid_comando += "Error: Al leer un Inode en el archivo " + val_rutadis + "\n"
					}
					arrayStrings := strings.Split(val_path, "/")
					tmp := inodeA
					cadena := ""
					podInode := 0
					for i := 1; i < len(arrayStrings); i++ {
						cadena = arrayStrings[i]
						siguienteInodo, existeValor := ReturnDirExist(disco, &sb, tmp, cadena)
						if existeValor {
							s_inodo_start := string(sb.S_inode_start[:])
							s_inodo_start = strings.Trim(s_inodo_start, "\x00")
							inodo_start, err := strconv.Atoi(s_inodo_start)
							if err != nil {
								Mens_error(err)
							}
							posicion := int32(inodo_start) + int32(siguienteInodo)*108
							podInode = siguienteInodo
							disco.Seek(int64(posicion), 0)
							errs := binary.Read(disco, binary.LittleEndian, &tmp)
							if errs != nil {
								Salid_comando += "Error: Al leer el archivo " + "\n"

							}
							fmt.Println(podInode)
						}
					}
				}

			}
		} else {
			Salid_comando += "Error: El path es obligatorio" + "\n"
		}
	}
}

func ReturnDirExist(file *os.File, superbloqu *estructuras.Super_bloque, Inode estructuras.Inodo, value string) (int, bool) {
	arrayDirectorios := strings.Split(value, "/")
	for i := 0; i < len(arrayDirectorios); i++ {
		siguienteInodo, existDir := ReturnExistValueInInode(file, superbloqu, &Inode, arrayDirectorios[i])
		if existDir {
			s_inodo_start := string(superbloqu.S_inode_start[:])
			s_inodo_start = strings.Trim(s_inodo_start, "\x00")
			inodo_start, err := strconv.Atoi(s_inodo_start)
			if err != nil {
				Mens_error(err)
			}
			posicion := int32(inodo_start) + int32(siguienteInodo)*108
			file.Seek(int64(posicion), 0)
			errs := binary.Read(file, binary.LittleEndian, &Inode)
			if errs != nil {
				Salid_comando += "Error: AL leer el archivo " + "\n"
				return -1, false
			}
			if i == len(arrayDirectorios)-1 {
				return siguienteInodo, true
			}
		}
	}
	return -1, false
}

func ReturnExistValueInInode(file *os.File, superBlock *estructuras.Super_bloque, inode *estructuras.Inodo, nameValue string) (int, bool) {
	i_type := string(inode.I_type[:])
	i_type = strings.Trim(i_type, "\x00")
	inodo_type, err := strconv.Atoi(i_type)
	if err != nil {
		Mens_error(err)
	}
	if inodo_type == 1 {
		return -1, false
	}
	for i := 0; i < 16; i++ {
		block_i := string(inode.I_block[i])
		block_i = strings.Trim(block_i, "\x00")
		if block_i != "$" {
			dirBloc := estructuras.Bloque_carpeta{}
			s_block_start := string(superBlock.S_block_start[:])
			s_block_start = strings.Trim(s_block_start, "\x00")
			block_start, errS := strconv.Atoi(s_block_start)
			if errS != nil {
				Mens_error(err)
			}
			I_block := string(inode.I_block[i])
			I_block = strings.Trim(I_block, "\x00")
			block_is, err := strconv.Atoi(I_block)
			if err != nil {
				Mens_error(err)
			}

			posicion := block_start + block_is*64
			file.Seek(int64(posicion), 0)
			erra := binary.Read(file, binary.LittleEndian, &dirBloc)
			if erra != nil {
				Salid_comando += "Error: Al leer el archivo " + "\n"
				return -1, false
			}
			/*nextInode, existDir := ReturnExistNameInBlock(&dirBloc, nameValue)
			if existDir {
				return nextInode, true
			}*/
		}
	}
	return -1, false
}
