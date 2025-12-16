package main

import "fmt"
import (
	"bufio"
	"net"
	"strings"
)

var table = s_table{ size: 19, captured_b: 0, captured_w: 0 }

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

}

func test() {
	fmt.Println("Test verifWinPoint")
	fmt.Println("-------------------")
	test_table := s_table{ size: 19, captured_b: 0, captured_w: 0 }
	putStone(&test_table, 14, 16, "b")
	// putStone(&test_table, 15, 16, "b")
	putStone(&test_table, 16, 16, "b")
	putStone(&test_table, 17, 16, "b")
	putStone(&test_table, 18, 16, "b")
	putStone(&test_table, 0, 17, "b")
	putStone(&test_table, 1, 17, "b")
	putStone(&test_table, 2, 17, "b")

	putStone(&test_table, 5, 5, "w")
	putStone(&test_table, 6, 6, "w")
	putStone(&test_table, 7, 7, "w")
	putStone(&test_table, 8, 8, "w")
	putStone(&test_table, 9, 9, "w")

	putStone(&test_table, 8, 12, "b")
	putStone(&test_table, 7, 13, "b")
	putStone(&test_table, 6, 14, "b")
	putStone(&test_table, 5, 15, "b")
	putStone(&test_table, 4, 16, "b")

	printTable(&test_table)

	dict := tableToDict(&test_table)
	fmt.Println("Data To Send:", dict)

	verif := verifWinPoint(&test_table, 16, 16, "b")
	fmt.Println("Test verifWinPoint b 16 16:", verif)
	verif = verifWinPoint(&test_table, 8, 12, "b")
	fmt.Println("Test verifWinPoint b 8 12:", verif)
	verif = verifWinPoint(&test_table, 6, 6, "w")
	fmt.Println("Test verifWinPoint w 6 6:", verif)

}

func testCapture() {
	fmt.Println("\nTest capture scenario")
	fmt.Println("---------------------")
	test_table := s_table{ size: 19, captured_b: 0, captured_w: 0 }
	putStone(&test_table, 11, 8, "w")
	putStone(&test_table, 11, 9, "w")
	putStone(&test_table, 11, 11, "w")
	putStone(&test_table, 11, 12, "w")

	putStone(&test_table, 10, 9, "w")
	putStone(&test_table, 11, 10, "w")
	putStone(&test_table, 12, 11, "b")

	printTable(&test_table)

	x := 11
	y := 10
	color := "w"
	test_table.captured_b = 4
	if getCapturedStones(&test_table, "b") == 4 {
		fmt.Println(test_table.captured_b)
		winPoint := verifWinPoint(&test_table, x, y, color)
		if winPoint.x_start != -1 {
			fmt.Println("Player", color, "wins with line from (", winPoint.x_start, ",", winPoint.y_start, ") to (", winPoint.x_end, ",", winPoint.y_end, ")")
		}
	}
}

func main() {
	test()
	testCapture()
	
	// RunServer()
}