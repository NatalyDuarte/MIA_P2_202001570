package comandos

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Mens_errorrm(err error) {
	Salid_comando += errors.New("Error: "+err.Error()).Error() + "\n"
}

var valorpathrm string = "/home/nataly/Documentos/Mia lab/Proyecto2/MIA_P2_202001570/Backend/Discos/MIA/P2"
var pathrm string = ""

func Rmdisk(arre_coman []string) {
	Salid_comando += "======================RMDISK=======================" + "\n"

	val_path := "/home/nataly/Documentos/Mia lab/Proyecto2/MIA_P2_202001570/Backend/Discos/MIA/P2"
	band_path := true

	val_driveletter := ""
	band_driveletter := false
	band_error := false

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
			Salid_comando += val_data + "\n"
			val_driveletter = val_data
			band_driveletter = true
		default:
			Salid_comando += "Error: Parametro no valido." + "\n"
		}
	}

	if !band_error {
		if band_path {
			if band_driveletter {
				val_path = val_path + "/" + val_driveletter + ".dsk"
				Salid_comando += val_path + "\n"
				_, e := os.Stat(val_path)
				if e != nil {
					if os.IsNotExist(e) {
						Salid_comando += "Error: No existe el disco que desea eliminar." + "\n"
						band_path = false
					}
				} else {
					fmt.Print("Desea eliminar el disco [S/N]?: ")
					var opcion string
					fmt.Scanln(&opcion)
					if opcion == "s" || opcion == "S" {
						cmd := exec.Command("/bin/sh", "-c", "rm \""+val_path+"\"")
						cmd.Dir = "/"
						_, err := cmd.Output()

						if err != nil {
							Mens_errorrm(err)
						} else {
							Salid_comando += "El Disco fue eliminado exitosamente" + "\n"
						}
						band_path = false
					} else if opcion == "n" || opcion == "N" {
						Salid_comando += "EL disco no fue eliminado" + "\n"
						band_path = false
					} else {
						Salid_comando += "Error: Opcion no valida intentalo de nuevo." + "\n"
					}
				}
			} else {
				Salid_comando += "Error: No se encontro el valor driveletter" + "\n"
			}
		} else {
			Salid_comando += "Error: No se encontro el path" + "\n"
		}
	}
}
