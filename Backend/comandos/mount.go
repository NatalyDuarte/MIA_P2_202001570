package comandos

import (
	"encoding/binary"
	"mimodulo/estructuras"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var valorpathmo string = "/home/nataly/Documentos/Mia lab/Proyecto2/MIA_P2_202001570/Backend/Discos/MIA/P2/"

func Mount(arre_coman []string) {
	Salid_comando += "===================MOUNT====================" + "\n"

	val_path := ""
	band_path := true
	expresionRegular := regexp.MustCompile(`\d+`)
	numerosComoString := ""

	val_driveletter := ""
	val_name := ""
	termina := "70"
	band_driveletter := false
	band_name := false
	band_error := false
	band_enc := false

	for i := 1; i < len(arre_coman); i++ {
		aux_data := strings.SplitAfter(arre_coman[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		switch {
		case strings.Contains(data, "driveletter="):
			if band_driveletter {
				Salid_comando += "Error: El parametro -driveletter ya fue ingresado." + "\n"
				band_error = true
				break
			}
			val_driveletter = val_data
			band_driveletter = true
		case strings.Contains(data, "name="):
			if band_name {
				Salid_comando += "Error: El parametro -name ya fue ingresado." + "\n"
				band_error = true
				break
			}
			val_name = val_data
			band_name = true
			coincidencias := expresionRegular.FindAllString(val_name, -1)
			numerosComoString = strings.Join(coincidencias, ",")
		default:
			Salid_comando += "Error: Parametro no valido." + "\n"
		}
	}

	if !band_error {
		if band_path {
			if band_driveletter {
				if band_name {
					Salid_comando += val_driveletter + numerosComoString + termina + "\n"
					val_path = valorpathmo + val_driveletter + ".dsk"
					if archivoExiste(val_path) {
						var empty [100]byte
						mbr := estructuras.Mbr{}
						disco, err := os.OpenFile(val_path, os.O_RDWR, 0660)

						if err != nil {
							Mens_error(err)
						}
						disco.Seek(0, 0)
						err = binary.Read(disco, binary.BigEndian, &mbr)
						posicion := 0

						if err != nil {
							Mens_error(err)
						}
						if mbr.Mbr_tamano != empty {
							for i := 0; i < 4; i++ {
								name := string(mbr.Mbr_partition[i].Part_name[:])
								name = strings.Trim(name, "\x00")
								types := string(mbr.Mbr_partition[i].Part_type[:])
								types = strings.Trim(types, "\x00")
								status := string(mbr.Mbr_partition[i].Part_status[:])
								status = strings.Trim(status, "\x00")
								if name == val_name {
									if types == "p" {
										if status != "1" {
											band_enc = true
											posicion = i
											break
										} else {
											Salid_comando += "Error: No se puede montar la particion ya que se encuentra montada" + "\n"
										}

									}

								}
							}

							if band_enc {
								Salid_comando += "Particion encontrada" + "\n"
								copy(mbr.Mbr_partition[posicion].Part_status[:], "1")

								copy(mbr.Mbr_partition[posicion].Part_id[:], []byte(val_driveletter+numerosComoString+termina))

								disco.Seek(0, os.SEEK_SET)
								err = binary.Write(disco, binary.BigEndian, &mbr)

								if err != nil {
									Mens_error(err)
								}

								Salid_comando += "Particion montada exitosamente" + "\n"
								imprimir(val_path)

							} else {
								Salid_comando += "Error: No existe la particion en el disco o la particion no es primaria o ya esta montada la particion" + "\n"
							}
							disco.Close()
						}
					} else {
						Salid_comando += "Error: El archivo no existe." + "\n"
					}

				} else {
					Salid_comando += "Error: No se encontro el valor de name" + "\n"
				}
			} else {
				Salid_comando += "Error: No se encontro el valor driveletter" + "\n"
			}
		} else {
			Salid_comando += "Error: No se encontro el path" + "\n"
		}
	}
}
func archivoExiste(ruta string) bool {
	_, err := os.Stat(ruta)
	return !os.IsNotExist(err)
}
func imprimir(ruta string) {
	var empty [100]byte
	mbr := estructuras.Mbr{}

	disco, err := os.OpenFile(ruta, os.O_RDWR, 0660)

	if err != nil {
		Mens_error(err)
	}
	defer func() {
		disco.Close()
	}()

	disco.Seek(0, 0)
	err = binary.Read(disco, binary.BigEndian, &mbr)

	if err != nil {
		Mens_error(err)
	}

	if mbr.Mbr_tamano != empty {
		for i := 0; i < 4; i++ {
			s_part_status := string(mbr.Mbr_partition[i].Part_status[:])
			s_part_status = strings.Trim(s_part_status, "\x00")

			if s_part_status == "1" {
				Salid_comando += "--------------------Particion: " + strconv.Itoa(i) + "--------------------" + "\n"
				part_sta := string(mbr.Mbr_partition[i].Part_status[:])
				part_sta = strings.Trim(part_sta, "\x00")
				Salid_comando += "Part_Status: " + part_sta + " \n"
				part_ty := string(mbr.Mbr_partition[i].Part_type[:])
				part_ty = strings.Trim(part_ty, "\x00")
				Salid_comando += "Part_Type: " + part_ty + " \n"
				part_fi := string(mbr.Mbr_partition[i].Part_fit[:])
				part_fi = strings.Trim(part_fi, "\x00")
				Salid_comando += "Part_Fit: " + part_fi + " \n"
				part_star := string(mbr.Mbr_partition[i].Part_start[:])
				part_star = strings.Trim(part_star, "\x00")
				Salid_comando += "Part_start: " + part_star + " \n"
				part_si := string(mbr.Mbr_partition[i].Part_size[:])
				part_si = strings.Trim(part_si, "\x00")
				Salid_comando += "Part_size: " + part_si + " \n"
				part_nam := string(mbr.Mbr_partition[i].Part_name[:])
				part_nam = strings.Trim(part_nam, "\x00")
				Salid_comando += "Part_Name: " + part_nam + " \n"
				part_corr := string(mbr.Mbr_partition[i].Part_correlative[:])
				part_corr = strings.Trim(part_corr, "\x00")
				Salid_comando += "Part_correlativo: " + part_corr + " \n"
				part_id := string(mbr.Mbr_partition[i].Part_id[:])
				part_id = strings.Trim(part_id, "\x00")
				Salid_comando += "Part_id: " + part_id + " \n"
			}
		}
	}
}
