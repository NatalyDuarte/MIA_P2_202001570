import './App.css';
import axios from 'axios';
import React, { useState, useEffect } from 'react';


function App() {
  


  const [data,setdata]= useState([]);
  useEffect(() => {
  }, [])
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

    var comando = document.getElementById('parametros').value;

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
            document.getElementById('parametros').value = '';
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
                  <input type="text" id="parametros" placeholder="Ingrese comandos"></input>
                  <button class="btn btn-primary rounded-pill" onClick={Mandar}>Analizar</button>
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
                            <input type="password" id="part" class="form-control" />
                            <label class="form-label" for="form2Example2">Particion</label>
                        </div>
                        <button id="log" type="button" class="btn btn-primary btn-block mb-4" onClick={Login}>Inicio</button>
                </div>
            </div>
        </div>
      </section>

      <section id="pantalla3">
        <div id="entra-y-salida">
          <h2>Archivos</h2>
          
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
      </section>
    </div>


    </div>
  );
}

export default App;