package comandos

func Logout(arre_coman []string) {
	Salid_comando += "===================LOGOUT====================" + "\n"

	band_error := false

	if !band_error {
		if iniciose != "" {
			iniciose = ""
			Salid_comando += "Se cerro sesion correctamente" + "\n"
		} else {
			Salid_comando += "Error: No se encontro una sesion iniciada" + "\n"
		}
	} else {
		Salid_comando += "Error: No se pudo cerrar sesion" + "\n"
	}
}
