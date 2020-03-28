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

func handleMensajesNuevos(c net.Conn) {
	var i int64 = 0
	for {
		var msg string
		err := gob.NewDecoder(c).Decode(&msg)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Print("\033[", (4 + i), ";32H", msg)
		i += 1
		fmt.Print("\033[9;5H")
	}
}

func imprimirTodoChat(c net.Conn) {
	var msgs []string
	err := gob.NewDecoder(c).Decode(&msgs)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i, msg := range msgs {
		fmt.Print("\033[", (4 + i), ";32H", msg)
	}
}

func menu(c net.Conn, nickname string) {
	fmt.Print("\033[H\033[2J")
	// Hilo para recibir mensajes del servidor o de otros usuarios
	go handleMensajesNuevos(c)
	// Ciclo para las opciones de Usuario
	var opc int64 = 1
	for opc != 0 {
		// System "Clear"
		fmt.Print("\033[H\033[2J")
		fmt.Println("+--------------------------+-----------------------------------------------------------------+")
		fmt.Println("|         WutsClient       |                       Chat                                      |")
		fmt.Println("+--------------------------+-----------------------------------------------------------------+")
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
			fmt.Print("\033[12;1HTu: ")
			msg = leerString()
			err = gob.NewEncoder(c).Encode(&msg)
			//imprimirTodoChat(c)
		case 2:
			err := gob.NewEncoder(c).Encode(&opc)
			if err != nil {
				fmt.Println(err)
			}
			/*var file string
			fmt.Print("\033[12;1H Dirección del archivo: ")
			file = leerString()*/
		case 3:
			fmt.Print("\033[11;1H Mensajes.")
			err := gob.NewEncoder(c).Encode(&opc)
			if err != nil {
				fmt.Println(err)
			}
		case 0:
			fmt.Println("+--------------------------+")
			fmt.Println("| Vuelva pronto.   :B      |")
			fmt.Println("+--------------------------+")
			err := gob.NewEncoder(c).Encode(&opc)
			if err != nil {
				fmt.Println(err)
			}
		default:
			opc = 20
		}
	}
}

func leerString() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func main() {
	// Crea el bucle para envíar el Nickname del usuario
	for {
		var nickname string
		// Conexión inicial entre cliente servidor
		c, err := net.Dial("tcp", "192.168.100.4:8080")
		if err != nil {
			fmt.Println(err)
			return
		}
		// El usuario ingresa el Nickname
		fmt.Print("Nickname: ")
		nickname = leerString()
		fmt.Print("Conectando con el servidor... ")
		// Se envía el nickname
		err = gob.NewEncoder(c).Encode(&nickname)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// Se verifica que el nickname no esté en uso
		var msg string
		err = gob.NewDecoder(c).Decode(&msg)
		//
		if err != nil {
			fmt.Println(err)
			continue
		}
		if msg == "Error" {
			fmt.Print("Nickname en uso, intente con uno nuevo.\n\n")
		} else {
			menu(c, nickname)
			break
		}
	}
}
