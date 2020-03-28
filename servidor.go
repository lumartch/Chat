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

func serverArchivo() {

}

func handleArchivo(c net.Conn) {

}

func existeNickname(listaUsuarios *list.List, nickname string) bool {
	var existe bool = false
	for e := listaUsuarios.Front(); e != nil; e = e.Next() {
		if e.Value == nickname {
			existe = true
			break
		}
	}
	return existe
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
			fmt.Println(nickname, ": ", msg)
		case 2:
			fmt.Println("Archivo")
		case 3:
			fmt.Println("Reespaldo de mensajes.")
		case 0:
			fmt.Println("Desconexión...")
		}
	}
	c.Close()
}

// Función de servidor que estará escuchando para cuando se conecte un Cliente en el puerto :9999
func server(listaUsuarios *list.List) {
	s, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		c, err := s.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		///
		var nickname string
		err = gob.NewDecoder(c).Decode(&nickname)
		// Verifica si es un usuario nuevo, en caso de serlo y existir en el servidor manda error
		if !existeNickname(listaUsuarios, nickname) {
			listaUsuarios.PushBack(nickname)
			fmt.Println("Se conectó: ", nickname)
			var msg string = "Conexion"
			err = gob.NewEncoder(c).Encode(&msg)
		} else {
			var msg string = "Error"
			err = gob.NewEncoder(c).Encode(&msg)
			err = nil
		}
		// Verifica que el paquete recibido no tenga errores
		if err != nil {
			fmt.Println(err)
			continue
		}
		// Crea una instancia hilo para el correspondiente usuario
		go handleUsuario(c, nickname)
	}
}

func main() {
	// Inicialización de la lista
	var listaUsuarios list.List
	//
	go server(&listaUsuarios)
	// Condicionante de paro
	var input string
	fmt.Scanln(&input)
}
