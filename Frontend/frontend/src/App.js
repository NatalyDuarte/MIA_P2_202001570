import './App.css';
import axios from 'axios';
import React, { useState, useEffect } from 'react';
import { graphviz } from "d3-graphviz"
import { useFileSystem } from 'use-file-system';
import Select from 'react-select';
import { FilePicker } from 'react-file-picker';

function App() {
  const [fileList, setFileList] = useState([]);
  const [selectedOption, setSelectedOption] = useState(null);

  const [fileLists, setFileLists] = useState([]);
  const [selectedOptions, setSelectedOptions] = useState(null);


  const [mostrarOpciones, setMostrarOpciones] = useState(false);
  const [carpetaSeleccionada, setCarpetaSeleccionada] = useState('');

  const toggleOpciones = () => {
    setMostrarOpciones(!mostrarOpciones);
  };

  const handleChangec = (event) => {
    setCarpetaSeleccionada(event.target.value);
  };

  const [mostrarOpcionesc, setMostrarOpcionesd] = useState(false);
  const [carpetaSeleccionadad, setCarpetaSeleccionadad] = useState('');

  const toggleOpcionesd = () => {
    setMostrarOpcionesd(!mostrarOpcionesc);
  };

  const handleChangear = (event) => {
    setCarpetaSeleccionadad(event.target.value);
  };


  const [selectedFile, setSelectedFile] = useState(null);

  const handleFileChange = (event) => {
    const file = event.target.files[0];
    setSelectedFile(file);
  };

  const handleUploadClick = () => {
    const fileInput = document.getElementById('fileInput');
    fileInput.click();
    fileInput.addEventListener('change', async (event) => {
      const file = event.target.files[0];
      if (!file) {
        return; // Handle empty file selection
      }
      try {
        const text = await file.text();
        console.log(text)
        var entradaActual = document.getElementById('entrada').value;
        var nuevoContenido = entradaActual + text + "\n \n" ;
        document.getElementById('entrada').value = nuevoContenido;
      } catch (error) {
        console.error('Error reading file:', error);
        // Handle file reading errors (e.g., display error message)
      }
    });
  };
  const handleChange = (selected) => {
    setSelectedOption(selected);
  };
  const handleChanges = (selected) => {
    setSelectedOptions(selected);
  };
  

  const outerStyle = {
    backgroundColor: 'black', // Using camelCase for backgroundColor
    padding: '1.5rem',
    borderRadius: '10px',
    margin: '2rem',
  };

  const innerStyle = {
    paddingLeft: '2%',
    paddingRight: '2%',
    paddingBottom: '2%',
  };

  function Parti(event){
    event.preventDefault();
    const comando = selectedOption.value

    var datos = {
        "Cmd": comando
    };

    axios.post("http://localhost:5000/verparti", datos)
    .then(
      (response) => {
        const files = response.data;
        setFileLists(files);
        
      }
    );
        
  }

  function VerCar(event){
    event.preventDefault();
    const comand = selectedOption.value
    const comandos = document.getElementById('part').value;
    const comando = comand +","+comandos
    var datos = {
        "Cmd": comando
    };

    axios.post("http://localhost:5000/vercar", datos)
    .then(
      (response) => {
        const files = response.data;
        setFileLists(files);
        
      }
    );
        
  }

  function ElectPar (event){
    event.preventDefault();
    let selectedValue= selectedOption.value
    let fileNameWithoutExtension = selectedValue.substring(0, selectedValue.lastIndexOf(".dsk"));
    let part = selectedOptions.value
    let regex = /\d+$/; // Matches one or more digits at the end of the string
    let match = part.match(regex);
    let numericPart = ""
    if (match) {
      numericPart = match[0]; // Get the first captured group (numeric part)
    } else {
      console.log("No numeric part found");
    }
    let combinedText = `${fileNameWithoutExtension}${numericPart}70`; 
    let textBox = document.getElementById("part");
    textBox.value = combinedText;

  }
  function VerD(event){
    event.preventDefault();
    const fetchFiles = async () => {
      try {
        const response = await axios.get('http://localhost:5000/files');
        const files = response.data.files;
        setFileList(files);
      } catch (error) {
        console.error('Error fetching files:', error);
      }
    };

    fetchFiles();
  }

  useEffect(() => {
  }, [])

  function Mandarss(event){
    event.preventDefault();

    var comando = document.getElementById('parametr').value;

    var datos = {
        "Cmd": comando
    };

    axios.post("http://localhost:5000/analizar", datos)
    .then(
        (response) => {
            console.log(response.data);
            var respo = response.data;
            const textoSinResultado = respo.replace(/"result":\s*"/, '');
            const variableSinLlave = textoSinResultado.replace(/^\{/, "");
            const variableSinLlaves = variableSinLlave.replace(/}$/, "");
            console.log(variableSinLlaves)
            graphviz("#graph").renderDot(String(variableSinLlaves))
            document.getElementById('parametr').value = '';
        }
    )
  }

  function Mandars(event){
    event.preventDefault();

    var comando = document.getElementById('parametross').value;

    var datos = {
        "Cmd": comando
    };

    axios.post("http://localhost:5000/analizar", datos)
    .then(
        (response) => {
            console.log(response.data);
            var respo = response.data;
            const textoSinResultado = respo.replace(/"result":\s*"/, '');
            const variableSinComillas = textoSinResultado.replace(/"/g, "");
            const variableSinLlave = variableSinComillas.replace(/\{/g, "");
            const variableSinLlaves = variableSinLlave.replace(/\}/g, "");
            document.getElementById('salidas').value = variableSinLlaves;

            // Agregar el comando al textarea de entrada
            var entradaActual = document.getElementById('entradas').value;
            var nuevoContenido = entradaActual + comando + "\n \n" ;
            document.getElementById('entradas').value = nuevoContenido;
            document.getElementById('parametross').value = '';
        }
    )
  }
  function Mandar(event){
    event.preventDefault();

    var comando = document.getElementById('entrada').value;
    
    var datos = {
        "Cmd": comando
    };


    axios.post("http://localhost:5000/analizar", datos)
    .then(
        (response) => {
            console.log(response.data);
            var respo = response.data;
            const textoSinResultado = respo.replace(/"result":\s*"/, '');
            const variableSinComillas = textoSinResultado.replace(/"/g, "");
            const variableSinLlave = variableSinComillas.replace(/\{/g, "");
            const variableSinLlaves = variableSinLlave.replace(/\}/g, "");
            document.getElementById('salida').value = variableSinLlaves;

            // Agregar el comando al textarea de entrada
            var entradaActual = document.getElementById('entrada').value;
            var nuevoContenido = entradaActual + comando + "\n \n" ;
            document.getElementById('entrada').value = nuevoContenido;
        }
    )
  }
  function Login(event) {
    event.preventDefault();
    var usuario = document.getElementById("user").value;
    var pass = document.getElementById("pass").value;
    var part = document.getElementById("part").value;
    var comando = "login -user="+usuario+" -pass="+pass+" -id="+part
    var datos = {
      "Cmd": comando
    };

    axios.post("http://localhost:5000/analizar", datos)
    .then(
        (response) => {
              console.log(response.data);
              var respo = response.data;
              const textoSinResultado = respo.replace(/"result":\s*"/, '');
              const variableSinComillas = textoSinResultado.replace(/"/g, "");
              const variableSinLlave = variableSinComillas.replace(/\{/g, "");
              const variableSinLlaves = variableSinLlave.replace(/\}/g, "");
              const sesionIniciada = "Sesion iniciada";

              const contieneSesion = variableSinLlaves.includes(sesionIniciada);
              
              if (contieneSesion) {
                alert("BIENVENIDO USUARIO")
                document.getElementById('UsuaPag').style.display = 'block';
                document.getElementById('HomePag').style.display = 'none';
                window.location.href = "#page-top";
                document.getElementById("user").value = "";
                document.getElementById("pass").value = "";
                document.getElementById("part").value = "";
                var tex = '<h1 class="masthead-heading mb-0">Bienvenido</h1>'
                tex += "<h2>" + usuario + "</h2>"
                document.getElementById("nombreusua").innerHTML = tex
              } else {
                alert("Usuario o contraseña incorrectos")
              }
        }
    )
  }
  function CerrarUsua(event) {
    event.preventDefault();
    alert("Cerrando Sesion Usuario....")
    document.getElementById("HomePag").style.display = "block";
    document.getElementById("UsuaPag").style.display = "none";
    window.location.href = "#page-top";
    var comando = "logout"
    var datos = {
      "Cmd": comando
    };

    axios.post("http://localhost:5000/analizar", datos)
    .then(
        (response) => {
              console.log(response.data);
              var respo = response.data;
              const textoSinResultado = respo.replace(/"result":\s*"/, '');
              const variableSinComillas = textoSinResultado.replace(/"/g, "");
              const variableSinLlave = variableSinComillas.replace(/\{/g, "");
              const variableSinLlaves = variableSinLlave.replace(/\}/g, "");
              const sesionIniciada = "Se cerro sesion correctamente";

              const contieneSesion = variableSinLlaves.includes(sesionIniciada);
              
              if (contieneSesion) {
                alert("Cerrando Sesion Usuario....")
              } else {
                alert("Ocurrio un error")
              }
        }
    )

  }


  return (
    <div className="App">
      <div id="HomePag" style={{ display: 'block' }}>
      <header class="masthead text-center text-white">
        <div class="masthead-content">
          <div class="container px-5">
              <h1 class="masthead-heading mb-0">Bienvenido</h1>
              <h2>Nataly Saraí Guzmán Duarte</h2>
              <h2>202001570</h2>
              <h5>Manejo e implentacion de archivos</h5>
          </div>
        </div>
        <div class="bg-circle-1 bg-circle"></div>
        <div class="bg-circle-2 bg-circle"></div>
        <div class="bg-circle-3 bg-circle"></div>
      </header>
      <nav class="navbar navbar-expand-lg navbar-dark navbar-custom fixed-top">
        <div class="container px-5">
            <a class="navbar-brand" href="#page-top">Proyecto2</a>
            <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarResponsive" aria-controls="navbarResponsive" aria-expanded="false" aria-label="Toggle navigation"><span class="navbar-toggler-icon"></span></button>
            <div class="collapse navbar-collapse" id="navbarResponsive">
                <ul class="navbar-nav ms-auto">
                    <li class="nav-item"><a class="nav-link" href="#pantalla1">Pantalla1</a></li>
                    <li class="nav-item"><a class="nav-link" href="#pantalla2">Pantalla2</a></li>
                    <li class="nav-item"><a class="nav-link" href="#pantalla3">Pantalla3</a></li>
                </ul>
            </div>
        </div>
      </nav>
      <section id="pantalla1">
        <div id="entra-y-salida">
          <div id="Entra">
            <h2>Entrada:</h2>
            <div class="background-color: black; height: 260px; width: 600px; border-radius: 1rem;">
              <textarea id="entrada" rows="20" cols="50"></textarea>
            </div>
            <div>
              <p>
                <div id="boton-y-caja">  
                  <button class="btn btn-primary rounded-pill" onClick={Mandar}>Analizar</button>
                </div>
                <div>
                    <input
                      id="fileInput"
                      type="file"
                      style={{ display: 'none' }}
                      onChange={handleFileChange}
                    />
                    <button onClick={handleUploadClick}>Cargar Archivo</button>
                    {selectedFile && <p>Archivo seleccionado: {selectedFile.name}</p>}
                  </div>
              </p>
            </div>
          </div>
          <div id="Salida">
            <h2>Salida:</h2>
            <div class="background-color: black; height: 260px; width: 600px; border-radius: 1rem;">
              <textarea id="salida" rows="20" cols="50"></textarea>
            </div>
          </div>
        </div>
      </section>
      <section id="pantalla2">
        <div class="container px-5">
            <div class="row gx-7 align-items-center">
              <div class="col-lg-6 order-lg-1">
               
                        <h2 class="display-4">Inicio Sesion</h2>
                        
                        <div class="form-outline mb-4">
                            <input type="text" id="user" class="form-control" />
                            <label class="form-label" for="form2Example1">User</label>
                        </div>
                        
                        <div class="form-outline mb-4">
                            <input type="password" id="pass" class="form-control" />
                            <label class="form-label" for="form2Example2">Password</label>
                        </div>

                        <div class="form-outline mb-4">
                            <input type="text" id="part" class="form-control" disabled value=""  />
                            <label class="form-label" for="form2Example2">Particion</label>
                            <div>
                            <button id="log" type="button" class="btn btn-primary btn-block mb-4" onClick={VerD}>VerDiscos</button>
                            <Select
                                id="fileSelect"
                                value={selectedOption} 
                                onChange={handleChange} 
                                options={fileList.map((fileName) => ({ value: fileName, label: fileName }))} 
                              />
                              <p>Opción seleccionada Disco: {selectedOption ? selectedOption.value : 'Ninguna'}</p>
                              <button id="log" type="button" class="btn btn-primary btn-block mb-4" onClick={Parti}>Elegir Particion</button>
                              <Select
                                id="fileSelectP"
                                value={selectedOptions} 
                                onChange={handleChanges} 
                                options={fileLists.map((fileName) => ({ value: fileName, label: fileName }))} 
                              />
                              <p>Opción seleccionada Particion: {selectedOptions ? selectedOptions.value : 'Ninguna'}</p>
                              <button id="log" type="button" class="btn btn-primary btn-block mb-4" onClick={ElectPar}>Confirmar</button>
                            </div>
                        </div>
                        <button id="log" type="button" class="btn btn-primary btn-block mb-4" onClick={Login}>Inicio</button>
                </div>
            </div>
        </div>
      </section>

      <section id="pantalla3">
          <h2>Visualizacion de Reportes</h2>
          <div id="boton-y-caja">  
                  <input type="text" id="parametr" placeholder="Ingrese comandos"></input>
                  <button class="btn btn-primary rounded-pill" onClick={Mandarss}>Analizar</button>
          </div>
          
          <div className="bg-dark p-4 rounded m-5" style={outerStyle}>
            <div style={innerStyle}>
              <h4 style={{ color: 'white', marginBottom: '1.5%', marginTop: '0.5%' }}>Reporte:</h4>
              <div id="graph"></div>
            </div>
          </div>
      </section>
      </div>

      <div id="UsuaPag" style={{ display: 'none' }}>
        <nav class="navbar navbar-expand-lg navbar-dark navbar-custom fixed-top">
            <div class="container px-5">
                <a class="navbar-brand" href="#page-top">GoDrive</a>
                <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarResponsive" aria-controls="navbarResponsive" aria-expanded="false" aria-label="Toggle navigation"><span class="navbar-toggler-icon"></span></button>
                <div class="collapse navbar-collapse" id="navbarResponsive">
                    <ul class="navbar-nav ms-auto">
                        <li class="nav-item">
                            <a class="btn btn-primary  rounded-pill " onClick={CerrarUsua}>CERRAR SESION</a>
                        </li>
                    </ul>
                </div>
            </div>
        </nav>
        <header class="masthead text-center text-white">
            <div class="masthead-content">
                <div id="nombreusua" class="container px-5">
                </div>
            </div>
            <div class="bg-circle-1 bg-circle"></div>
            <div class="bg-circle-2 bg-circle"></div>
            <div class="bg-circle-3 bg-circle"></div>
        </header>
        <section id="pantalla1">
        <div id="entra-y-salida">
          <div id="Entra">
            <h2>Entrada:</h2>
            <div class="background-color: black; height: 260px; width: 600px; border-radius: 1rem;">
              <textarea id="entradas" rows="20" cols="50"></textarea>
            </div>
            <div>
              <p>
                <div id="boton-y-caja">  
                  <input type="text" id="parametross" placeholder="Ingrese comandos"></input>
                  <button class="btn btn-primary rounded-pill" onClick={Mandars}>Analizar</button>
                </div>
              </p>
            </div>
          </div>
          <div id="Salida">
            <h2>Salida:</h2>
            <div class="background-color: black; height: 260px; width: 600px; border-radius: 1rem;">
              <textarea id="salidas" rows="20" cols="50"></textarea>
            </div>
          </div>
          
         
        </div>
        <div>
      <button onClick={toggleOpciones}>
        {mostrarOpciones ? 'Ocultar opciones' : 'Mostrar Carpetas'}
      </button>
      {mostrarOpciones && (
        <select value={carpetaSeleccionada} onChange={handleChangec}>
          <option value="carpeta0">Carpeta 0</option>
          <option value="carpeta1">Carpeta 1</option>
        </select>
      )}
    </div>
    <div>
      <button onClick={toggleOpcionesd}>
        {mostrarOpcionesc ? 'Ocultar opciones' : 'Mostrar Archivos'}
      </button>
      {mostrarOpcionesc && (
        <select value={carpetaSeleccionadad} onChange={handleChangear}>
          <option value="carpeta0">user.txt</option>
        </select>
      )}
    </div>
      </section>
    </div>


    </div>
  );
}

export default App;