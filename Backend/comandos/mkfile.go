package comandos

import (
	"encoding/binary"
	"fmt"
	"mimodulo/estructuras"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unsafe"
)

func Mkfile(arre_coman []string) {
	Salid_comando += "===================MKFILE====================" + "\n"

	val_path := ""
	val_r := ""
	val_size := ""
	val_cont := ""

	band_path := false
	band_r := false
	band_size := false
	band_cont := false
	band_enc := false

	band_error := false

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
			val_r = val_data
			band_r = true
		case strings.Contains(data, "size="):
			if band_size {
				Salid_comando += "Error: El parametro -size ya fue ingresado." + "\n"
				band_error = true
				break
			}
			val_size = val_data
			band_size = true
		case strings.Contains(data, "cont="):
			if band_cont {
				Salid_comando += "Error: El parametro -cont ya fue ingresado." + "\n"
				band_error = true
				break
			}
			val_cont = val_data
			band_cont = true
		default:
			Salid_comando += "Error: Parametro no valido." + "\n"
		}
	}

	if !band_error {
		if band_path {
			if val_r == "" {
				PathNewFile := ReturnValueWithoutMarks(val_path)
				mbr := estructuras.Mbr{}
				superBlock := estructuras.Super_bloque{}
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
					err = binary.Read(disco, binary.BigEndian, &superBlock)
					if superBlock.S_filesystem_type != empty {
						inodePrincipal := estructuras.Inodo{}
						s_inodo_start := string(superBlock.S_inode_start[:])
						s_inodo_start = strings.Trim(s_inodo_start, "\x00")
						inodo_start, err := strconv.Atoi(s_inodo_start)
						if err != nil {
							Mens_error(err)
						}
						disco.Seek(int64(int32(inodo_start)+int32(unsafe.Sizeof(estructuras.Inodo{}))), 0)
						err = binary.Read(disco, binary.LittleEndian, &inodePrincipal)
						if err != nil {
							Salid_comando += "Error: Al leer un Inode en el archivo " + val_rutadis + "\n"
						}
						arrayStrings := strings.Split(PathNewFile, "/")
						tmp := inodePrincipal
						cadena := ""
						podInode := 0
						for i := 1; i < len(arrayStrings); i++ {
							cadena = arrayStrings[i]
							siguienteInodo, existeValor := mkfile.ReturnDirExist(file, &superBlock, tmp, cadena)
							if existeValor {
								if i != len(arrayStrings)-1 {
									posicion := superBlock.S_inode_start + int32(siguienteInodo)*108
									podInode = siguienteInodo
									file.Seek(int64(posicion), 0)
									err := binary.Read(file, binary.LittleEndian, &tmp)
									if err != nil {
										fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
										return ""
									}
								} else {
									return "El archivo ya existe"
								}
							}
						}

					}
				}
			} else {
				Salid_comando += "Error: El parametro r no debe contener ningun valor" + "\n"
			}
		} else {
			Salid_comando += "Error: No se encontro el path" + "\n"
		}
	}
}

func ReturnValueWithoutMarks(value string) string {
	var tmpString string
	remplaceString := regexp.MustCompile("\"")
	tmpString = remplaceString.ReplaceAllString(value, "")
	tmpString = strings.TrimSpace(tmpString)
	return tmpString
}

func ReturnDirExist(file *os.File, superbloqu *estructuras.Super_bloque, Inode estructuras.Inodo, value string) (int, bool) {
	arrayDirectorios := strings.Split(value, "/")
	for i := 0; i < len(arrayDirectorios); i++ {
		siguienteInodo, existDir := ReturnExistValueInInode(file, superbloqu, &Inode, arrayDirectorios[i])
		if existDir {
			posicion := superbloqu.S_inode_start + int32(siguienteInodo)*108
			file.Seek(int64(posicion), 0)
			err := binary.Read(file, binary.LittleEndian, &Inode)
			if err != nil {
				Salid_comando += "Error: Al leer un Inode en el archivo " + "\n"
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
	i_types, err := strconv.Atoi(i_type)
	if err != nil {
		Mens_error(err)
	}
	if i_types == 1 {
		return -1, false
	}
	for i := 0; i < 16; i++ {
		block_i := string(inode.I_block[i])
		block_i = strings.Trim(block_i, "\x00")
		if block_i != "$" {
			dirBloc := estructuras.Bloque_carpeta{}
			block_st := string(superBlock.S_block_start[:])
			block_st = strings.Trim(block_st, "\x00")
			block_sta, err := strconv.Atoi(block_st)
			if err != nil {
				Mens_error(err)
			}
			block_i_int, err := strconv.Atoi(block_i)
			if err != nil {
				Mens_error(err)
			}
			posicion := block_sta + block_i_int*64
			file.Seek(int64(posicion), 0)
			errs := binary.Read(file, binary.LittleEndian, &dirBloc)
			if errs != nil {
				Salid_comando += "Error: Al leer un bloque de carpeta en el archivo " + "\n"
				return -1, false
			}
			nextInode, existDir := ReturnExistNameInBlock(&dirBloc, nameValue)
			if existDir {
				return nextInode, true
			}
		}
	}
	return -1, false
}

func ReturnExistNameInBlock(bloc *estructuras.Bloque_carpeta, nameValue string) (int, bool) {
	// El bloque contiene 4 elementos content que tienen nombre y el inodo referencia
	// vamos a validar si existe el nombre
	tmp := ReturnValueArr12BytesOFString(bloc.B_Content[2].B_name[:])
	if tmp == nameValue {
		return int(bloc.B_Content[2].B_inodp), true
	}
	tmp = ReturnValueArr12BytesOFString(bloc.B_Content[3].B_name[:])
	if tmp == nameValue {
		return int(bloc.B_Content[3].B_inodp), true
	}
	return -1, false
}
