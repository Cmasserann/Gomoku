package main

import "fmt"
import (
	"bufio"
	"net"
	"strings"
)


var table = s_table{ size: 19}


func RunServer() {
	fmt.Println("Server is running...")

	ln, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()
	
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message Received:", string(message))

		data := tableToDict(&table)
		newmessage := fmt.Sprintf("%v", data)
		newmessage = strings.TrimSpace(newmessage)
		fmt.Print("Sending data:", newmessage)

		conn.Write([]byte(newmessage + "\n"))
		break
	}
}

func init() {
	fmt.Println("Initializing server...")

}

func test() {
	// table := s_table{ size: 19}
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
	test()
	
	RunServer()
}