// Luis Fernando Martínez Castellanos
// 216787787
// Sistemas concurrentes y distribuidos
// L - M  7AM

package main

import (
	"fmt"
)

// Struct para los procesos
type Usuario struct {
	Nickname string
}

func main() {
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
		fmt.Scan(&opc)
		switch opc {
		case 1:
			var mensaje string
			fmt.Print("\033[11;1H Tu:")
			fmt.Scan(&mensaje)
		case 2:
			var archivo string
			fmt.Print("\033[11;1H Dirección del archivo:")
			fmt.Scan(&archivo)
		case 3:

		case 0:
			fmt.Println("+-------------------------------------+")
			fmt.Println("| Gracias por usar el software.       |")
			fmt.Println("+-------------------------------------+")
		default:
		}
	}
}
