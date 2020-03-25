// Luis Fernando Martínez Castellanos
// 216787787
// Sistemas concurrentes y distribuidos
// L - M  7AM

package main

import (
	"container/list"
	"fmt"
	"net"
)

// Struct para los procesos
type Usuario struct {
	Nickname string
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
		if err != nil || listaUsuarios.Front() == nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(c)
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
