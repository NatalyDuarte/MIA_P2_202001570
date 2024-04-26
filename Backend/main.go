package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mimodulo/comandos"
	"mimodulo/estructuras"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/rs/cors"
)

var Salida_comando string = ""
var Salida_parti string = ""
var GraphDot string = ""

type FileList struct {
	Files []string `json:"files"`
}

type Particion struct {
	Nombre string
}

type Carpeta struct {
	Nombre string
}

type Archivo struct {
	Nombre string
}

var carpetas = []Carpeta{}
var archivos = []Archivo{}

func main() {
	Analizar()
}

func mens_error(err error) {
	Salida_comando += errors.New("Error: "+err.Error()).Error() + "\n"
}

func Analizar() {
	fmt.Println("Bienvenido al API del Proyecto2")

	mux := http.NewServeMux()

	mux.HandleFunc("/analizar", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var Content estructuras.Cmd_API
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &Content)
		split_cmd(Content.Cmd)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "` + Salida_comando + `" }`))
		Salida_comando = ""
	})

	mux.HandleFunc("/verparti", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var Content estructuras.Cmd_API
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &Content)
		particiones, err := Discoselec(Content.Cmd)
		if err != nil {
			fmt.Println("Ocurrio un error")
			return
		}
		jsonData, err := json.Marshal(particiones)
		if err != nil {
			fmt.Fprintf(w, "Error marshalling data to JSON: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	mux.HandleFunc("/vercar", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		var Content estructuras.Cmd_API
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &Content)
		particiones, err := Carpetas(Content.Cmd)
		if err != nil {
			fmt.Println("Ocurrio un error")
			return
		}
		fmt.Println(particiones)
		jsonData, err := json.Marshal(particiones)
		if err != nil {
			fmt.Fprintf(w, "Error marshalling data to JSON: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	mux.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request) {
		files, err := ioutil.ReadDir("Discos/MIA/P2")
		if err != nil {
			fmt.Fprintf(w, "Error reading directory: %v", err)
			return
		}

		var fileList FileList
		for _, file := range files {
			if !file.IsDir() {
				fileList.Files = append(fileList.Files, file.Name())
			}
		}
		jsonData, err := json.Marshal(fileList)
		if err != nil {
			fmt.Fprintf(w, "Error marshalling data to JSON: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	fmt.Println("Servidor en el puerto 5000")
	handler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(":5000", handler))
}

func Discoselec(discona string) ([]string, error) {
	var empty [100]byte
	mbr := estructuras.Mbr{}
	path := "/home/nataly/Documentos/Mia lab/Proyecto2/MIA_P2_202001570/Backend/Discos/MIA/P2/" + discona
	disco, err := os.OpenFile(path, os.O_RDWR, 0660)
	particiones := []Particion{}
	if err != nil {
		mens_error(err)
	}
	disco.Seek(0, 0)
	err = binary.Read(disco, binary.BigEndian, &mbr)

	if err != nil {
		mens_error(err)
	}
	if mbr.Mbr_tamano != empty {
		for i := 0; i < 4; i++ {
			name := string(mbr.Mbr_partition[i].Part_name[:])
			name = strings.Trim(name, "\x00")
			if name != "" {
				fmt.Println(name)
				particion := Particion{Nombre: name}
				particiones = append(particiones, particion)
			}
		}
	}

	// Convert particiones slice to []string
	var particionesStr []string
	for _, p := range particiones {
		particionesStr = append(particionesStr, p.Nombre)
	}
	return particionesStr, nil
}

func split_cmd(cmd string) {
	arr_com := strings.Split(cmd, "\n")

	for i := 0; i < len(arr_com); i++ {
		if arr_com[i] != "" {
			Split_comando(arr_com[i])

		}
	}
}

func Split_comando(comando string) {
	var arre_coman []string
	comando = strings.Replace(comando, "\n", "", 1)
	comando = strings.Replace(comando, "\r", "", 1)
	band_comentario := false
	if strings.Contains(comando, "pause") {
		arre_coman = append(arre_coman, comando)
	} else if strings.Contains(comando, "PAUSE") {
		arre_coman = append(arre_coman, comando)
	} else if strings.Contains(comando, "#") {
		band_comentario = true
		//Salida_comando += comando + "\n"
	} else {
		arre_coman = strings.Split(comando, " -")
	}

	if !band_comentario {
		Ejecutar_comando(arre_coman)
	}
}

func Ejecutar_comando(arre_coman []string) {
	data := strings.ToLower(arre_coman[0])

	if data == "mkdisk" {
		/*=======================MKDISK================== */
		comandos.Salid_comando = ""
		comandos.Mkdisk(arre_coman)
		Salida_comando += comandos.Salid_comando
	} else if data == "rmdisk" {
		/*=======================RMDISK================== */
		comandos.Salid_comando = ""
		comandos.Rmdisk(arre_coman)
		Salida_comando += comandos.Salid_comando
	} else if data == "fdisk" {
		/*=======================FDISK=================== */
		// faltan cosas
		comandos.Salid_comando = ""
		comandos.Fdisk(arre_coman)
		Salida_comando += comandos.Salid_comando
	} else if data == "rep" {
		/*=======================REP===================== */
		comandos.Salid_comando = ""
		comandos.Rep(arre_coman)
		Salida_comando += comandos.Salid_comando
	} else if data == "execute" {
		/*=======================EXECUTE================= */
		Execute(arre_coman)
	} else if data == "mkfile" {
		/*========================PAUSE================== */
		//pause()
	} else if data == "mount" {
		/*========================MOUNT================== */
		comandos.Salid_comando = ""
		comandos.Mount(arre_coman)
		Salida_comando += comandos.Salid_comando
	} else if data == "pause" {
		/*========================PAUSE================== */
		pause()
	} else if data == "unmount" {
		/*========================UNMOUNT================== */
		comandos.Salid_comando = ""
		comandos.Unmount(arre_coman)
		Salida_comando += comandos.Salid_comando
	} else if data == "mkfs" {
		/*========================MKFS================== */
		comandos.Salid_comando = ""
		comandos.Mkfs(arre_coman)
		Salida_comando += comandos.Salid_comando
	} else if data == "login" {
		/*========================LOGIN================== */
		comandos.Salid_comando = ""
		comandos.Login(arre_coman)
		Salida_comando += comandos.Salid_comando
	} else if data == "logout" {
		/*========================LOGOUT================== */
		comandos.Salid_comando = ""
		comandos.Logout(arre_coman)
		Salida_comando += comandos.Salid_comando
	} else if data == "mkgrp" {
		/*========================MKFS================== */
		comandos.Salid_comando = ""
		comandos.Mkgrp(arre_coman)
		Salida_comando += comandos.Salid_comando
	} else if data == "mkdir" {
		/*========================MKFS================== */
		comandos.Salid_comando = ""
		comandos.Mkdir(arre_coman)
		Salida_comando += comandos.Salid_comando
	} else {
		/*=======================ERROR=================== */
		Salida_comando += "Error: El comando no fue reconocido." + "\n"
	}
}

func pause() {
	Salida_comando += "Pausa: Presiona enter para continuar..." + "\n"
}

func Execute(arre_coman []string) {
	Salida_comando += "==================EXECUTE=======================" + "\n"
	val_path := ""

	band_path := false
	band_error := false

	for i := 1; i < len(arre_coman); i++ {
		aux_data := strings.SplitAfter(arre_coman[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		switch {
		case strings.Contains(data, "path="):
			if band_path {
				Salida_comando += "Error: El parametro -path ya fue ingresado." + "\n"
				band_error = true
				break
			}

			band_path = true

			val_path = strings.Replace(val_data, "\"", "", 2)
		default:
			Salida_comando += "Error: Parametro no valido." + "\n"
		}
	}

	if !band_error {
		if band_path {

			file, err := os.Open(*&val_path)
			if err != nil {
				Salida_comando += errors.New("Error al abrir el archivo: "+err.Error()).Error() + "\n"
				return
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if line == "" {
					continue
				}
				if strings.HasPrefix(line, "#") {
					continue
				}
				Salida_comando += "\n Comando:" + line + "\n"
				Split_comando(line)
			}
			if err := scanner.Err(); err != nil {
				Salida_comando += errors.New("Error al leer el archivo: "+err.Error()).Error() + "\n"
			}
		}
	}
}

func Carpetas(cadena string) ([]string, error) {
	arreglo := strings.Split(cadena, ",")
	discona := arreglo[0]
	id := arreglo[1]

	path := "/home/nataly/Documentos/Mia lab/Proyecto2/MIA_P2_202001570/Backend/Discos/MIA/P2/" + discona
	mbr := estructuras.Mbr{}
	sb := estructuras.Super_bloque{}
	disco, err := os.OpenFile(path, os.O_RDWR, 0660)
	var empty [100]byte
	if err != nil {
		mens_error(err)
	}
	defer func() {
		disco.Close()
	}()
	disco.Seek(0, 0)
	err = binary.Read(disco, binary.BigEndian, &mbr)
	if mbr.Mbr_tamano != empty {

	}
	var buffer string
	numero_parti := 0
	band_crea := false
	for i := 0; i < 4; i++ {
		s_part_id := string(mbr.Mbr_partition[i].Part_id[:])
		s_part_id = strings.Trim(s_part_id, "\x00")

		if s_part_id == id {
			numero_parti = i
			band_crea = true
			break
		}

	}

	if band_crea {
		s_part_startas := string(mbr.Mbr_partition[numero_parti].Part_start[:])
		s_part_startas = strings.Trim(s_part_startas, "\x00")
		part_starta, err := strconv.Atoi(s_part_startas)
		if err != nil {
			mens_error(err)
		}

		disco.Seek(int64(part_starta), 0)
		err = binary.Read(disco, binary.BigEndian, &sb)
		if sb.S_filesystem_type != empty {
			s_blocks_count := string(sb.S_blocks_count[:])
			s_blocks_count = strings.Trim(s_blocks_count, "\x00")
			numero_block, err := strconv.Atoi(s_blocks_count)
			if err != nil {
				mens_error(err)
			}
			s_free_blocks := string(sb.S_free_blocks_count[:])
			s_free_blocks = strings.Trim(s_free_blocks, "\x00")
			numero_block_libres, err := strconv.Atoi(s_free_blocks)
			if err != nil {
				mens_error(err)
			}

			s_block_start := string(sb.S_block_start[:])
			s_block_start = strings.Trim(s_block_start, "\x00")
			start_block, err := strconv.Atoi(s_block_start)
			if err != nil {
				mens_error(err)
			}

			numero_block_uso := numero_block - numero_block_libres
			buffer += "digraph G{\n"
			for c := 0; c < int(numero_block_uso); c++ {
				blocks, tipo := leerBloques(int64(start_block), path)
				fmt.Println(tipo)
				carpeta := Carpeta{}
				archivo := Archivo{}
				switch b := blocks.(type) {
				case estructuras.Bloque_archivo:
					s_block_cont := string(b.B_content[:])
					s_block_cont = strings.Trim(s_block_cont, "\x00")
					name := "Archivo" + strconv.Itoa(c)
					archivo = Archivo{Nombre: name}
					archivos = append(archivos, archivo)
				case estructuras.Bloque_carpeta:
					for _, c := range b.B_content {
						s_block_name := string(c.B_name[:])
						s_block_name = strings.Trim(s_block_name, "\x00")

						archivo = Archivo{Nombre: s_block_name}
						archivos = append(archivos, archivo)
					}
					name := "Carpeta" + strconv.Itoa(c)
					carpeta = Carpeta{Nombre: name}
					carpetas = append(carpetas, carpeta)
				default:
					fmt.Println("Tipo desconocido")
				}

			}

		}

	}
	var carpetass []string
	for _, p := range carpetas {
		carpetass = append(carpetass, p.Nombre)
	}
	return carpetass, nil

}

func leerBloques(start int64, pathDisco string) (interface{}, int) {
	file, err := os.Open(pathDisco)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Seek(start, 0)

	var char byte
	err = binary.Read(file, binary.LittleEndian, &char)
	if err != nil {
		panic(err)
	}
	if char != '.' {
		file.Seek(start, 0)

		bloque_archivo := estructuras.Bloque_archivo{}
		err = binary.Read(file, binary.LittleEndian, &bloque_archivo.B_content)
		if err != nil {
			panic(err)
		}
		return bloque_archivo, 0
	} else {
		file.Seek(start, 0)

		bloque_carpeta := estructuras.Bloque_carpeta{}

		for c := 0; c < 4; c++ {
			err = binary.Read(file, binary.LittleEndian, &bloque_carpeta.B_content[c].B_name)
			err = binary.Read(file, binary.LittleEndian, &bloque_carpeta.B_content[c].B_inodo)
			if err != nil {
				panic(err)
			}
		}
		return bloque_carpeta, 1
	}
}
