package comandos

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"mimodulo/estructuras"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func Mkgrp(arre_coman []string) {
	Salid_comando += "===================MKGRP====================" + "\n"

	val_name := ""

	band_name := false

	band_error := false
	band_enc := false

	for i := 1; i < len(arre_coman); i++ {
		aux_data := strings.SplitAfter(arre_coman[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		switch {
		case strings.Contains(data, "name="):
			if band_name {
				Salid_comando += "Error: El parametro -name ya fue ingresado." + "\n"
				band_error = true
				break
			}
			val_name = val_data
			band_name = true
		default:
			Salid_comando += "Error: Parametro no valido." + "\n"
		}
	}

	if !band_error {
		if band_name {
			fmt.Println(iniciose)
			if iniciose == "root" {
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
						tmpString := ""
						blockFile := estructuras.Bloque_archivo{}
						for i := 0; i < 15; i++ {
							block_i := string(inodeA.I_block[i])
							block_i = strings.Trim(block_i, "\x00")
							if block_i != "$" {
								blockFile := estructuras.Bloque_archivo{}
								s_block_start := string(sb.S_block_start[:])
								s_block_start = strings.Trim(s_block_start, "\x00")
								block_start, err := strconv.Atoi(s_block_start)
								if err != nil {
									Mens_error(err)
								}
								I_block := string(inodeA.I_block[i])
								I_block = strings.Trim(I_block, "\x00")
								block_i, err := strconv.Atoi(I_block)
								if err != nil {
									Mens_error(err)
								}
								pos := int32(block_start) + int32(block_i)*int32(unsafe.Sizeof(estructuras.Bloque_carpeta{}))
								disco.Seek(int64(pos), 0)
								err = binary.Read(disco, binary.LittleEndian, &blockFile)
								if err != nil {
									Salid_comando += "Error al leer un bloque_archivo en el archivo"
								}
								tmp1 := string(blockFile.B_content[:])
								tmpString += tmp1
							}
						}
						res1 := strings.Split(tmpString, "\n")
						finalString := ""
						numGrupExist := 0
						for i := 0; i < len(res1); i++ {
							res2 := strings.Split(res1[i], ",")
							if len(res2) > 2 && res2[1] == "G" {
								numGrup, err := strconv.Atoi(res2[0])
								if err != nil {
									Mens_error(err)
								}
								numGrupExist = numGrup
								finalString += res1[i] + "\n"
							} else if len(res2) > 2 {
								finalString += res1[i] + "\n"
							}
						}
						numGrupExist++
						stringIsert := "," + strconv.Itoa(numGrupExist) + ",G," + val_name
						finalString += stringIsert
						if len(finalString) > 64 {
							CantM := len(finalString) / 64
							numBlocks := math.Floor(float64(CantM))
							for int(numBlocks*64) < len(finalString) {
								numBlocks++
							}
							arrBlocksFree := make([]int, int64(numBlocks))
							s_blockm_start := string(sb.S_bm_block_start[:])
							s_blockm_start = strings.Trim(s_blockm_start, "\x00")
							blockm_start, err := strconv.Atoi(s_blockm_start)
							if err != nil {
								Mens_error(err)
							}
							disco.Seek(int64(blockm_start+64), 0)
							i, z := 0, 0
							for i < int(numBlocks) {
								var statusByte byte
								err = binary.Read(disco, binary.LittleEndian, &statusByte)
								if err != nil {
									log.Fatal(err)
								}
								if statusByte == '0' {
									i++
									arrBlocksFree = append(arrBlocksFree, z)
								}
								z++
								s_block_mstart := string(sb.S_bm_block_start[:])
								s_block_mstart = strings.Trim(s_block_mstart, "\x00")
								block_mstart, err := strconv.Atoi(s_block_mstart)
								if err != nil {
									Mens_error(err)
								}
								s_block_start := string(sb.S_block_start[:])
								s_block_start = strings.Trim(s_block_start, "\x00")
								block_start, err := strconv.Atoi(s_block_start)
								if err != nil {
									Mens_error(err)
								}
								if (z + int(block_mstart)) >= int(block_start) {
									break
								}
							}
							//reescribimos el inodo archivo
							emptyPartstat := make([]byte, len(inodeA.I_size))
							copy(inodeA.I_size[:], emptyPartstat)
							lenBytes := make([]byte, 4)
							binary.LittleEndian.PutUint32(lenBytes, uint32(len(finalString)))
							copy(inodeA.I_size[:], lenBytes)
							retura := ReturnDate8Bytes()
							emptyPartstast := make([]byte, len(inodeA.I_mtime))
							copy(inodeA.I_mtime[:], emptyPartstast)
							copy(inodeA.I_mtime[:], retura[:])

						} else {
							emptli := make([]byte, len(blockFile.B_content))
							copy(blockFile.B_content[:], emptli)
							respu := ReturnValueArr64Bytes(finalString)
							copy(blockFile.B_content[:], respu[:])
							s_block_s := string(sb.S_block_start[:])
							s_block_s = strings.Trim(s_block_s, "\x00")
							block_s, err := strconv.Atoi(s_block_s)
							if err != nil {
								Mens_error(err)
							}
							disco.Seek(int64(block_s), 0)
							err = binary.Write(disco, binary.LittleEndian, &blockFile)
							if err != nil {
								log.Println(err)
							}
							emptyPartstatw := make([]byte, len(inodeA.I_size))
							copy(inodeA.I_size[:], emptyPartstatw)
							lenBytes := make([]byte, 4)
							binary.LittleEndian.PutUint32(lenBytes, uint32(len(finalString)))
							copy(inodeA.I_size[:], lenBytes)
							retura := ReturnDate8Bytes()
							emptyPartstast := make([]byte, len(inodeA.I_mtime))
							copy(inodeA.I_mtime[:], emptyPartstast)
							copy(inodeA.I_mtime[:], retura[:])
							s_inod_start := string(sb.S_inode_start[:])
							s_inod_start = strings.Trim(s_inod_start, "\x00")
							inod_start, err := strconv.Atoi(s_inod_start)
							if err != nil {
								Mens_error(err)
							}
							disco.Seek(int64(inod_start)+int64(unsafe.Sizeof(estructuras.Inodo{})), 0)
							err = binary.Write(disco, binary.LittleEndian, &inodeA)
							if err != nil {
								log.Println(err)
							}
						}
						Salid_comando += "Grupo creado con exito" + "\n"
					}
				}
			} else {
				Salid_comando += "Error: No se puede crear un grupo con un usuario que no sea root" + "\n"
			}
		} else {
			Salid_comando += "Error: No se encontro el valor name" + "\n"
		}

	}
}
func ReturnDate8Bytes() [8]byte {
	t := string(time.Now().Format("02012006"))
	tmpT := []byte(t)
	return [8]byte(tmpT)
}

func ReturnValueArr64Bytes(value string) [64]byte {
	var tmp [64]byte
	for i := 0; i < 64; i++ {
		if i >= len(value) {
			break
		}
		tmp[i] = value[i]
	}
	return tmp
}
