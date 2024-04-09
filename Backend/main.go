package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mimodulo/comandos"
	"mimodulo/estructuras"
	"net/http"
	"os"
	"strings"

	"github.com/rs/cors"
)

var Salida_comando string = ""
var GraphDot string = ""

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
		Split_comando(Content.Cmd)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "` + Salida_comando + `" }`))
		Salida_comando = ""
	})

	fmt.Println("Servidor en el puerto 5000")
	handler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(":5000", handler))
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
		Salida_comando += comando + "\n"
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
		/*========================MKFS================== */
		comandos.Salid_comando = ""
		comandos.Login(arre_coman)
		Salida_comando += comandos.Salid_comando
	} else {
		/*=======================ERROR=================== */
		Salida_comando += "Error: El comando no fue reconocido." + "\n"
	}
}

func pause() {
	Salida_comando += "[MENSAJE] Presiona enter para continuar..." + "\n"
}

func Execute(arre_coman []string) {
	Salida_comando += "==============EXECUTE=======================" + "\n"
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
