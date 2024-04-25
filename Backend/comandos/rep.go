package comandos

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"mimodulo/estructuras"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var val_rutadis string = "/home/nataly/Documentos/Mia lab/Proyecto2/MIA_P2_202001570/Backend/Discos/MIA/P2/"

func Rep(arre_coman []string) {
	//Salid_comando += "=======================REP=======================" + "\n"

	val_path := ""
	band_path := false

	val_id := ""
	band_id := false

	val_ruta := ""
	band_ruta := false

	val_name := ""
	band_name := false

	band_error := false

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
			val_name = strings.Replace(val_data, "\"", "", 2)
			val_name = strings.ToLower(val_name)
			if val_name == "mbr" || val_name == "disk" || val_name == "inode" || val_name == "journaling" {
				band_name = true
			} else if val_name == "block" || val_name == "bm_inode" || val_name == "bm_block" || val_name == "tree" {
				band_name = true
			} else if val_name == "sb" || val_name == "file" || val_name == "ls" {
				band_name = true
			} else {
				Salid_comando += "Error: El Valor del parametro -name no es valido" + "\n"
				band_error = true
				break
			}
		case strings.Contains(data, "path="):
			if band_path {
				Salid_comando += "Error: El parametro -path ya fue ingresado." + "\n"
				band_error = true
				break
			}
			val_path = val_data
			band_path = true
		case strings.Contains(data, "id="):
			if band_id {
				Salid_comando += "Error: El parametro -id ya fue ingresado." + "\n"
				band_error = true
				break
			}
			val_id = val_data
			band_id = true
		case strings.Contains(data, "ruta="):
			if band_ruta {
				Salid_comando += "Error: El parametro -ruta ya fue ingresado." + "\n"
				band_error = true
				break
			}
			val_ruta = val_data
			band_ruta = true
			fmt.Print("Ruta: ", val_ruta)
		default:
			Salid_comando += "Error: Parametro no valido." + "\n"
		}
	}

	if !band_error {
		if band_path {
			if band_name {
				if band_id {
					if val_name == "mbr" {
						RepoMbr(val_path, val_id)
					} else if val_name == "disk" {
						RepoDisk(val_path, val_id)
					} else if val_name == "sb" {
						ReporteSb(val_path, val_id)
					} else if val_name == "inode" {
						ReporteInode(val_path, val_id)
					} else if val_name == "block" {
						ReporteBloque(val_path, val_id)
					} else if val_name == "bm_inode" {
						ReporteBmInode(val_path, val_id)
					} else if val_name == "bm_block" {
						ReporteBmBlock(val_path, val_id)
					} else if val_name == "tree" {
						ReporteTree(val_path, val_id)
					}
				} else {
					Salid_comando += "Error: No se encontro el id" + "\n"
				}
			} else {
				Salid_comando += "Error: No se encontro el name" + "\n"
			}
		} else {
			Salid_comando += "Error: No se encontro el path" + "\n"
		}
	}
}

func RepoMbr(path string, id string) {
	carpeta := ""
	archivo := ""

	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			carpeta = path[:i]
			archivo = path[i+1:]
			break
		}
	}
	fmt.Println(archivo)
	//Salid_comando += "Carpeta:" + carpeta + "\n"
	//Salid_comando += "Archivo:" + archivo + "\n"
	//Salid_comando += "Disco:" + id + "\n"
	val_rutadis = val_rutadis + id + ".dsk"

	_, err := os.Stat(carpeta)

	// Si la carpeta no existe, crearla
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(carpeta, 0755)
			if err != nil {
				Salid_comando += errors.New("Error al crear la carpeta: "+err.Error()).Error() + "\n"
				return
			}
		} else {
			Salid_comando += errors.New("Error al verificar la carpeta: "+err.Error()).Error() + "\n"
			return
		}
	}

	//Salid_comando += "La carpeta" + carpeta + "ya existe o se ha creado correctamente" + "\n"

	var buffer string
	buffer += "digraph G{\nsubgraph cluster{\nlabel=\"REPORTE DE MBR\"\ntbl[shape=box,label=<\n<table border='0' cellborder='1' cellspacing='0' width='300'  height='200' >\n"
	var empty [100]byte
	mbr := estructuras.Mbr{}

	bandera_extendida := false
	numero_exten := 0

	disco, err := os.OpenFile(val_rutadis, os.O_RDWR, 0660)

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
		s_part_tamaño := string(mbr.Mbr_tamano[:])
		s_part_tamaño = strings.Trim(s_part_tamaño, "\x00")
		buffer += "<tr> <td width='150' bgcolor=\"pink\"><b>Mbr_tamano</b></td><td width='150' bgcolor=\"pink\">" + s_part_tamaño + "</td>  </tr>\n"
		s_part_crea := string(mbr.Mbr_fecha_creacion[:])
		s_part_crea = strings.Trim(s_part_crea, "\x00")
		buffer += "<tr>  <td bgcolor=\"pink\"><b>Mbr_fecha_creacion</b></td><td bgcolor=\"pink\">" + s_part_crea + "</td>  </tr>\n"
		s_part_sig := string(mbr.Mbr_dsk_signature[:])
		s_part_sig = strings.Trim(s_part_sig, "\x00")
		buffer += "<tr>  <td bgcolor=\"pink\"><b>Mbr_dsk_signature</b></td><td bgcolor=\"pink\">" + s_part_sig + "</td>  </tr>\n"
		s_part_fit := string(mbr.Dsk_fit[:])
		s_part_fit = strings.Trim(s_part_fit, "\x00")
		buffer += "<tr>  <td bgcolor=\"pink\"><b>Dsk_fit</b></td><td bgcolor=\"pink\">" + s_part_fit + "</td>  </tr>"
		for i := 0; i < 4; i++ {
			s_part_statas := string(mbr.Mbr_partition[i].Part_status[:])
			s_part_statas = strings.Trim(s_part_statas, "\x00")
			s_part_startas := string(mbr.Mbr_partition[i].Part_start[:])
			s_part_startas = strings.Trim(s_part_startas, "\x00")
			if s_part_statas != "E" && s_part_startas != "-1" {
				buffer += "<tr>"
				buffer += "<td colspan=\"2\" bgcolor=\"purple\">Particion</td>"
				buffer += "</tr>"
				s_part_stat := string(mbr.Mbr_partition[i].Part_status[:])
				s_part_stat = strings.Trim(s_part_stat, "\x00")
				buffer += "<tr>  <td bgcolor=\"pink\"><b>Part_status" + "</b></td><td bgcolor=\"pink\">" + s_part_stat + "</td>  </tr>\n"
				s_part_type := string(mbr.Mbr_partition[i].Part_type[:])
				s_part_type = strings.Trim(s_part_type, "\x00")
				buffer += "<tr>  <td bgcolor=\"pink\"><b>Part_type" + "</b></td><td bgcolor=\"pink\">" + s_part_type + "</td>  </tr>\n"
				s_part_fi := string(mbr.Mbr_partition[i].Part_fit[:])
				s_part_fi = strings.Trim(s_part_fi, "\x00")
				buffer += "<tr>  <td bgcolor=\"pink\"><b>Part_fit" + "</b></td><td bgcolor=\"pink\">" + s_part_fi + "</td>  </tr>\n"
				s_part_st := string(mbr.Mbr_partition[i].Part_start[:])
				s_part_st = strings.Trim(s_part_st, "\x00")
				buffer += "<tr>  <td bgcolor=\"pink\"><b>Part_start " + "</b></td><td bgcolor=\"pink\">" + s_part_st + "</td>  </tr>\n"
				s_part_siz := string(mbr.Mbr_partition[i].Part_size[:])
				s_part_siz = strings.Trim(s_part_siz, "\x00")
				buffer += "<tr>  <td bgcolor=\"pink\"><b>Part_size" + "</b></td><td bgcolor=\"pink\">" + s_part_siz + "</td>  </tr>\n"
				s_part_name := string(mbr.Mbr_partition[i].Part_name[:])
				s_part_name = strings.Trim(s_part_name, "\x00")
				buffer += "<tr>  <td bgcolor=\"pink\"><b>Part_name" + "</b></td><td bgcolor=\"pink\">" + s_part_name + "</td>  </tr>\n"
				if s_part_type == "e" {
					numero_exten = i
					bandera_extendida = true
				}
			}

		}
		buffer += ("</table>\n>];\n}")
		if bandera_extendida {
			s_part_start := string(mbr.Mbr_partition[numero_exten].Part_start[:])
			s_part_start = strings.Trim(s_part_start, "\x00")
			part_start, err := strconv.Atoi(s_part_start)
			if err != nil {
				Mens_error(err)
			}
			ebr := obtener_ebr(val_rutadis, part_start)
			bandera_bucle := false

			u := 0
			buffer += ("subgraph cluster_1{\n label=\"Particiones Logicas dentro de la extendida\"\ntbl_1[shape=box, label=<\n<table border='0' cellborder='1' cellspacing='0'  width='300' height='160' >\n")
			for !bandera_bucle {
				if ebr.Part_status != empty {
					s_part_ebnexts := string(ebr.Part_next[:])
					s_part_ebnexts = strings.Trim(s_part_ebnexts, "\x00")
					if s_part_ebnexts != "-1" && s_part_ebnexts != "0" {
						buffer += "<tr>"
						buffer += "<td width='150' colspan=\"2\" bgcolor=\"magenta\">Particion</td>"
						buffer += "</tr>"
						s_part_ebfit := string(ebr.Part_fit[:])
						s_part_ebfit = strings.Trim(s_part_ebfit, "\x00")
						buffer += "<tr>  <td width='150' bgcolor=\"pink\"><b>Part_fit" + "</b></td><td width='150' bgcolor=\"pink\">" + s_part_ebfit + "</td>  </tr>\n"
						s_part_ebst := string(ebr.Part_start[:])
						s_part_ebst = strings.Trim(s_part_ebst, "\x00")
						buffer += "<tr>  <td bgcolor=\"pink\"><b>Part_start " + "</b></td><td bgcolor=\"pink\">" + s_part_ebst + "</td>  </tr>\n"
						s_part_ebnext := string(ebr.Part_next[:])
						s_part_ebnext = strings.Trim(s_part_ebnext, "\x00")
						buffer += "<tr>  <td bgcolor=\"pink\"><b>Part_next" + "</b></td><td bgcolor=\"pink\">" + s_part_ebnext + "</td>  </tr>\n"
						s_part_ebsiz := string(ebr.Part_size[:])
						s_part_ebsiz = strings.Trim(s_part_ebsiz, "\x00")
						buffer += "<tr>  <td bgcolor=\"pink\"><b>Part_size " + "</b></td><td bgcolor=\"pink\">" + s_part_ebsiz + "</td>  </tr>\n"
						s_part_ebname := string(ebr.Part_name[:])
						s_part_ebname = strings.Trim(s_part_ebname, "\x00")
						buffer += "<tr>  <td bgcolor=\"pink\"><b>Part_name " + "</b></td><td bgcolor=\"pink\">" + s_part_ebname + "</td>  </tr>\n"
					}

					u += 1
				} else {
					break
				}
				s_part_next := string(ebr.Part_next[:])
				s_part_next = strings.Trim(s_part_next, "\x00")
				part_next, err := strconv.Atoi(s_part_next)
				if err != nil {
					Mens_error(err)
				}
				if part_next == -1 {
					bandera_bucle = true
				} else {
					ebr = obtener_ebr(val_rutadis, part_next)
				}

			}
			if u != 0 {
				buffer += "</table>\n>];\n}"
			}
		}
		buffer += "}\n"
		Salid_comando += buffer
		file, err2 := os.Create("Mbr.dot")
		if err2 != nil && !os.IsExist(err) {
			log.Fatal(err2)
		}
		defer file.Close()

		err = os.Chmod("Mbr.dot", 0777)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}

		_, err = file.WriteString(buffer)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}
		cmd := exec.Command("dot", "-Tpng", "Mbr.dot", "-o", path)
		err = cmd.Run()
		if err != nil {
			fmt.Errorf("no se pudo generar la imagen: %v", err)
		}
		val_rutadis = "/home/nataly/Documentos/Mia lab/Proyecto2/MIA_P2_202001570/Backend/Discos/MIA/P2/"
		//Salid_comando += "Reporte Mbr-Ebr creado exitosamente" + "\n"
	}

}

func obtener_ebr(ruta string, part_start int) estructuras.Ebr {
	var empty [100]byte
	ebr_empty := estructuras.Ebr{}
	ebr := estructuras.Ebr{}

	// Apertura de archivo
	disco, err := os.OpenFile(ruta, os.O_RDWR, 0660)

	// ERROR
	if err != nil {
		Mens_error(err)
	}
	defer func() {
		disco.Close()
	}()
	disco.Seek(int64(part_start), 0)
	err = binary.Read(disco, binary.BigEndian, &ebr)

	if err != nil {
		Mens_error(err)
	}

	if ebr.Part_next != empty {
		return ebr
	} else {
		return ebr_empty
	}
}

func ReporteSb(path string, id string) {
	carpeta := ""
	archivo := ""
	regex := regexp.MustCompile(`^[a-zA-Z]+`)
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			carpeta = path[:i]
			archivo = path[i+1:]
			break
		}
	}
	fmt.Println(archivo)
	//Salid_comando += "Carpeta:" + carpeta + "\n"
	//Salid_comando += "Archivo:" + archivo + "\n"
	letters := regex.FindString(id)
	//Salid_comando += "Disco:" + letters + "\n"

	val_rutadis = val_rutadis + letters + ".dsk"

	_, err := os.Stat(carpeta)

	// Si la carpeta no existe, crearla
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(carpeta, 0755)
			if err != nil {
				Salid_comando += errors.New("Error al crear la carpeta: "+err.Error()).Error() + "\n"
				return
			}
		} else {
			Salid_comando += errors.New("Error al verificar la carpeta: "+err.Error()).Error() + "\n"
			return
		}
	}

	//Salid_comando += "La carpeta" + carpeta + "ya existe o se ha creado correctamente" + "\n"
	mbr := estructuras.Mbr{}
	sb := estructuras.Super_bloque{}
	disco, err := os.OpenFile(val_rutadis, os.O_RDWR, 0660)
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
			Mens_error(err)
		}

		disco.Seek(int64(part_starta), 0)
		err = binary.Read(disco, binary.BigEndian, &sb)
		if sb.S_filesystem_type != empty {
			buffer += "digraph G{\nsubgraph cluster{\nlabel=\"Super Bloque\"\ntbl[shape=box,label=<\n<table border='0' cellborder='1' cellspacing='0' width='300'  height='200' >\n"
			S_filesystem_type := string(sb.S_filesystem_type[:])
			S_filesystem_type = strings.Trim(S_filesystem_type, "\x00")
			buffer += "<tr> <td bgcolor=\"green\"><b>S_filesystem_type </b></td><td bgcolor=\"green\">" + S_filesystem_type + "</td> </tr>*"
			S_inodes_count := string(sb.S_inodes_count[:])
			S_inodes_count = strings.Trim(S_inodes_count, "\x00")
			buffer += "<tr><td><b>S_inodes_count </b></td><td>" + S_inodes_count + "</td> </tr>\n"
			S_blocks_count := string(sb.S_blocks_count[:])
			S_blocks_count = strings.Trim(S_blocks_count, "\x00")
			buffer += "<tr><td bgcolor=\"green\"><b>S_blocks_count </b></td><td bgcolor=\"green\">" + S_blocks_count + "</td> </tr>\n"
			S_free_blocks_count := string(sb.S_free_blocks_count[:])
			S_free_blocks_count = strings.Trim(S_free_blocks_count, "\x00")
			buffer += "<tr><td><b>S_free_blocks_count</b></td><td>" + S_free_blocks_count + "</td> </tr>\n"
			S_free_inodes_count := string(sb.S_free_inodes_count[:])
			S_free_inodes_count = strings.Trim(S_free_inodes_count, "\x00")
			buffer += "<tr><td bgcolor=\"green\"><b>S_free_inodes_count</b></td><td bgcolor=\"green\">" + S_free_inodes_count + "</td> </tr>\n"
			S_mtime := string(sb.S_mtime[:])
			S_mtime = strings.Trim(S_mtime, "\x00")
			buffer += "<tr><td><b>S_mtime</b></td><td>" + S_mtime + "</td> </tr>\n"
			S_umtime := string(sb.S_umtime[:])
			S_umtime = strings.Trim(S_umtime, "\x00")
			buffer += "<tr><td bgcolor=\"green\"><b>S_umtime</b></td><td bgcolor=\"green\">" + S_umtime + "</td> </tr>\n"
			S_mnt_count := string(sb.S_mnt_count[:])
			S_mnt_count = strings.Trim(S_mnt_count, "\x00")
			buffer += "<tr><td><b>S_mnt_count</b></td><td>" + S_mnt_count + "</td> </tr>\n"
			S_magic := string(sb.S_magic[:])
			S_magic = strings.Trim(S_magic, "\x00")
			buffer += "<tr><td bgcolor=\"green\"><b>S_magic</b></td><td bgcolor=\"green\">" + S_magic + "</td> </tr>\n"
			S_inode_size := string(sb.S_inode_size[:])
			S_inode_size = strings.Trim(S_inode_size, "\x00")
			buffer += "<tr><td><b>S_inode_size</b></td><td>" + S_inode_size + "</td> </tr>\n"
			S_block_size := string(sb.S_block_size[:])
			S_block_size = strings.Trim(S_block_size, "\x00")
			buffer += "<tr><td bgcolor=\"green\"><b>Sb_date_ultimo_montaje</b></td><td bgcolor=\"green\">" + S_block_size + "</td> </tr>\n"
			S_firts_ino := string(sb.S_firts_ino[:])
			S_firts_ino = strings.Trim(S_firts_ino, "\x00")
			buffer += "<tr><td><b>S_firts_ino</b></td><td>" + S_firts_ino + "</td> </tr>\n"
			S_first_blo := string(sb.S_first_blo[:])
			S_first_blo = strings.Trim(S_first_blo, "\x00")
			buffer += "<tr><td bgcolor=\"green\"><b>S_first_blo</b></td><td bgcolor=\"green\">" + S_first_blo + "</td> </tr>\n"
			S_bm_inode_start := string(sb.S_bm_inode_start[:])
			S_bm_inode_start = strings.Trim(S_bm_inode_start, "\x00")
			buffer += "<tr><td><b>S_bm_inode_start</b></td><td>" + S_bm_inode_start + "</td> </tr>\n"
			S_bm_block_start := string(sb.S_bm_block_start[:])
			S_bm_block_start = strings.Trim(S_bm_block_start, "\x00")
			buffer += "<tr><td bgcolor=\"green\"><b>S_bm_block_start</b></td><td bgcolor=\"green\">" + S_bm_block_start + "</td> </tr>\n"
			S_inode_start := string(sb.S_inode_start[:])
			S_inode_start = strings.Trim(S_inode_start, "\x00")
			buffer += "<tr><td><b>S_inode_start</b></td><td>" + S_inode_start + "</td> </tr>\n"
			S_block_start := string(sb.S_block_start[:])
			S_block_start = strings.Trim(S_block_start, "\x00")
			buffer += "<tr><td bgcolor=\"green\"><b>S_block_start</b></td><td bgcolor=\"green\">" + S_block_start + "</td> </tr>\n"
			buffer += "</table>\n>];\n}"
			buffer += "}\n"
		}

	}
	Salid_comando += buffer
	file, err2 := os.Create("Super.dot")
	if err2 != nil && !os.IsExist(err) {
		log.Fatal(err2)
	}
	defer file.Close()

	err = os.Chmod("Super.dot", 0777)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	_, err = file.WriteString(buffer)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	cmd := exec.Command("dot", "-Tpng", "Super.dot", "-o", path)
	err = cmd.Run()
	if err != nil {
		fmt.Errorf("no se pudo generar la imagen: %v", err)
	}
	val_rutadis = "/home/nataly/Documentos/Mia lab/Proyecto2/MIA_P2_202001570/Backend/Discos/MIA/P2/"
	//Salid_comando += "Reporte Super Blocke creado exitosamente" + "\n"
}

func RepoDisk(path string, id string) {
	var buffer string
	mbr := estructuras.Mbr{}
	val_rutadis = val_rutadis + id + ".dsk"
	buffer += "digraph G{\ntbl [\nshape=box\nlabel=<\n<table border='0' cellborder='2' width='100' height=\"30\" color='lightblue4'>\n<tr>"
	f, err := os.OpenFile(val_rutadis, os.O_RDWR, 0755)
	if err != nil {
		Salid_comando += "Error: No existe la ruta" + "\n"
	}
	defer f.Close()
	PorcentajeUtilizao := 0.0
	var EspacioUtilizado int64 = 0
	disco, err := os.OpenFile(val_rutadis, os.O_RDWR, 0660)

	if err != nil {
		Mens_error(err)
	}
	defer func() {
		disco.Close()
	}()

	disco.Seek(0, 0)
	err = binary.Read(disco, binary.BigEndian, &mbr)

	TamanioDisco := string(mbr.Mbr_tamano[:])
	TamanioDisco = strings.Trim(TamanioDisco, "\x00")
	tamdis, err := strconv.Atoi(TamanioDisco)
	if err != nil {
		print("Aqw")
		Mens_error(err)
	}
	buffer += "<td height='30' width='75'> MBR </td>"
	for i := 0; i < 4; i++ {
		s_part_partina := string(mbr.Mbr_partition[i].Part_type[:])
		s_part_partina = strings.Trim(s_part_partina, "\x00")
		s_part_partna := string(mbr.Mbr_partition[i].Part_name[:])
		s_part_partna = strings.Trim(s_part_partna, "\x00")
		s_part_partstatus := string(mbr.Mbr_partition[i].Part_status[:])
		s_part_partstatus = strings.Trim(s_part_partstatus, "\x00")
		if s_part_partina == "p" {
			s_part_partsiz := string(mbr.Mbr_partition[i].Part_size[:])
			s_part_partsiz = strings.Trim(s_part_partsiz, "\x00")
			part_size, err := strconv.Atoi(s_part_partsiz)
			if err != nil {
				print("Aqw1")
				Mens_error(err)
			}
			PorcentajeUtilizao = (float64(part_size) / float64(tamdis)) * 100
			buffer += "<td height='30' width='75.0'>PRIMARIA <br/>" + s_part_partna + " <br/> Ocupado: " + strconv.Itoa(int(PorcentajeUtilizao)) + "%</td>"
			EspacioUtilizado += int64(part_size)
		} else if s_part_partstatus == "0" {
			buffer += "<td height='30' width='75.0'>Libre</td>"
		}
		if s_part_partina == "e" {
			s_part_partsiz := string(mbr.Mbr_partition[i].Part_size[:])
			s_part_partsiz = strings.Trim(s_part_partsiz, "\x00")
			part_size, err := strconv.Atoi(s_part_partsiz)
			if err != nil {
				print("Aqw3")
				Mens_error(err)
			}
			EspacioUtilizado += int64(part_size)
			PorcentajeUtilizao = (float64(part_size) / float64(tamdis)) * 100
			buffer += "<td  height='30' width='15.0'>\n"
			buffer += "<table border='5'  height='30' WIDTH='15.0' cellborder='1'>\n"
			buffer += " <tr>  <td height='60' colspan='100%'>EXTENDIDA <br/>" + s_part_partna + " <br/> Ocupado:" + strconv.Itoa(int(PorcentajeUtilizao)) + "%</td>  </tr>\n<tr>"
			s_part_partsart := string(mbr.Mbr_partition[i].Part_start[:])
			s_part_partsart = strings.Trim(s_part_partsart, "\x00")
			part_start, err := strconv.Atoi(s_part_partsart)
			if err != nil {
				print("Aqw4")
				Mens_error(err)
			}
			var InicioExtendida int64 = int64(part_start)
			disco.Seek(int64(InicioExtendida), 0)
			ebr := estructuras.Ebr{}
			err = binary.Read(disco, binary.BigEndian, &ebr)

			if err != nil {
				Mens_error(err)
			}
			bandera_bucle := false
			var empty [100]byte
			for !bandera_bucle {
				ebr_next := string(ebr.Part_next[:])
				ebr_next = strings.Trim(ebr_next, "\x00")
				if ebr.Part_status != empty {
					if ebr_next != "-1" {
						ebr_size := string(ebr.Part_size[:])
						ebr_size = strings.Trim(ebr_size, "\x00")
						part_sizes, err := strconv.Atoi(ebr_size)
						if err != nil {
							print("Aqw5")
							Mens_error(err)
						}

						EspacioUtilizado += int64(part_sizes)
						PorcentajeUtilizao = (float64(part_sizes) / float64(part_size)) * 100
						ebr_name := string(ebr.Part_name[:])
						ebr_name = strings.Trim(ebr_name, "\x00")
						buffer += "<td height='30'>EBR</td><td height='30'> Logica:  " + ebr_name + " " + strconv.Itoa(int(PorcentajeUtilizao)) + "%</td>"
					}

				} else {
					break
				}
				s_part_next := string(ebr.Part_next[:])
				s_part_next = strings.Trim(s_part_next, "\x00")
				part_next, err := strconv.Atoi(s_part_next)
				if err != nil {
					print("Aqw6")
					Mens_error(err)
				}
				if part_next == -1 {
					bandera_bucle = true
				} else {
					ebr = obtener_ebr(val_rutadis, part_next)
				}

			}
			if (int64(part_size) - EspacioUtilizado) > 0 {
				PorcentajeUtilizao = (float64(int64(tamdis)-EspacioUtilizado) / float64(tamdis)) * 100
				buffer += "<td height='30' width='100%'>Libre: " + strconv.Itoa(int(PorcentajeUtilizao)) + "%</td>"
			}

			buffer += "</tr>\n"
			buffer += "</table>\n</td>"

		}
	}
	if (int64(tamdis) - EspacioUtilizado) > 0 {
		PorcentajeUtilizao = (float64(int64(tamdis)-EspacioUtilizado) / float64(tamdis)) * 100
		buffer += "<td height='30' width='75.0'>Libre: " + strconv.Itoa(int(PorcentajeUtilizao)) + "%</td>"
	}
	buffer += "     </tr>\n</table>\n>];\n}"
	Salid_comando += buffer
	file, err2 := os.Create("Disk.dot")
	if err2 != nil && !os.IsExist(err) {
		log.Fatal(err2)
	}
	defer file.Close()

	err = os.Chmod("Disk.dot", 0777)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	_, err = file.WriteString(buffer)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	cmd := exec.Command("dot", "-Tpng", "Disk.dot", "-o", path)
	err = cmd.Run()
	if err != nil {
		fmt.Errorf("no se pudo generar la imagen: %v", err)
	}
	val_rutadis = "/home/nataly/Documentos/Mia lab/Proyecto2/MIA_P2_202001570/Backend/Discos/MIA/P2/"
	//Salid_comando += "Reporte Disk creado exitosamente" + "\n"
}

func ReporteInode(path string, id string) {
	carpeta := ""
	archivo := ""
	regex := regexp.MustCompile(`^[a-zA-Z]+`)
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			carpeta = path[:i]
			archivo = path[i+1:]
			break
		}
	}
	fmt.Println(archivo)
	//Salid_comando += "Carpeta:" + carpeta + "\n"
	//Salid_comando += "Archivo:" + archivo + "\n"
	letters := regex.FindString(id)
	//Salid_comando += "Disco:" + letters + "\n"

	val_rutadis = val_rutadis + letters + ".dsk"

	_, err := os.Stat(carpeta)

	// Si la carpeta no existe, crearla
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(carpeta, 0755)
			if err != nil {
				Salid_comando += errors.New("Error al crear la carpeta: "+err.Error()).Error() + "\n"
				return
			}
		} else {
			Salid_comando += errors.New("Error al verificar la carpeta: "+err.Error()).Error() + "\n"
			return
		}
	}
	//Salid_comando += "La carpeta" + carpeta + "ya existe o se ha creado correctamente" + "\n"
	mbr := estructuras.Mbr{}
	sb := estructuras.Super_bloque{}
	disco, err := os.OpenFile(val_rutadis, os.O_RDWR, 0660)
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
			Mens_error(err)
		}

		disco.Seek(int64(part_starta), 0)
		err = binary.Read(disco, binary.BigEndian, &sb)
		if sb.S_filesystem_type != empty {
			s_inodos_count := string(sb.S_inodes_count[:])
			s_inodos_count = strings.Trim(s_inodos_count, "\x00")
			numero_inodos, err := strconv.Atoi(s_inodos_count)
			if err != nil {
				Mens_error(err)
			}
			s_free_inodos := string(sb.S_free_inodes_count[:])
			s_free_inodos = strings.Trim(s_free_inodos, "\x00")
			numero_inodos_libres, err := strconv.Atoi(s_free_inodos)
			if err != nil {
				Mens_error(err)
			}

			s_inodos_start := string(sb.S_inode_start[:])
			s_inodos_start = strings.Trim(s_inodos_start, "\x00")
			start_inodos, err := strconv.Atoi(s_inodos_start)
			if err != nil {
				Mens_error(err)
			}

			numero_inodos_uso := numero_inodos - numero_inodos_libres
			buffer += "digraph G{\n"
			for c := 0; c < int(numero_inodos_uso); c++ {
				inodos := leerInodos(int64(start_inodos), val_rutadis)
				buffer += "subgraph cluster_" + strconv.Itoa(c) + "{\n label=\"Inodo" + strconv.Itoa(c) + "\"\ntbl_" + strconv.Itoa(c) + "[shape=box, label=<\n<table border='0' cellborder='1' cellspacing='0'  width='300' height='160' >\n"
				s_inodo_uid := string(inodos.I_uid[:])
				s_inodo_uid = strings.Trim(s_inodo_uid, "\x00")
				buffer += "<tr><td width='150' bgcolor=\"pink\"><b>I_uid" + "</b></td><td width='150' bgcolor=\"pink\">" + s_inodo_uid + "</td></tr>\n"
				s_inodo_gid := string(inodos.I_gid[:])
				s_inodo_gid = strings.Trim(s_inodo_gid, "\x00")
				buffer += "<tr><td width='150' bgcolor=\"pink\"><b>I_gid" + "</b></td><td width='150' bgcolor=\"pink\">" + s_inodo_gid + "</td></tr>\n"
				s_inodo_size := string(inodos.I_size[:])
				s_inodo_size = strings.Trim(s_inodo_size, "\x00")
				buffer += "<tr><td width='150' bgcolor=\"pink\"><b>I_size" + "</b></td><td width='150' bgcolor=\"pink\">" + s_inodo_size + "</td></tr>\n"
				s_inodo_atime := string(inodos.I_atime[:])
				s_inodo_atime = strings.Trim(s_inodo_atime, "\x00")
				buffer += "<tr><td width='150' bgcolor=\"pink\"><b>I_atime" + "</b></td><td width='150' bgcolor=\"pink\">\"" + s_inodo_atime + "\"</td></tr>\n"
				s_inodo_ctime := string(inodos.I_ctime[:])
				s_inodo_ctime = strings.Trim(s_inodo_ctime, "\x00")
				buffer += "<tr><td width='150' bgcolor=\"pink\"><b>I_ctime" + "</b></td><td width='150' bgcolor=\"pink\">\"" + s_inodo_ctime + "\"</td></tr>\n"
				s_inodo_mtime := string(inodos.I_mtime[:])
				s_inodo_mtime = strings.Trim(s_inodo_mtime, "\x00")
				buffer += "<tr><td width='150' bgcolor=\"pink\"><b>I_mtime" + "</b></td><td width='150' bgcolor=\"pink\">\"" + s_inodo_mtime + "\"</td></tr>\n"
				// bloques
				for i := 1; i < 15; i++ {
					s_inodo_bloc := string(inodos.I_block[i])
					s_inodo_bloc = strings.Trim(s_inodo_bloc, "\x00")
					buffer += "<tr>  <td width='150' bgcolor=\"pink\"><b>I_bloc" + strconv.Itoa(i) + "</b></td><td width='150' bgcolor=\"pink\">" + s_inodo_bloc + "</td>  </tr>\n"

				}
				s_inodo_type := string(inodos.I_type[:])
				s_inodo_type = strings.Trim(s_inodo_type, "\x00")
				buffer += "<tr><td width='150' bgcolor=\"pink\"><b>I_type" + "</b></td><td width='150' bgcolor=\"pink\">" + s_inodo_type + "</td></tr>\n"
				s_inodo_perm := string(inodos.I_perm[:])
				s_inodo_perm = strings.Trim(s_inodo_perm, "\x00")
				buffer += "<tr><td width='150' bgcolor=\"pink\"><b>I_perm" + "</b></td><td width='150' bgcolor=\"pink\">" + s_inodo_perm + "</td></tr>\n"
				buffer += "</table>>];}\n"
			}
			buffer += "}\n"
			Salid_comando += buffer
			file, err2 := os.Create("Inodo.dot")
			if err2 != nil && !os.IsExist(err) {
				log.Fatal(err2)
			}
			defer file.Close()

			err = os.Chmod("Inodo.dot", 0777)
			if err != nil && !os.IsExist(err) {
				log.Fatal(err)
			}

			_, err = file.WriteString(buffer)
			if err != nil && !os.IsExist(err) {
				log.Fatal(err)
			}
			cmd := exec.Command("dot", "-Tpng", "Inodo.dot", "-o", path)
			err = cmd.Run()
			if err != nil {
				fmt.Errorf("no se pudo generar la imagen: %v", err)
			}
			val_rutadis = "/home/nataly/Documentos/Mia lab/Proyecto2/MIA_P2_202001570/Backend/Discos/MIA/P2/"
			//Salid_comando += "Reporte Inodo creado exitosamente" + "\n"
		}

	}

}

func leerInodos(start int64, pathDisco string) estructuras.Inodo {
	file, err := os.Open(pathDisco)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Seek(start, 0)

	inodo := estructuras.Inodo{}
	err = binary.Read(file, binary.LittleEndian, &inodo)
	if err != nil {
		Mens_error(err)
	}

	return inodo

}

func ReporteBloque(path string, id string) {
	carpeta := ""
	archivo := ""
	regex := regexp.MustCompile(`^[a-zA-Z]+`)
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			carpeta = path[:i]
			archivo = path[i+1:]
			break
		}
	}
	fmt.Println(archivo)
	//Salid_comando += "Carpeta:" + carpeta + "\n"
	//Salid_comando += "Archivo:" + archivo + "\n"
	letters := regex.FindString(id)
	//Salid_comando += "Disco:" + letters + "\n"

	val_rutadis = val_rutadis + letters + ".dsk"

	_, err := os.Stat(carpeta)

	// Si la carpeta no existe, crearla
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(carpeta, 0755)
			if err != nil {
				Salid_comando += errors.New("Error al crear la carpeta: "+err.Error()).Error() + "\n"
				return
			}
		} else {
			Salid_comando += errors.New("Error al verificar la carpeta: "+err.Error()).Error() + "\n"
			return
		}
	}
	//Salid_comando += "La carpeta" + carpeta + "ya existe o se ha creado correctamente" + "\n"
	mbr := estructuras.Mbr{}
	sb := estructuras.Super_bloque{}
	disco, err := os.OpenFile(val_rutadis, os.O_RDWR, 0660)
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
			Mens_error(err)
		}

		disco.Seek(int64(part_starta), 0)
		err = binary.Read(disco, binary.BigEndian, &sb)
		if sb.S_filesystem_type != empty {
			s_blocks_count := string(sb.S_blocks_count[:])
			s_blocks_count = strings.Trim(s_blocks_count, "\x00")
			numero_block, err := strconv.Atoi(s_blocks_count)
			if err != nil {
				Mens_error(err)
			}
			s_free_blocks := string(sb.S_free_blocks_count[:])
			s_free_blocks = strings.Trim(s_free_blocks, "\x00")
			numero_block_libres, err := strconv.Atoi(s_free_blocks)
			if err != nil {
				Mens_error(err)
			}

			s_block_start := string(sb.S_block_start[:])
			s_block_start = strings.Trim(s_block_start, "\x00")
			start_block, err := strconv.Atoi(s_block_start)
			if err != nil {
				Mens_error(err)
			}

			numero_block_uso := numero_block - numero_block_libres
			buffer += "digraph G{\n"
			for c := 0; c < int(numero_block_uso); c++ {
				blocks, tipo := leerBloques(int64(start_block), val_rutadis)
				fmt.Println(tipo)
				switch b := blocks.(type) {
				case estructuras.Bloque_archivo:
					buffer += "subgraph cluster_" + strconv.Itoa(c) + "{\n label=\"Bloque Archivo" + strconv.Itoa(c) + "\"\ntbl_" + strconv.Itoa(c) + "[shape=box, label=<\n<table border='0' cellborder='1' cellspacing='0'  width='300' height='160' >\n"
					s_block_cont := string(b.B_content[:])
					s_block_cont = strings.Trim(s_block_cont, "\x00")
					buffer += "<tr><td width='150' bgcolor=\"pink\"><b>B_content" + "</b></td><td width='150' bgcolor=\"pink\">" + s_block_cont + "</td></tr>\n"
					buffer += "</table>>];}\n"
				case estructuras.Bloque_carpeta:
					buffer += "subgraph cluster_" + strconv.Itoa(c) + "{\n label=\"Bloque Carpeta" + strconv.Itoa(c) + "\"\ntbl_" + strconv.Itoa(c) + "[shape=box, label=<\n<table border='0' cellborder='1' cellspacing='0'  width='300' height='160' >\n"
					for _, c := range b.B_content {
						s_block_name := string(c.B_name[:])
						s_block_name = strings.Trim(s_block_name, "\x00")
						buffer += "<tr><td width='150' bgcolor=\"pink\"><b>B_name" + "</b></td><td width='150' bgcolor=\"pink\">" + s_block_name + "</td></tr>\n"
						s_block_inodo := string(c.B_inodo[:])
						s_block_inodo = strings.Trim(s_block_inodo, "\x00")
						buffer += "<tr><td width='150' bgcolor=\"pink\"><b>B_inodo" + "</b></td><td width='150' bgcolor=\"pink\">" + s_block_inodo + "</td></tr>\n"
					}
					buffer += "</table>>];}\n"
				default:
					fmt.Println("Tipo desconocido")
					//Salid_comando += "Tipo desconocido" + "\n"
				}

			}
			buffer += "}\n"
			Salid_comando += buffer
			file, err2 := os.Create("Bloques.dot")
			if err2 != nil && !os.IsExist(err) {
				log.Fatal(err2)
			}
			defer file.Close()

			err = os.Chmod("Bloques.dot", 0777)
			if err != nil && !os.IsExist(err) {
				log.Fatal(err)
			}

			_, err = file.WriteString(buffer)
			if err != nil && !os.IsExist(err) {
				log.Fatal(err)
			}
			cmd := exec.Command("dot", "-Tpng", "Bloques.dot", "-o", path)
			err = cmd.Run()
			if err != nil {
				fmt.Errorf("no se pudo generar la imagen: %v", err)
			}
			val_rutadis = "/home/nataly/Documentos/Mia lab/Proyecto2/MIA_P2_202001570/Backend/Discos/MIA/P2/"
			//Salid_comando += "Reporte Bloques creado exitosamente" + "\n"
		}

	}

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

func ReporteBmInode(path string, id string) {
	carpeta := ""
	archivo := ""
	regex := regexp.MustCompile(`^[a-zA-Z]+`)
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			carpeta = path[:i]
			archivo = path[i+1:]
			break
		}
	}
	fmt.Println(archivo)
	//Salid_comando += "Carpeta:" + carpeta + "\n"
	//Salid_comando += "Archivo:" + archivo + "\n"
	letters := regex.FindString(id)
	//Salid_comando += "Disco:" + letters + "\n"

	val_rutadis = val_rutadis + letters + ".dsk"

	_, err := os.Stat(carpeta)

	// Si la carpeta no existe, crearla
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(carpeta, 0755)
			if err != nil {
				Salid_comando += errors.New("Error al crear la carpeta: "+err.Error()).Error() + "\n"
				return
			}
		} else {
			Salid_comando += errors.New("Error al verificar la carpeta: "+err.Error()).Error() + "\n"
			return
		}
	}
	//Salid_comando += "La carpeta" + carpeta + "ya existe o se ha creado correctamente" + "\n"
	mbr := estructuras.Mbr{}
	sb := estructuras.Super_bloque{}
	disco, err := os.OpenFile(val_rutadis, os.O_RDWR, 0660)
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
			Mens_error(err)
		}
		disco.Seek(int64(part_starta), 0)
		err = binary.Read(disco, binary.BigEndian, &sb)
		if sb.S_filesystem_type != empty {
			s_inodos_count := string(sb.S_inodes_count[:])
			s_inodos_count = strings.Trim(s_inodos_count, "\x00")
			numero_inodos, err := strconv.Atoi(s_inodos_count)
			if err != nil {
				Mens_error(err)
			}
			s_start_bit := string(sb.S_bm_inode_start[:])
			s_start_bit = strings.Trim(s_start_bit, "\x00")
			inodos, err := strconv.Atoi(s_start_bit)
			if err != nil {
				Mens_error(err)
			}
			buffer += "digraph G{\n"
			buffer += "subgraph cluster_s {\n label=\"Inodo\"\ntbl_s [shape=box, label=<\n<table border='0' cellborder='1' cellspacing='0'  width='300' height='160' >\n"
			pos := 0
			buffer += "<tr><td width='150' bgcolor=\"pink\">"
			for c := 0; c < numero_inodos; c++ {
				disco.Seek(int64(inodos), 0)
				var char byte
				err = binary.Read(disco, binary.LittleEndian, &char)
				if err != nil {
					panic(err)
				}

				val := string(char)

				if pos < 20 {
					buffer += val + ""
					pos++
				} else {
					buffer += "</td></tr>\n"
					buffer += "<tr><td width='150' bgcolor=\"pink\">" + val + " "
					pos = 1
				}

				inodos++

			}
			buffer += "</td></tr>\n"
			buffer += "</table>>];}\n"
			buffer += "}\n"
			Salid_comando += buffer
			file, err2 := os.Create("Inodobt.dot")
			if err2 != nil && !os.IsExist(err) {
				log.Fatal(err2)
			}
			defer file.Close()

			err = os.Chmod("Inodobt.dot", 0777)
			if err != nil && !os.IsExist(err) {
				log.Fatal(err)
			}

			_, err = file.WriteString(buffer)
			if err != nil && !os.IsExist(err) {
				log.Fatal(err)
			}
			cmd := exec.Command("dot", "-Tpng", "Inodobt.dot", "-o", path)
			err = cmd.Run()
			if err != nil {
				fmt.Errorf("no se pudo generar la imagen: %v", err)
			}
			val_rutadis = "/home/nataly/Documentos/Mia lab/Proyecto2/MIA_P2_202001570/Backend/Discos/MIA/P2/"
			//Salid_comando += "Reporte de Bitmap Inodo creado exitosamente" + "\n"
		}
	}
	defer func() {
		disco.Close()
	}()
}

func ReporteBmBlock(path string, id string) {
	carpeta := ""
	archivo := ""
	regex := regexp.MustCompile(`^[a-zA-Z]+`)
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			carpeta = path[:i]
			archivo = path[i+1:]
			break
		}
	}
	fmt.Println(archivo)
	//Salid_comando += "Carpeta:" + carpeta + "\n"
	//Salid_comando += "Archivo:" + archivo + "\n"
	letters := regex.FindString(id)
	//Salid_comando += "Disco:" + letters + "\n"

	val_rutadis = val_rutadis + letters + ".dsk"

	_, err := os.Stat(carpeta)

	// Si la carpeta no existe, crearla
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(carpeta, 0755)
			if err != nil {
				Salid_comando += errors.New("Error al crear la carpeta: "+err.Error()).Error() + "\n"
				return
			}
		} else {
			Salid_comando += errors.New("Error al verificar la carpeta: "+err.Error()).Error() + "\n"
			return
		}
	}
	//Salid_comando += "La carpeta" + carpeta + "ya existe o se ha creado correctamente" + "\n"
	mbr := estructuras.Mbr{}
	sb := estructuras.Super_bloque{}
	disco, err := os.OpenFile(val_rutadis, os.O_RDWR, 0660)
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
			Mens_error(err)
		}
		disco.Seek(int64(part_starta), 0)
		err = binary.Read(disco, binary.BigEndian, &sb)
		if sb.S_filesystem_type != empty {
			s_blocks_count := string(sb.S_blocks_count[:])
			s_blocks_count = strings.Trim(s_blocks_count, "\x00")
			numero_blocks, err := strconv.Atoi(s_blocks_count)
			if err != nil {
				Mens_error(err)
			}
			s_start_blocks := string(sb.S_bm_block_start[:])
			s_start_blocks = strings.Trim(s_start_blocks, "\x00")
			blocks, err := strconv.Atoi(s_start_blocks)
			if err != nil {
				Mens_error(err)
			}
			buffer += "digraph G{\n"
			buffer += "subgraph cluster_s {\n label=\"Blocks\"\ntbl_s [shape=box, label=<\n<table border='0' cellborder='1' cellspacing='0'  width='300' height='160' >\n"
			pos := 0
			buffer += "<tr><td width='150' bgcolor=\"pink\">"
			for c := 0; c < numero_blocks; c++ {
				disco.Seek(int64(blocks), 0)
				var char byte
				err = binary.Read(disco, binary.LittleEndian, &char)
				if err != nil {
					panic(err)
				}

				val := string(char)

				if pos < 20 {
					buffer += val + ""
					pos++
				} else {
					buffer += "</td></tr>\n"
					buffer += "<tr><td width='150' bgcolor=\"pink\">" + val + " "
					pos = 1
				}

				blocks++

			}
			buffer += "</td></tr>\n"
			buffer += "</table>>];}\n"
			buffer += "}\n"
			Salid_comando += buffer
			file, err2 := os.Create("Bloquesbt.dot")
			if err2 != nil && !os.IsExist(err) {
				log.Fatal(err2)
			}
			defer file.Close()

			err = os.Chmod("Bloquesbt.dot", 0777)
			if err != nil && !os.IsExist(err) {
				log.Fatal(err)
			}

			_, err = file.WriteString(buffer)
			if err != nil && !os.IsExist(err) {
				log.Fatal(err)
			}
			cmd := exec.Command("dot", "-Tpng", "Bloquesbt.dot", "-o", path)
			err = cmd.Run()
			if err != nil {
				fmt.Errorf("no se pudo generar la imagen: %v", err)
			}
			val_rutadis = "/home/nataly/Documentos/Mia lab/Proyecto2/MIA_P2_202001570/Backend/Discos/MIA/P2/"
			//Salid_comando += "Reporte de Bitmap Bloques creado exitosamente" + "\n"
		}
	}
	defer func() {
		disco.Close()
	}()
}

func ReporteTree(path string, id string) {
	carpeta := ""
	archivo := ""
	regex := regexp.MustCompile(`^[a-zA-Z]+`)
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			carpeta = path[:i]
			archivo = path[i+1:]
			break
		}
	}
	fmt.Println(archivo)
	//Salid_comando += "Carpeta:" + carpeta + "\n"
	//Salid_comando += "Archivo:" + archivo + "\n"
	letters := regex.FindString(id)
	//Salid_comando += "Disco:" + letters + "\n"

	val_rutadis = val_rutadis + letters + ".dsk"

	_, err := os.Stat(carpeta)

	// Si la carpeta no existe, crearla
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(carpeta, 0755)
			if err != nil {
				Salid_comando += errors.New("Error al crear la carpeta: "+err.Error()).Error() + "\n"
				return
			}
		} else {
			Salid_comando += errors.New("Error al verficar la carpeta: "+err.Error()).Error() + "\n"
			return
		}
	}
	//Salid_comando += "La carpeta" + carpeta + "ya existe o se ha creado correctamente" + "\n"
	mbr := estructuras.Mbr{}
	sb := estructuras.Super_bloque{}
	disco, err := os.OpenFile(val_rutadis, os.O_RDWR, 0660)
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
			Mens_error(err)
		}
		disco.Seek(int64(part_starta), 0)
		err = binary.Read(disco, binary.BigEndian, &sb)
		if sb.S_filesystem_type != empty {
			buffer += "digraph grafica{\nrankdir=TB;\nnode [shape = record, style=filled, fillcolor=seashell2];\n"
			buffer += "subgraph cluster_s {\n label=\"\"\ntbl_s [shape=box, label=<\n"
			s_inodo_start := string(sb.S_inode_start[:])
			s_inodo_start = strings.Trim(s_inodo_start, "\x00")
			part_starta, err := strconv.Atoi(s_inodo_start)
			if err != nil {
				Mens_error(err)
			}
			s_block_start := string(sb.S_block_start[:])
			s_block_start = strings.Trim(s_block_start, "\x00")
			/*part_blocksta, err := strconv.Atoi(s_block_start)
			if err != nil {
				Mens_error(err)
			}*/
			inodo := leerInodos(int64(part_starta), val_rutadis)
			etiquetaInicio := "<TABLE>"
			etiquetaFinal := "</TABLE>"

			// Crear el primer nodo del inodo en graphviz
			celdas := ""
			celdas += "<tr> <td colspan=\"2\" bgcolor=\"yellow\">Inodo 1</td> </tr>"
			pos := 1
			for _, c := range inodo.I_block {
				celdas += fmt.Sprintf(`<tr>
					<td>apt: %d</td>
					<td>%d</td>
				</tr>`, pos, c)
				pos++
			}
			buffer += etiquetaInicio + celdas + etiquetaFinal
			buffer += ">];}\n"
			numBloque := 0
			for c, _ := range inodo.I_block {
				if c != -1 {
					numBloque = c - 1
					s_block_siz := string(sb.S_block_size[:])
					s_block_siz = strings.Trim(s_block_siz, "\x00")
					block_size, err := strconv.Atoi(s_block_siz)
					if err != nil {
						Mens_error(err)
					}
					inicioBloques := numBloque * block_size
					blocks, tipo := leerBloques(int64(inicioBloques), val_rutadis)
					fmt.Println(tipo)
					switch b := blocks.(type) {
					case estructuras.Bloque_archivo:
						fmt.Println("Bloque Archivo")
						//Salid_comando += "bloque archivo" + "\n"
					case estructuras.Bloque_carpeta:
						buffer += "subgraph cluster_" + strconv.Itoa(c) + "{\n label=\"Bloque Carpeta" + strconv.Itoa(c) + "\"\ntbl_" + strconv.Itoa(c) + "[shape=box, label=<\n"
						celdas := "<TR><TD bgcolor=\"pink\">Bloque 1</TD><TD bgcolor=\"pink\"></TD></TR>"
						for _, c := range b.B_content {
							s_block_name := string(c.B_name[:])
							s_block_name = strings.Trim(s_block_name, "\x00")
							s_block_inodo := string(c.B_inodo[:])
							s_block_inodo = strings.Trim(s_block_inodo, "\x00")
							celdas += "<TR><TD >" + s_block_name + "</TD><TD >" + s_block_inodo + "</TD></TR>"
						}
						buffer += "<" + etiquetaInicio + celdas + etiquetaFinal + ">"
						buffer += ">];}\n"
					default:
						fmt.Println("Tipo desconocido")
						//Salid_comando += "Tipo desconocido" + "\n"
					}
				}
			}

			buffer += "}\n"
			Salid_comando += buffer
			file, err2 := os.Create("Tree.dot")
			if err2 != nil && !os.IsExist(err) {
				log.Fatal(err2)
			}
			defer file.Close()

			err = os.Chmod("Tree.dot", 0777)
			if err != nil && !os.IsExist(err) {
				log.Fatal(err)
			}

			_, err = file.WriteString(buffer)
			if err != nil && !os.IsExist(err) {
				log.Fatal(err)
			}
			cmd := exec.Command("dot", "-Tpng", "Tree.dot", "-o", path)
			err = cmd.Run()
			if err != nil {
				fmt.Errorf("no se pudo generar la imagen: %v", err)
			}

			val_rutadis = "/home/nataly/Documentos/Mia lab/Proyecto2/MIA_P2_202001570/Backend/Discos/MIA/P2/"
			//Salid_comando += "Reporte de Arbol creado exitosamente" + "\n"
		}
	}
	defer func() {
		disco.Close()
	}()
}
