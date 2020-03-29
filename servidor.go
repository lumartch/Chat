// Luis Fernando Martínez Castellanos
// 216787787
// Sistemas concurrentes y distribuidos
// L - M  7AM

package main

import (
	"container/list"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

const BUFFER_SIZE = 1024

type Usuario struct {
	Nickname string
	Conexion net.Conn
}

type Archivo struct {
	Nombre string
	Datos  []byte
}

// Inicialización de la lista
var listaUsuarios list.List
var msgs []string

func handleArchivo(arc Archivo) {
	// Crea una copia del archivo envíado dentro de la carpeta de cada usuario conectado
	for e := listaUsuarios.Front(); e != nil; e = e.Next() {
		err := gob.NewEncoder(e.Value.(Usuario).Conexion).Encode(&arc)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// Función para corroborar que el usuario existe
func existeNickname(nickname string) bool {
	var existe bool = false
	for e := listaUsuarios.Front(); e != nil; e = e.Next() {
		if e.Value.(Usuario).Nickname == nickname {
			existe = true
			break
		}
	}
	return existe
}

// Función para eliminar de la lista al usuario
func eliminarNickname(nickname string) {
	for e := listaUsuarios.Front(); e != nil; e = e.Next() {
		if e.Value.(Usuario).Nickname == nickname {
			listaUsuarios.Remove(e)
			break
		}
	}
}

// Handle para envíar el mensaje a todas las conexiones dentro del servidor
func handleMensajes(msg string) {
	for e := listaUsuarios.Front(); e != nil; e = e.Next() {
		err := gob.NewEncoder(e.Value.(Usuario).Conexion).Encode(&msg)
		if err != nil {
			fmt.Println(err)
		}
	}
	msgs = append(msgs, msg)
}

func handleUsuario(c net.Conn, nickname string) {
	var opc int = 1
	for opc != 0 {
		err := gob.NewDecoder(c).Decode(&opc)
		if err != nil {
			fmt.Println(err)
		}
		/// Switch para los handlers de acciones
		switch opc {
		// Si se captura un 1 envía un "echo" del mensaje al resto de usuarios
		case 1:
			var msg string
			err = gob.NewDecoder(c).Decode(&msg)
			fmt.Println(nickname, ":", msg)
			handleMensajes(nickname + ": " + msg)
		// Si se captura un 2 reenvía el archivo al resto de usuarios
		case 2:
			var arc Archivo
			err := gob.NewDecoder(c).Decode(&arc)
			if err != nil {
				fmt.Println(err)
			}
			err = ioutil.WriteFile("files/"+arc.Nombre, arc.Datos, 0644)
			if err != nil {
				fmt.Println("Error creating", arc.Nombre)
				fmt.Println(err)
				return
			}
			handleMensajes(nickname + " envío: " + arc.Nombre)
			//handleArchivo(arc)
		// Si captura un se termina la conexión con el usuario
		case 0:
			fmt.Println(nickname, " se desconectó.")
		}
	}
	c.Close()
	eliminarNickname(nickname)
	handleMensajes(nickname + " se desconectó.")
}

func handleConexion(c net.Conn) {
	// Captura el nickname ingresado por el usuario
	var nickname string
	err := gob.NewDecoder(c).Decode(&nickname)
	// Verifica que el paquete recibido no tenga errores
	if err != nil {
		fmt.Println(err)
	}
	// Verifica si es un usuario nuevo, en caso de serlo y existir en el servidor manda error
	if !existeNickname(nickname) {
		// Ingresa el usuario a la lista y se imprime dentro del servidor
		listaUsuarios.PushBack(Usuario{Nickname: nickname, Conexion: c})
		fmt.Println("Se conectó: ", nickname)
		// Se manda la notificación a los usuarios conectados actualmente
		handleMensajes(nickname + " se conectó.")
		// Se notifica a la conexión del usuario actual que se conectó sin errores
		var msg string = "Conexion"
		err = gob.NewEncoder(c).Encode(&msg)
		// Crea una instancia hilo para el correspondiente usuario
		go handleUsuario(c, nickname)
	} else {
		// Se manda un mensaje de error al usuario y se termina la conexión
		var msg string = "Error"
		err = gob.NewEncoder(c).Encode(&msg)
		c.Close()
	}
}

// Función de servidor que estará escuchando para cuando se conecte un Cliente en el puerto :9999
func server() {
	s, err := net.Listen("tcp", "192.168.100.4:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		// Acepta el request del usuario
		c, err := s.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		// Crea un hilo para la conexión actual y espera una nueva
		go handleConexion(c)
	}
}

func main() {
	// Crea un directorio para el servidor donde almacena todos los archivos
	_, err := os.Stat("files")
	if os.IsNotExist(err) {
		errDir := os.MkdirAll("files", 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}
	// Lanza hilo para el servidor
	go server()
	// Condicionante de paro
	var input string
	fmt.Scanln(&input)
}
