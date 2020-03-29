// Luis Fernando Martínez Castellanos
// 216787787
// Sistemas concurrentes y distribuidos
// L - M  7AM

package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
)

const BUFFER_SIZE = 1024

func enviarArchivo(fileName string) {

}

// Cada mensaje que es recibido del servidor se guarda dentro de un Slice
func handleMensajesNuevos(c net.Conn, msgs ...string) {
	for {
		var msg string
		err := gob.NewDecoder(c).Decode(&msg)
		if err != nil {
			fmt.Println(err)
		}
		msgs = append(msgs, msg)
		imprimirMensajes(msgs...)
	}
}

func limpiarChat() {
	for i := 0; i < 21; i++ {
		fmt.Print("\033[", (4 + i), ";29H                                                               ")
	}
}

// Imprime todos los mensajes dentro de la interfáz de usuario
func imprimirMensajes(msgs ...string) {
	limpiarChat()
	var i int64 = 20
	for j := 0; j < len(msgs) && j < 20; j++ {
		if i >= 0 {
			if len(msgs[len(msgs)-j-1]) > 60 {
				msg := []string{}
				var str string
				for k, ch := range msgs[len(msgs)-j-1] {
					if (k+1)%60 == 0 {
						msg = append(msg, str)
						str = ""
						str = str + string(ch)
					} else {
						str = str + string(ch)
					}
				}
				msg = append(msg, str)
				for k := 0; k < len(msg); k++ {
					if i >= 0 {
						fmt.Print("\033[", (4 + i), ";32H", msg[len(msg)-k-1])
						i--
					}
				}
			} else {
				fmt.Print("\033[", (4 + i), ";32H", msgs[len(msgs)-j-1])
				i--
			}
		}
	}
	/*for _, msg := range msgs {
		if len(msg) > 60 {
			var str string
			for j, ch := range msg {
				if (j+1)%60 == 0 {
					fmt.Print("\033[", (4 + i), ";32H", str)
					str = ""
					str = str + string(ch)
					i++
				} else {
					str = str + string(ch)
				}
			}
			fmt.Print("\033[", (4 + i), ";32H", str)
			i++
		} else {
			fmt.Print("\033[", (4 + i), ";32H", msg)
			i++
		}
	}*/
	fmt.Print("\033[9;5H")
}

func imprimirInterfaz() {
	// System Clear
	fmt.Print("\033[H\033[2J")
	// Se imprime la interfaz para el usuario
	fmt.Println("+--------------------------+-----------------------------------------------------------------+")
	fmt.Println("|         WutsClient       |                       Chat                                      |")
	fmt.Println("+--------------------------+-----------------------------------------------------------------+")
	fmt.Println("| 1.- Envíar mensaje.      |")
	fmt.Println("| 2.- Enviar archivo.      |")
	fmt.Println("| 3.- Guardar chat.        |")
	fmt.Println("| 0.- Salir.               |")
	fmt.Println("+--------------------------+")
	fmt.Println("| :                        |")
	fmt.Println("+--------------------------+")
}

func menu(c net.Conn, nickname string) {
	// Crea un directorio para el cliente donde almacena sus archivos
	_, err := os.Stat(nickname)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(nickname, 0755)
		if errDir != nil {
			log.Fatal(err)
		}

	}
	msgs := []string{}
	// Hilo para recibir mensajes del servidor o de otros usuarios
	go handleMensajesNuevos(c, msgs...)
	imprimirInterfaz()
	// Ciclo para las opciones de Usuario
	var opc int64 = 1
	for opc != 0 {
		fmt.Print("\033[9;1H| :                        |")
		fmt.Print("\033[12;1H                            ")
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
			opc = 20
		case 2:
			err := gob.NewEncoder(c).Encode(&opc)
			if err != nil {
				fmt.Println(err)
			}
			/*var file string
			fmt.Print("\033[12;1H Dirección del archivo: ")
			file = leerString()*/
			opc = 20
		//case 3:

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
		var m string
		err = gob.NewDecoder(c).Decode(&m)
		//
		if err != nil {
			fmt.Println(err)
			continue
		}
		if m == "Error" {
			fmt.Print("Nickname en uso, intente con uno nuevo.\n\n")
		} else {
			menu(c, nickname)
			break
		}
	}
}
