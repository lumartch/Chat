// Luis Fernando Martínez Castellanos
// 216787787
// Sistemas concurrentes y distribuidos
// L - M  7AM

package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

type Archivo struct {
	Nombre string
	Datos  []byte
}

type Mensaje struct {
	Opcion  int64
	Mensaje string
	File    Archivo
}

// Se envía archivo al servidor
func enviarArchivo(c net.Conn, dirArchivo string, nickname string) {
	// Abre el archivo si existe en el directorio
	f, err := os.Open(dirArchivo)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()
	// Obtiene la información del archivo
	fileStat, err := os.Stat(dirArchivo)
	if err != nil {
		log.Fatal(err)
		return
	}
	// Lee el archivo desde el origen
	input, err := ioutil.ReadFile(dirArchivo)
	err = gob.NewEncoder(c).Encode(&Archivo{Nombre: fileStat.Name(), Datos: input})
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Cada mensaje que es recibido del servidor se guarda dentro de un Slice
func handleRespuestaServidor(c net.Conn, nickname string, msgs ...string) {
	for {
		var msg Mensaje
		err := gob.NewDecoder(c).Decode(&msg)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if msg.Opcion == 1 {
			msgs = append(msgs, msg.Mensaje)
			imprimirMensajes(msgs...)
		} else {
			fileName := msg.File.Nombre
			err = ioutil.WriteFile(nickname+"/"+fileName, msg.File.Datos, 0644)
			if err != nil {
				fmt.Println("Error creating", fileName)
				fmt.Println(err)
				return
			}
		}
	}
}

// Limpia la ventana de chat
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
	// Se crea el historial de conversación del usuario
	msgs := []string{}
	// Hilo para recibir mensajes del servidor o de otros usuarios
	go handleRespuestaServidor(c, nickname, msgs...)
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
			var dirArchivo string
			fmt.Print("\033[12;1H Dirección del archivo: ")
			dirArchivo = leerString()
			enviarArchivo(c, dirArchivo, nickname)
			opc = 20
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
			continue
		}
		menu(c, nickname)
		break
	}
}
