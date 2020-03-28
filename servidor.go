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

// Struct para los procesos
type Usuario struct {
	Nickname string
	Opcion   int64
	Mensaje  string
}

func serverArchivo() {

}

func handleArchivo(c net.Conn) {
}

func existeNickname(listaUsuarios *list.List, nickname string) bool {
	var existe bool = false
	for e := listaUsuarios.Front(); e != nil; e = e.Next() {
		if e.Value.(Usuario).Nickname == nickname {
			existe = true
			break
		}
	}
	return existe
}

func handlerUsuarios() {

}

// Función de servidor que estará escuchando para cuando se conecte un Cliente en el puerto :9999
func server(listaUsuarios *list.List) {
	s, err := net.Listen("tcp", ":9999")
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
		var usr Usuario
		err = gob.NewDecoder(c).Decode(&usr)
		// Verifica si es un usuario nuevo, en caso de serlo y existir en el servidor manda error
		if usr.Opcion == 0 && !existeNickname(listaUsuarios, usr.Nickname) {
			listaUsuarios.PushBack(usr)
			fmt.Println("Se conectó: ", usr.Nickname)
			var msg string = "Conexion"
			err = gob.NewEncoder(c).Encode(&msg)
		} else {
			var msg string = "Error"
			err = gob.NewEncoder(c).Encode(&msg)
			err = nil
		}
		//
		if err != nil {
			fmt.Println(err)
			continue
		}
		/// Switch para los handlers de acciones
		switch usr.Opcion {
		case 1:
			fmt.Println("Mensaje")
		case 2:
			fmt.Println("Archivo")
		case 3:
			fmt.Println("Reespaldo de mensajes.")
		case 4:
			fmt.Println("Desconexión...")
		}
		c.Close()
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
