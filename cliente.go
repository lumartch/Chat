// Luis Fernando Martínez Castellanos
// 216787787
// Sistemas concurrentes y distribuidos
// L - M  7AM

package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

// Struct para los procesos
type Usuario struct {
	Nickname string
	Opcion   int64
	Mensaje  string
}

func enviarMensaje(usr *Usuario) {
	// Conexión inicial entre cliente servidor
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = gob.NewEncoder(c).Encode(&usr)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func menu(usr *Usuario) {
	var opc int64
	opc = 1
	for opc != 0 {
		// System "Clear"
		fmt.Print("\033[H\033[2J")
		fmt.Println("+--------------------------+")
		fmt.Println("|         WutsClient       |")
		fmt.Println("+--------------------------+")
		fmt.Println("| 1.- Envíar mensaje.      |")
		fmt.Println("| 2.- Enviar archivo.      |")
		fmt.Println("| 3.- Mostrar chat.        |")
		fmt.Println("| 0.- Salir.               |")
		fmt.Println("+--------------------------+")
		fmt.Println("| :                        |")
		fmt.Println("+--------------------------+")
		fmt.Print("\033[9;5H")
		fmt.Scanln(&opc)
		switch opc {
		case 1:
			fmt.Print("\033[12;1H Tu: ")
			usr.Mensaje = leerString()
			usr.Opcion = 1
			go enviarMensaje(usr)
		case 2:
			fmt.Print("\033[12;1H Dirección del archivo:")
		case 3:
			fmt.Print("\033[11;1H Mensajes.")
		case 0:
			fmt.Println("+-------------------------------------+")
			fmt.Println("| Gracias por usar el software.       |")
			fmt.Println("+-------------------------------------+")
		default:
			fmt.Println("+-------------------------------------+")
			fmt.Println("| Opción inválida.                    |")
			fmt.Println("+-------------------------------------+")
		}
	}
}

func enviarNickname(usr *Usuario) {
	// Conexión inicial entre cliente servidor
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = gob.NewEncoder(c).Encode(&usr)
	if err != nil {
		fmt.Println(err)
		return
	}
	menu(usr)
}

func leerString() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func main() {
	var nickname string
	fmt.Print("Nickname: ")
	nickname = leerString()
	fmt.Println(nickname)
	enviarNickname(&Usuario{Nickname: nickname, Opcion: 0, Mensaje: ""})
}
