package comandos

import (
	"encoding/binary"
	"mimodulo/estructuras"
	"os"
	"regexp"
	"strings"
)

func Unmount(arre_coman []string) {
	Salid_comando += "===================UNMOUNT====================" + "\n"

	val_path := ""
	band_path := true
	//val_disco := ""
	posicion := 0
	val_id := ""
	regex := regexp.MustCompile(`^[a-zA-Z]+`)
	band_id := false
	band_error := false
	band_enc := false
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
			val_path = "/home/nataly/Documentos/Mia lab/Proyecto2/MIA_P2_202001570/Backend/Discos/MIA/P2/" + letters + ".dsk"
			splitted := regexp.MustCompile(`^[a-zA-Z](\d+)`)
			match := splitted.FindStringSubmatch(val_id)

			if len(match) > 1 {
				numbe := match[1][0:1]
				number = numbe
			}
			Salid_comando += number + "\n"

		default:
			Salid_comando += "Error: Parametro no valido." + "\n"
		}
	}

	if !band_error {
		if band_path {
			if band_id {
				if archivoExiste(val_path) {
					var empty [100]byte
					mbr := estructuras.Mbr{}
					disco, err := os.OpenFile(val_path, os.O_RDWR, 0660)

					if err != nil {
						Mens_error(err)
					}
					disco.Seek(0, 0)
					err = binary.Read(disco, binary.BigEndian, &mbr)

					if err != nil {
						Mens_error(err)
					}
					if mbr.Mbr_tamano != empty {
						for i := 0; i < 4; i++ {
							corre := string(mbr.Mbr_partition[i].Part_id[:])
							corre = strings.Trim(corre, "\x00")
							if corre == val_id {
								posicion = i
								band_enc = true
							}
						}

						if band_enc {
							emptyPartName := make([]byte, len(mbr.Mbr_partition[posicion].Part_id))
							copy(mbr.Mbr_partition[posicion].Part_id[:], emptyPartName)

							emptyPartstat := make([]byte, len(mbr.Mbr_partition[posicion].Part_status))
							copy(mbr.Mbr_partition[posicion].Part_status[:], emptyPartstat)
							copy(mbr.Mbr_partition[posicion].Part_status[:], "0")
							disco.Seek(0, os.SEEK_SET)
							err = binary.Write(disco, binary.BigEndian, &mbr)
							Salid_comando += "Particion desmontada exitosamente" + "\n"
						} else {
							Salid_comando += "Error: No existe la particion en el disco o la particion no es primaria" + "\n"
						}
						disco.Close()
					}
					disco.Close()
				} else {
					Salid_comando += "Error: El disco no existe." + "\n"
				}

			} else {
				Salid_comando += "Error: No se encontro el valor id" + "\n"
			}
		} else {
			Salid_comando += "Error: No se encontro el path" + "\n"
		}
	}
}
