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
	putStone(&test_table, 11, 10, "w")
	putStone(&test_table, 11, 11, "w")
	putStone(&test_table, 11, 12, "w")

	putStone(&test_table, 10, 10, "w")
	putStone(&test_table, 12, 10, "b")

	printTable(&test_table)

	x := 11
	y := 10
	color := "w"
	opponentColor := "b"
	test_table.captured_b = 4
	if verifWinPoint(&test_table, x, y, color) {
		fmt.Println(test_table.captured_b)
		fmt.Println(color, " win by five in a row at (", x, ",", y, ")")
		result := verifCapturePossible(&test_table, opponentColor)
		if getCapturedStones(&test_table, opponentColor) == 4 && result.x != -1 {
				fmt.Println("But ", opponentColor, " win by capturing possible before placing stone at (", result.x, ",", result.y, ")")
		}
	}

	putStone(&test_table, 9, 10, "b")
	printTable(&test_table)
	
	capture(&test_table, 9, 10, "b", "b")
	fmt.Println("After capture:")
	printTable(&test_table)
	fmt.Println("Captured black stones:", getCapturedStones(&test_table, "b"))
}

func testIllegalMove() {
	fmt.Println("\nTest illegal move scenario")
	fmt.Println("-------------------------")
	test_table := s_table{ size: 19, captured_b: 0, captured_w: 0 }
	putStone(&test_table, 10, 8, "w")
	putStone(&test_table, 10, 9, "w")
	putStone(&test_table, 9, 10, "w")
	putStone(&test_table, 7, 10, "w")
	test_table.cells[10*19+10] = "T"

	putStone(&test_table, 1, 0, "w")
	putStone(&test_table, 1, 2, "w")
	putStone(&test_table, 0, 0, "w")
	putStone(&test_table, 3, 3, "w")
	putStone(&test_table, 3, 1, "w")
	putStone(&test_table, 4, 1, "w")
	test_table.cells[1*19+1] = "T"
	printTable(&test_table)
	result := illegalMove(&test_table, 10, 10, "w")
	if result {
		fmt.Println("Move at (10,10) for white is illegal (double free three)")
	} else {
		fmt.Println("Move at (10,10) for white is legal")
	}
	result = illegalMove(&test_table, 10, 10, "b")
	if result {
		fmt.Println("Move at (10,10) for blqck is illegal (double free three)")
	} else {
		fmt.Println("Move at (10,10) for blqck is legal")
	}

	result = illegalMove(&test_table, 1, 1, "w")
	if result {
		fmt.Println("Move at (1,1) for white is illegal (double free three)")
	} else {
		fmt.Println("Move at (1,1) for white is legal")
	}

}

func testWin() {
	fmt.Println("\nTest win scenario")
	fmt.Println("-----------------")
	test_table := s_table{ size: 19, captured_b: 0, captured_w: 0 }
	for i := 0; i < 5; i++ {
		putStone(&test_table, i, 0, "b")
		putStone(&test_table, 0, i, "b")
		putStone(&test_table, i, i, "b")
		putStone(&test_table, i, 4-i, "b")
	}
	printTable(&test_table)
}


func main() {
	// test()
	testCapture()
	// testIllegalMove()
	// testWin()


	
	// RunServer()
}