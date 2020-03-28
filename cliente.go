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

const BUFFER_SIZE = 1024

func enviarArchivo(fileName string) {

}

func handleMensajes(c net.Conn) {
	var msg string
	err := gob.NewDecoder(c).Decode(&msg)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(msg)
	}
}

func menu(c net.Conn) {
	go handleMensajes(c)
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
			err := gob.NewEncoder(c).Encode(&opc)
			if err != nil {
				fmt.Println(err)
			}
			var msg string
			fmt.Print("\033[12;1H Tu:")
			msg = leerString()
			err = gob.NewEncoder(c).Encode(&msg)
		case 2:
			fmt.Print("\033[12;1H Dirección del archivo:")
			err := gob.NewEncoder(c).Encode(&opc)
			if err != nil {
				fmt.Println(err)
			}
		case 3:
			fmt.Print("\033[11;1H Mensajes.")
			err := gob.NewEncoder(c).Encode(&opc)
			if err != nil {
				fmt.Println(err)
			}
		case 0:
			fmt.Println("+-------------------------------------+")
			fmt.Println("| Gracias por usar el software.       |")
			fmt.Println("+-------------------------------------+")
			err := gob.NewEncoder(c).Encode(&opc)
			if err != nil {
				fmt.Println(err)
			}
		default:
			fmt.Println("+-------------------------------------+")
			fmt.Println("| Opción inválida.                    |")
			fmt.Println("+-------------------------------------+")
		}
	}
}

func leerString() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func main() {
	var nickname string
	// Conexión inicial entre cliente servidor
	c, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	// Crea el bucle para envíar el Nickname del usuario
	for {
		fmt.Print("Nickname: ")
		nickname = leerString()
		fmt.Print("Conectando con el servidor... ")
		err = gob.NewEncoder(c).Encode(&nickname)
		if err != nil {
			fmt.Println(err)
			return
		}
		var msg string
		err = gob.NewDecoder(c).Decode(&msg)
		if msg == "Error" {
			fmt.Print("Nickname en uso, intente con uno nuevo.\n\n")
		} else {
			break
		}
	}
	menu(c)
}
