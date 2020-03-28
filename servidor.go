// Luis Fernando Martínez Castellanos
// 216787787
// Sistemas concurrentes y distribuidos
// L - M  7AM

package main

import (
	"container/list"
	"encoding/gob"
	"fmt"
	"net"
)

const BUFFER_SIZE = 1024

// Inicialización de la lista
var listaUsuarios list.List

func serverArchivo() {

}

func handleArchivo(c net.Conn) {

}

func existeNickname(nickname string) bool {
	var existe bool = false
	for e := listaUsuarios.Front(); e != nil; e = e.Next() {
		if e.Value == nickname {
			existe = true
			break
		}
	}
	return existe
}

func eliminarNickname(nickname string) {
	for e := listaUsuarios.Front(); e != nil; e = e.Next() {
		if e.Value == nickname {
			listaUsuarios.Remove(e)
			break
		}
	}
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
		case 1:
			var msg string
			err = gob.NewDecoder(c).Decode(&msg)
			fmt.Println(nickname, ":", msg)
		case 2:
			fmt.Println("Archivo")
		case 3:
			fmt.Println("Reespaldo de mensajes.")
		case 0:
			fmt.Println(nickname, " se desconectó.")
		}
	}
	eliminarNickname(nickname)
	c.Close()
}

// Función de servidor que estará escuchando para cuando se conecte un Cliente en el puerto :9999
func server() {
	s, err := net.Listen("tcp", ":8080")
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
		// Captura el nickname ingresado por el usuario
		var nickname string
		err = gob.NewDecoder(c).Decode(&nickname)
		// Verifica que el paquete recibido no tenga errores
		if err != nil {
			fmt.Println(err)
			continue
		}
		// Verifica si es un usuario nuevo, en caso de serlo y existir en el servidor manda error
		if !existeNickname(nickname) {
			// Ingresa el nickname a la lista
			listaUsuarios.PushBack(nickname)
			fmt.Println("Se conectó: ", nickname)
			var msg string = "Conexion"
			err = gob.NewEncoder(c).Encode(&msg)
			// Crea una instancia hilo para el correspondiente usuario
			go handleUsuario(c, nickname)
		} else {
			var msg string = "Error"
			err = gob.NewEncoder(c).Encode(&msg)
			c.Close()
		}
	}
}

func main() {
	// Lanza hilo para el servidor
	go server()
	// Condicionante de paro
	var input string
	fmt.Scanln(&input)
}
