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

// Struct para los procesos
type Usuario struct {
	Nickname string
	Opcion   int64
	Mensaje  string
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
		if err != nil {
			fmt.Println(err)
			continue
		}
		/// Switch para los handlers de acciones
		switch usr.Opcion {
		case 0:
			fmt.Println("Usuario nuevo")
		case 1:
			fmt.Println("Mensaje")
		case 2:
			fmt.Println("Archivo")
		case 3:
			fmt.Println("Reespaldo de mensajes.")
		}
		c.Close()
		// Envía los mensajes a los demas usuarios
		//go handlerUsuarios()
		/*listaUsuarios.Remove(listaUsuarios.Front())
		go handleCliente(c, p, listaUsuarios)*/
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
