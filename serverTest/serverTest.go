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
	putStone(&test_table, 14, 16, 1)
	// putStone(&test_table, 15, 16, "b")
	putStone(&test_table, 16, 16, 1)
	putStone(&test_table, 17, 16, 1)
	putStone(&test_table, 18, 16, 1)
	putStone(&test_table, 0, 17, 1)
	putStone(&test_table, 1, 17, 1)
	putStone(&test_table, 2, 17, 1)

	putStone(&test_table, 5, 5, 2)
	putStone(&test_table, 6, 6, 2)
	putStone(&test_table, 7, 7, 2)
	putStone(&test_table, 8, 8, 2)
	putStone(&test_table, 9, 9, 2)

	putStone(&test_table, 8, 12, 1)
	putStone(&test_table, 7, 13, 1)
	putStone(&test_table, 6, 14, 1)
	putStone(&test_table, 5, 15, 1)
	putStone(&test_table, 4, 16, 1)
	printTable(&test_table)

	dict := tableToDict(&test_table)
	fmt.Println("Data To Send:", dict)

	verif := verifWinPoint(&test_table, 16, 16, 1)
	fmt.Println("Test verifWinPoint b 16 16:", verif)
	verif = verifWinPoint(&test_table, 8, 12, 1)
	fmt.Println("Test verifWinPoint b 8 12:", verif)
	verif = verifWinPoint(&test_table, 6, 6, 2)
	fmt.Println("Test verifWinPoint w 6 6:", verif)

}

func testCapture() {
	fmt.Println("\nTest capture scenario")
	fmt.Println("---------------------")
	test_table := s_table{ size: 19, captured_b: 0, captured_w: 0 }
	putStone(&test_table, 11, 8, 2)
	putStone(&test_table, 11, 9, 2)
	putStone(&test_table, 11, 10, 2)
	putStone(&test_table, 11, 11, 2)
	putStone(&test_table, 11, 12, 2)

	putStone(&test_table, 10, 10, 2)
	putStone(&test_table, 9, 10, 1)

	printTable(&test_table)

	x := 11
	y := 10
	color := uint8(2)
	opponentColor := uint8(1)
	test_table.captured_b = 4
	if verifWinPoint(&test_table, x, y, color) {
		fmt.Println(test_table.captured_b)
		fmt.Println(color, " win by five in a row at (", x, ",", y, ")")
		result := verifCapturePossible(&test_table, opponentColor)
		if getCapturedStones(&test_table, opponentColor) == 4 && len(result) > 0 {
				fmt.Println("But ", opponentColor, " win by capturing possible before placing stone at (", result, ")")
		}
	}

	putStone(&test_table, 12, 10, 1)
	printTable(&test_table)
	
	capture(&test_table, 12, 10, 1, 1)
	fmt.Println("After capture:")
	printTable(&test_table)
	fmt.Println("Captured black stones:", getCapturedStones(&test_table, 1))
}

func testIllegalMove() {
	fmt.Println("\nTest illegal move scenario")
	fmt.Println("-------------------------")
	test_table := s_table{ size: 19, captured_b: 0, captured_w: 0 }
	putStone(&test_table, 10, 8, 2)
	putStone(&test_table, 10, 9, 2)
	putStone(&test_table, 9, 10, 2)
	putStone(&test_table, 7, 10, 2)

	putStone(&test_table, 1, 0, 2)
	putStone(&test_table, 1, 2, 2)
	putStone(&test_table, 0, 0, 2)
	putStone(&test_table, 3, 3, 2)
	putStone(&test_table, 3, 1, 2)
	putStone(&test_table, 4, 1, 2)
	printTable(&test_table)

	result := illegalMove(&test_table, 10, 10, 2)
	if result {
		fmt.Println("Move at (10,10) for white is illegal (double free three)")
	} else {
		fmt.Println("Move at (10,10) for white is legal")
	}
	result = illegalMove(&test_table, 10, 10, 1)
	if result {
		fmt.Println("Move at (10,10) for black is illegal (double free three)")
	} else {
		fmt.Println("Move at (10,10) for black is legal")
	}

	result = illegalMove(&test_table, 1, 1, 2)
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
		putStone(&test_table, i, 0, 1)
		putStone(&test_table, 0, i, 1)
		putStone(&test_table, i, i, 1)
		putStone(&test_table, i, 4-i, 1)
	}
	printTable(&test_table)
}

func testIA() {
	fmt.Println("\nTest IA scenario")
	fmt.Println("-----------------")
	test_table := s_table{ size: 19, captured_b: 0, captured_w: 0 }
	putStone(&test_table, 5, 5, 1)
	putStone(&test_table, 6, 5, 2)
	printTable(&test_table)
	IAMainNoThread(test_table, 2)
}

func main() {
	// test()
	// testCapture()
	// testIllegalMove()
	// testWin()
	testIA()


	
	// RunServer()
}