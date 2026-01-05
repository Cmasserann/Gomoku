package main

import (
	// "github.com/gin-gonic/gin"
	"fmt"
)

const gobanWidth = 19
const gobanSize = gobanWidth * gobanWidth

var goban = s_table{ size: gobanWidth, captured_b: 0, captured_w: 0 }

func main() {
	// router := gin.Default()

	// setRouter(router)
	// Output: olleh
	testIA()
}


func testIA() {
	fmt.Println("\nTest IA scenario")
	fmt.Println("-----------------")
	test_table := s_table{ size: 19, captured_b: 0, captured_w: 0 }
	putStone(&test_table, 5, 5, 1)
	putStone(&test_table, 5, 6, 1)
	putStone(&test_table, 6, 5, 1)
	putStone(&test_table, 5, 7, 1)
	putStone(&test_table, 7, 5, 2)
	printTable(&test_table)
	IAMainNoThread(test_table, 2)
}