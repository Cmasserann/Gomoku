package main

import "fmt"
import (
	"bufio"
	"net"
	"strings"
)

func RunServer() {
	fmt.Println("Server is running...")

	ln, _ := net.Listen("tcp", ":8080")

	conn, _ := ln.Accept()

	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message Received:", string(message))

        newmessage := strings.ToUpper(message)

		conn.Write([]byte(newmessage + "\n"))
		return
	}
}

func init() {
	fmt.Println("Initializing server...")
	table := s_table{ size: 19}
	// putStone(&table, 3, 3, "w")
	// // putStone(&table, 6, 6, "w")
	// putStone(&table, 6, 5, "w")
	// putStone(&table, 6, 4, "w")
	// putStone(&table, 6, 3, "w")
	// putStone(&table, 6, 2, "w")

	putStone(&table, 14, 16, "b")
	// putStone(&table, 15, 16, "b")
	putStone(&table, 16, 16, "b")
	putStone(&table, 17, 16, "b")
	putStone(&table, 18, 16, "b")
	putStone(&table, 0, 17, "b")
	putStone(&table, 1, 17, "b")
	putStone(&table, 2, 17, "b")


	putStone(&table, 18, 18, "dd")
	putStone(&table, 0, 0, "dd")
	putStone(&table, 0, 18, "rd")
	putStone(&table, 18, 0, "rd")
	
	
	putStone(&table, 5, 5, "w")
	putStone(&table, 6, 6, "w")
	putStone(&table, 7, 7, "w")
	putStone(&table, 8, 8, "w")
	putStone(&table, 9, 9, "w")

	putStone(&table, 8, 12, "b")
	putStone(&table, 7, 13, "b")
	putStone(&table, 6, 14, "b")
	putStone(&table, 5, 15, "b")
	putStone(&table, 4, 16, "b")

	printTable(&table)

	dict := tableToDict(&table)
	fmt.Println("Data To Send:", dict)

	player := []string{"b", "w"}
	end := GameEnded(&table, player)
	fmt.Println("Game Ended:", end)

	test := verifWinPoint(&table, 16, 16, "b")
	fmt.Println("Test verifWinPoint b :", test)
	test = verifWinPoint(&table, 6, 4, "w")
	fmt.Println("Test verifWinPoint w :", test)

}

func main() {
}