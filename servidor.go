// Luis Fernando Martínez Castellanos
// 216787787
// Sistemas concurrentes y distribuidos
// L - M  7AM

package main

import (
	"container/list"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

type Usuario struct {
	Nickname string
	Conexion net.Conn
}

type Archivo struct {
	Nombre string
	Datos  []byte
}

type Mensaje struct {
	Opcion  int64
	Mensaje string
	File    Archivo
}

// Inicialización de la lista
var lUsrMensajes list.List
var msgs []string
var files []string

// Función para corroborar que el usuario existe
func existeNickname(nickname string) bool {
	var existe bool = false
	for e := lUsrMensajes.Front(); e != nil; e = e.Next() {
		if e.Value.(Usuario).Nickname == nickname {
			existe = true
			break
		}
	}
	return existe
}

// Función para eliminar de la lista al usuario
func eliminarNickname(nickname string) {
	for e := lUsrMensajes.Front(); e != nil; e = e.Next() {
		if e.Value.(Usuario).Nickname == nickname {
			lUsrMensajes.Remove(e)
			break
		}
	}
}

// Handle para envíar el mensaje a todas las conexiones dentro del servidor
func enviarMensaje(msg string) {
	for e := lUsrMensajes.Front(); e != nil; e = e.Next() {
		// Envía el mensaje a los usuarios
		err := gob.NewEncoder(e.Value.(Usuario).Conexion).Encode(&Mensaje{Opcion: 1, Mensaje: msg})
		if err != nil {
			fmt.Println(err)
		}
	}
	msgs = append(msgs, msg)
}

//
func handleArchivo(fileName string) {
	// Crea una copia del archivo envíado dentro de la carpeta de cada usuario conectado
	for e := lUsrMensajes.Front(); e != nil; e = e.Next() {
		// Lee el archivo desde el origen
		input, err := ioutil.ReadFile("files/" + fileName)
		// Envía el mensaje a los usuarios
		err = gob.NewEncoder(e.Value.(Usuario).Conexion).Encode(&Mensaje{Opcion: 2, File: Archivo{Nombre: fileName, Datos: input}})
		if err != nil {
			fmt.Println(err)
		}
	}
}

func handleUsuario(c net.Conn, nickname string) {
	// Se manda la notificación a los usuarios conectados actualmente
	enviarMensaje(nickname + " se conectó.")
	//fmt.Println("Se conectó: ", nickname)
	var opc int = 1
	for opc != 0 {
		err := gob.NewDecoder(c).Decode(&opc)
		if err != nil {
			fmt.Println(err)
		}
		/// Switch para los handlers de acciones
		switch opc {
		// Si se captura un 1 envía un "echo" del mensaje al resto de usuarios
		case 1:
			var msg string
			err = gob.NewDecoder(c).Decode(&msg)
			//fmt.Println(nickname, ":", msg)
			enviarMensaje(nickname + ": " + msg)
		// Si se captura un 2 reenvía el archivo al resto de usuarios
		case 2:
			// Crea un directorio para el servidor donde almacena todos los archivos
			_, err := os.Stat("files")
			if os.IsNotExist(err) {
				errDir := os.MkdirAll("files", 0755)
				if errDir != nil {
					fmt.Println(err)
				}
			}
			// Recibe el archivo por la conexión con el usuario
			var arc Archivo
			err = gob.NewDecoder(c).Decode(&arc)
			if err != nil {
				fmt.Println(err)
			}
			err = ioutil.WriteFile("files/"+arc.Nombre, arc.Datos, 0644)
			if err != nil {
				fmt.Println("Error creating", arc.Nombre)
				fmt.Println(err)
				return
			}
			// Agrega al Slice el nombre del archivo
			files = append(files, arc.Nombre)
			// Manda el mensaje al resto de usuarios y archivo a los usuarios conectados.
			//fmt.Println(nickname, " envío: ", arc.Nombre)
			enviarMensaje(nickname + " envío: " + arc.Nombre)
			handleArchivo(arc.Nombre)
		// Si captura un se termina la conexión con el usuario
		case 0:
			//fmt.Println(nickname, " se desconectó.")
			enviarMensaje(nickname + " se desconectó.")
		}
	}
	c.Close()
	eliminarNickname(nickname)
}

func handleConexion(c net.Conn) {
	// Captura el nickname ingresado por el usuario
	var nickname string
	err := gob.NewDecoder(c).Decode(&nickname)
	// Verifica que el paquete recibido no tenga errores
	if err != nil {
		fmt.Println(err)
	}
	// Verifica si es un usuario nuevo, en caso de serlo y existir en el servidor manda error
	if !existeNickname(nickname) {
		// Ingresa el usuario a la lista y se imprime dentro del servidor
		lUsrMensajes.PushBack(Usuario{Nickname: nickname, Conexion: c})
		// Se notifica a la conexión del usuario actual que se conectó sin errores
		var msg string = "Conexion"
		err = gob.NewEncoder(c).Encode(&msg)
		// Crea una instancia hilo para el correspondiente usuario
		handleUsuario(c, nickname)
	} else {
		// Se manda un mensaje de error al usuario y se termina la conexión
		var msg string = "Error"
		err = gob.NewEncoder(c).Encode(&msg)
		c.Close()
	}
}

// Función de servidor que estará escuchando para cuando se conecte un Cliente en el puerto :9999
func server() {
	s, err := net.Listen("tcp", "192.168.100.4:8080")
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
		// Crea un hilo para la conexión actual y espera una nueva
		go handleConexion(c)
	}
}

func main() {
	// Crea un directorio para el servidor donde almacena todos los archivos
	_, err := os.Stat("files")
	if os.IsNotExist(err) {
		errDir := os.MkdirAll("files", 0755)
		if errDir != nil {
			fmt.Println(err)
		}
	}
	// Lanza hilo para el servidor
	go server()
	// Condicionante de paro
	var opc int64 = 1
	for opc != 0 {
		// System Clear
		fmt.Print("\033[H\033[2J")
		// Se imprime la interfaz para el usuario
		fmt.Println("+--------------------------+")
		fmt.Println("|         WutsServer       |")
		fmt.Println("+--------------------------+")
		fmt.Println("| 1.- Mostrar mensajes.    |")
		fmt.Println("| 2.- Mostrar archivos.    |")
		fmt.Println("| 3.- Respaldar mensajes   |")
		fmt.Println("|   y nombres de archivos. |")
		fmt.Println("| 0.- Salir.               |")
		fmt.Println("+--------------------------+")
		fmt.Println("| :                        |")
		fmt.Println("+--------------------------+")
		fmt.Print("\033[10;5H")
		fmt.Scan(&opc)
		switch opc {
		case 1:
			// System Clear
			fmt.Print("\033[H\033[2J")
			fmt.Println("+------------------------------------+")
			fmt.Println("|    Mensajes envíados al servidor   |")
			fmt.Println("+------------------------------------+")
			for _, msg := range msgs {
				fmt.Println(msg)
			}
			fmt.Println("Presione [Enter] para continuar...")
			fmt.Scanln()
			opc = 20
		case 2:
			// System Clear
			fmt.Print("\033[H\033[2J")
			fmt.Println("+------------------------------------+")
			fmt.Println("|    Archivos envíados al servidor   |")
			fmt.Println("+------------------------------------+")
			for _, f := range files {
				fmt.Println(f)
			}
			fmt.Println("Presione [Enter] para continuar...")
			fmt.Scanln()
			opc = 20
		case 3:
			fmt.Println("+--------------------------+")
			fmt.Println("| Mensajes respaldados en  |")
			fmt.Println("| files/Server_Msg.txt     |")
			fmt.Println("+--------------------------+")
			// Se sobreescribe el archivo
			var file, _ = os.Create("files/Server_Msg.txt")
			defer file.Close()
			// Se abre dicho archivo para lectura y escritura
			file, _ = os.OpenFile("files/Server_Msg.txt", os.O_RDWR, 0644)
			defer file.Close()
			// Se ingresan los mensajes actuales del servidor
			for _, msg := range msgs {
				_, err = file.WriteString(msg + "\n")
				if err != nil {
					return
				}
			}
			// Se sincroniza para verificar que no haya error
			err = file.Sync()
			if err != nil {
				return
			}
			fmt.Println("Presione [Enter] para continuar...")
			fmt.Scanln()
			opc = 20
		case 0:
			fmt.Println("+--------------------------+")
			fmt.Println("| Vuelva pronto :D         |")
			fmt.Println("+--------------------------+")
		}
	}
}
