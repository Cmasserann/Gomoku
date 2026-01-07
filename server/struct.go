package main

import "fmt"

type s_table struct {
	size int
	cells [gobanSize]uint8
	captured_b int
	captured_w int
}

type s_StonesPos struct {
	x int
	y int
}

type s_Stones struct {
	b []s_StonesPos
	w []s_StonesPos
	captured_b int
	captured_w int
}

var directions = [][2]int{
			{1, 0},  // horizontal
			{0, 1},  // vertical
			{1, 1},  // diagonal \
			{1, -1}, // diagonal /
		}

// Place a stone on the table at (x, y) with the given color (1 for black, 2 for white).
func putStone(table *s_table, x int, y int, color uint8) bool {
	size := table.size
	if illegalMove(table, x, y, color) {
		return false
	}
	table.cells[y*size+x] = color
	return true
}

func playTurn(table *s_table, x int, y int, color uint8) int {
	if illegalMove(table, x, y, color) {
		return -1
	}
	
	table.cells[y*table.size+x] = color
	capture(table, x, y, color, color)

	if verifWinPoint(table, x, y, color) || getCapturedStones(table, color) >= 5 {
		return 1
	}

	return 0
}

func printTable(table *s_table) {
	size := table.size

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			cell := table.cells[y*size+x]
			if cell == 0 {
				fmt.Print(". ")
			} else {
				fmt.Printf("%d ", cell)
			}
		}	
		fmt.Println()
	}
}

func getCapturedStones(table *s_table, color uint8) int {
	if color == 1 {
		return table.captured_b 
	}
	return table.captured_w


}

func tableToDict(table *s_table) s_Stones {
	result := s_Stones{}
	size := table.size

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			cell := table.cells[y*size+x]
			if cell == 1 {
				result.b = append(result.b, s_StonesPos{x: x, y: y})
			} else if cell == 2 {
				result.w = append(result.w, s_StonesPos{x: x, y: y})
			}
		}
	}
	result.captured_b = table.captured_b
	result.captured_w = table.captured_w

	return result
}

func inbounds(size int, x int, y int) bool {
	return x >= 0 && x < size && y >= 0 && y < size
}

func verifWinPoint(table *s_table, x int, y int, color uint8) bool {
	size := table.size
	count_x := 0
	count_y := 0
	count_d1 := 0
	count_d2 := 0

	for i := -4; i <= 4; i++ {
		if inbounds(size, x+i, y) {
			if table.cells[y*size+(x+i)] == color {
				count_x++
				if count_x >= 5 {
					return true
				}
			} else {
				count_x = 0
			}
		}

		if inbounds(size, x, y+i) {
			if table.cells[(y+i)*size+x] == color {
				count_y++
				if count_y >= 5 {
					return true
				}
			} else {
				count_y = 0
			}
		}

		if inbounds(size, x+i, y+i) {
			if table.cells[(y+i)*size+(x+i)] == color {
				count_d1++
				if count_d1 >= 5 {
					return true
				}
			} else {
				count_d1 = 0
			}
		}

		if inbounds(size, x+i, y-i) {
			if table.cells[(y-i)*size+(x+i)] == color {
				count_d2++
				if count_d2 >= 5 {
					return true
				}
			} else {
				count_d2 = 0
			}
		}
	}

	return false
}


func verifCapturePossible(table *s_table, color uint8) []s_StonesPos {
	size := table.size
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if table.cells[y*size+x] == color {
				result := capture(table, x, y, color, 0)
				if len(result) > 0 {
					return result
				}
			}
		}
	}
	return []s_StonesPos{}
}

func capture(table *s_table, x int, y int, color uint8, endColor uint8) []s_StonesPos {
	size := table.size
	result := []s_StonesPos{}
	opponent := uint8(1)
	if color == 1 {
		opponent = 2
	}

	for i := -1; i <= 1; i += 2 {
		for _, dir := range directions {
			dx := dir[0] * i
			dy := dir[1] * i

			next_x := x + dx
			next_y := y + dy
			mid_x := x + 2 * dx
			mid_y := y + 2 * dy
			end_x := x + 3 * dx
			end_y := y + 3 * dy


			if inbounds(size, next_x, next_y) &&
				inbounds(size, mid_x, mid_y) &&
				inbounds(size, end_x, end_y) {
				if table.cells[next_y*size+next_x] == opponent &&
					table.cells[mid_y*size+mid_x] == opponent &&
					table.cells[end_y*size+end_x] == endColor {

					if endColor == color {
						// Capture the opponent stone
						table.cells[next_y*size+next_x] = 0
						table.cells[mid_y*size+mid_x] = 0
						if color == 1 {
							table.captured_b++
							} else {
								table.captured_w++
						}
					}
					result = append(result, s_StonesPos{x: end_x, y: end_y})
				}
			}
		}
	}
	return result
}

func freeThrees(table s_table, x int, y int, color uint8) int {
	newTable := table
	newTable.cells[y*table.size+x] = color
	
	size := table.size
	count := 0
	for i := -1; i <= 1; i += 2 {
		for _, dir := range directions {

			dx := dir[0] * i
			dy := dir[1] * i

			next_x := x + dx
			next_y := y + dy
			mid_x := x + 2 * dx
			mid_y := y + 2 * dy
			end_x := x + 3 * dx
			end_y := y + 3 * dy

			if inbounds(size, next_x, next_y) &&
				inbounds(size, mid_x, mid_y) &&
				inbounds(size, end_x, end_y) {
				if newTable.cells[next_y*size+next_x] == color &&
					newTable.cells[mid_y*size+mid_x] == color &&
					newTable.cells[end_y*size+end_x] == 0 {
						count++
				} else if newTable.cells[next_y*size+next_x] == color &&
					newTable.cells[mid_y*size+mid_x] == 0 &&
					newTable.cells[end_y*size+end_x] == color {
						count++
				} else if newTable.cells[next_y*size+next_x] == 0 &&
					newTable.cells[mid_y*size+mid_x] == color &&
					newTable.cells[end_y*size+end_x] == color {
						count++
				}
			}

			next_x = x - dx
			next_y = y - dy
			mid_x = x + 1 * dx
			mid_y = y + 1 * dy
			end_x = x + 2 * dx
			end_y = y + 2 * dy

			if inbounds(size, next_x, next_y) &&
				inbounds(size, mid_x, mid_y) &&
				inbounds(size, end_x, end_y) &&
				verifWinPoint(&newTable, x, y, color) == false {
				if newTable.cells[next_y*size+next_x] == color &&
					newTable.cells[mid_y*size+mid_x] == color &&
					newTable.cells[end_y*size+end_x] == 0  && i == -1 {
						count++
				} else if newTable.cells[next_y*size+next_x] == color &&
					newTable.cells[mid_y*size+mid_x] == 0 &&
					newTable.cells[end_y*size+end_x] == color {
						count++
				} 
			}
		}
	}
	return count
}

func illegalMove(table *s_table, x int, y int, color uint8) bool {
	if !inbounds(table.size, x, y) || table.cells[y*table.size+x] != 0 ||
		color != 1 && color != 2 || freeThrees(*table, x, y, color) >= 2 {
		return true
	}
	
	// if table.cells[y*table.size+x] != 0 {
	// 	return true
	// }
	
	// if color != 1 && color != 2 {
	// 	return true
	// }
	
	// if freeThrees(table, x, y, color) >= 2 {
	// 	return true
	// }
	return false
}
