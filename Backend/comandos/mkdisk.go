package comandos

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"io"
	"math/rand"
	"mimodulo/estructuras"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func Mens_error(err error) {
	Salid_comando += errors.New("Error: "+err.Error()).Error() + "\n"
}

var valorpath string = "/home/nataly/Documentos/Mia lab/Proyecto2/MIA_P2_202001570/Backend/Discos/MIA/P2"
var path string = ""
var alfabeto []string = []string{
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J",
	"K", "L", "M", "N", "Ñ", "O", "P", "Q", "R", "S",
	"T", "U", "V", "W", "X", "Y", "Z",
}
var posicion = 0
var Salid_comando string = ""

func Mkdisk(arre_coman []string) {
	Salid_comando += "=========================MKDISK==========================" + "\n"

	val_size := 0
	val_fit := "ff"
	val_unit := "m"
	val_path := "/home/nataly/Documentos/Mia lab/Proyecto2/MIA_P2_202001570/Backend/Discos/MIA/P2"

	band_size := false
	band_fit := false
	band_unit := false
	band_path := true
	band_error := false

	for i := 1; i < len(arre_coman); i++ {
		aux_data := strings.SplitAfter(arre_coman[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		switch {
		case strings.Contains(data, "size="):
			if band_size {
				Salid_comando += "Error: El parametro -size ya fue ingresado." + "\n"
				band_error = true
				break
			}

			band_size = true

			aux_size, err := strconv.Atoi(val_data)
			val_size = aux_size

			if err != nil {
				Mens_error(err)
			}

			if val_size < 0 {
				band_error = true
				Salid_comando += "Error: El parametro -size es negativo." + "\n"
				break
			}
		case strings.Contains(data, "fit="):
			if band_fit {
				Salid_comando += "Error: El parametro -fit ya fue ingresado." + "\n"
				band_error = true
				break
			}

			if val_fit == "bf" {
				band_fit = true
				val_fit = "b"
			} else if val_fit == "ff" {
				band_fit = true
				val_fit = "f"
			} else if val_fit == "wf" {
				band_fit = true
				val_fit = "w"
			} else {
				Salid_comando += "Error: El Valor del parametro -fit no es valido." + "\n"
				band_error = true
				break
			}
		case strings.Contains(data, "unit="):
			if band_unit {
				Salid_comando += "Error: El parametro -unit ya fue ingresado." + "\n"
				band_error = true
				break
			}

			val_unit = strings.Replace(val_data, "\"", "", 2)
			val_unit = strings.ToLower(val_unit)

			if val_unit == "k" || val_unit == "m" {
				band_unit = true
			} else {
				Salid_comando += "Error: El Valor del parametro -unit no es valido." + "\n"
				band_error = true
				break
			}
		default:
			Salid_comando += "Error: Parametro no valido." + "\n"
		}
	}

	if !band_error {
		if band_path {
			if band_size {
				total_size := 1024
				master_boot_record := estructuras.Mbr{}

				Crear_disco(val_path)

				fecha := time.Now()
				str_fecha := fecha.Format("02/01/2006 15:04:05")

				copy(master_boot_record.Mbr_fecha_creacion[:], str_fecha)

				rand.Seed(time.Now().UnixNano())
				min := 0
				max := 100
				num_random := rand.Intn(max-min+1) + min

				copy(master_boot_record.Mbr_dsk_signature[:], strconv.Itoa(int(num_random)))

				if band_fit {
					copy(master_boot_record.Dsk_fit[:], val_fit)
				} else {
					copy(master_boot_record.Dsk_fit[:], "f")
				}

				if band_unit {
					if val_unit == "m" {
						copy(master_boot_record.Mbr_tamano[:], strconv.Itoa(int(val_size*1024*1024)))
						total_size = val_size * 1024
					} else {
						copy(master_boot_record.Mbr_tamano[:], strconv.Itoa(int(val_size*1024)))
						total_size = val_size
					}
				} else {
					copy(master_boot_record.Mbr_tamano[:], strconv.Itoa(int(val_size*1024*1024)))
					total_size = val_size * 1024
				}

				for i := 0; i < 4; i++ {
					copy(master_boot_record.Mbr_partition[i].Part_status[:], "0")
					copy(master_boot_record.Mbr_partition[i].Part_type[:], "0")
					copy(master_boot_record.Mbr_partition[i].Part_fit[:], "0")
					copy(master_boot_record.Mbr_partition[i].Part_start[:], "-1")
					copy(master_boot_record.Mbr_partition[i].Part_size[:], "0")
					copy(master_boot_record.Mbr_partition[i].Part_name[:], "")
				}

				str_total_size := strconv.Itoa(total_size)

				cmd := exec.Command("/bin/sh", "-c", "dd if=/dev/zero of=\""+path+"\" bs=1024 count="+str_total_size)
				cmd.Dir = "/"
				_, err := cmd.Output()

				if err != nil {
					Mens_error(err)
				}

				disco, err := os.OpenFile(path, os.O_RDWR, 0660)

				if err != nil {
					Mens_error(err)
				}

				disco.Seek(0, os.SEEK_SET)

				err = binary.Write(disco, binary.BigEndian, master_boot_record)
				if err != nil {
					Mens_error(err)
				}
				disco.Close()
				Salid_comando += "El disco fue creado exitosamente." + "\n"
			} else {
				Salid_comando += "Error: No se encontro el valor size" + "\n"
			}
		} else {
			Salid_comando += "Error: No se encontro el path" + "\n"
		}
	}
}

func Crear_disco(ruta string) {
	nome := NombreDisk()
	if nome == "error" {
		Salid_comando += "Error: No se encontro el nombre del disco" + "\n"
	} else {
		ruta = ruta + "/" + nome
	}
	path = ruta
	aux, err := filepath.Abs(ruta)

	if err != nil {
		Mens_error(err)
	}

	cmd1 := exec.Command("/bin/sh", "-c", "sudo mkdir -p '"+filepath.Dir(aux)+"'")
	cmd1.Dir = "/"
	_, err1 := cmd1.Output()

	if err1 != nil {
		Mens_error(err)
	}

	cmd2 := exec.Command("/bin/sh", "-c", "sudo chmod -R 777 '"+filepath.Dir(aux)+"'")
	cmd2.Dir = "/"
	_, err2 := cmd2.Output()

	if err2 != nil {
		Mens_error(err)
	}

	if _, err := os.Stat(filepath.Dir(aux)); errors.Is(err, os.ErrNotExist) {
		if err != nil {
			Salid_comando += "[FAILURE] No se pudo crear el disco." + "\n"
		}
	}
}
func Struct_a_bytes(p interface{}) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)

	if err != nil && err != io.EOF {
		Mens_error(err)
	}

	return buf.Bytes()
}

func NombreDisk() string {
	nombre := ""
	if posicion < 0 || posicion >= len(alfabeto) {
		nombre = "error"
	} else {
		nombre = alfabeto[posicion]
		nombre = nombre + ".dsk"
		posicion++
	}
	return nombre
}
