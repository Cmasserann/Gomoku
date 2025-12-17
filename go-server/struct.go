package main

import "fmt"


type s_table struct {
	size int
	cells [19*19]string
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

func putStone(table *s_table, x int, y int, color string) bool {
	size := table.size
	if illegalMove(table, x, y, color) {
		return false
	}
	table.cells[y*size+x] = color
	return true
}

func printTable(table *s_table) {
	size := table.size

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			cell := table.cells[y*size+x]
			if cell == "" {
				fmt.Print(". ")
			} else {
				fmt.Printf("%s ", cell)
			}
		}	
		fmt.Println()
	}
}

func getCapturedStones(table *s_table, color string) int {
	if color == "b" {
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
			if cell == "b" {
				result.b = append(result.b, s_StonesPos{x: x, y: y})
			} else if cell == "w" {
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

func verifWinPoint(table *s_table, x int, y int, color string) bool {
	size := table.	size
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


func verifCapturePossible(table *s_table, color string) s_StonesPos {
	size := table.size
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if table.cells[y*size+x] == color {
				result := capture(table, x, y, color, "")
				if result.x != -1 {
					return result
				}
			}
		}
	}
	return s_StonesPos{x: -1, y: -1}
}

func capture(table *s_table, x int, y int, color string, endColor string) s_StonesPos {
	size := table.size
	opponent := "b"
	if color == "b" {
		opponent = "w"
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

					if endColor != color {
						return s_StonesPos{x: end_x, y: end_y}
					}
					// Capture the opponent stone
					table.cells[next_y*size+next_x] = ""
					table.cells[mid_y*size+mid_x] = ""
					if color == "b" {
						table.captured_b++
					} else {
						table.captured_w++
					}
					return s_StonesPos{x: end_x, y: end_y}
				}
			}
		}
	}
	return s_StonesPos{x: -1, y: -1}
}

func freeThrees(table *s_table, x int, y int, color string) int {
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
				if table.cells[next_y*size+next_x] == color &&
					table.cells[mid_y*size+mid_x] == color &&
					table.cells[end_y*size+end_x] == "" {
						count++
				} else if table.cells[next_y*size+next_x] == color &&
					table.cells[mid_y*size+mid_x] == "" &&
					table.cells[end_y*size+end_x] == color {
						count++
				} else if table.cells[next_y*size+next_x] == "" &&
					table.cells[mid_y*size+mid_x] == color &&
					table.cells[end_y*size+end_x] == color {
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
				inbounds(size, end_x, end_y) {
				if table.cells[next_y*size+next_x] == color &&
					table.cells[mid_y*size+mid_x] == color &&
					table.cells[end_y*size+end_x] == "" {
						count++
				} else if table.cells[next_y*size+next_x] == color &&
					table.cells[mid_y*size+mid_x] == "" &&
					table.cells[end_y*size+end_x] == color {
						count++
				} 
			}
		}
	}
	return count
}

func illegalMove(table *s_table, x int, y int, color string) bool {
	if !inbounds(table.size, x, y) {
		return true
	}
	
	if table.cells[y*table.size+x] != "" {
		return true
	}
	
	if color != "b" && color != "w" {
		return true
	}
	
	if freeThrees(table, x, y, color) >= 2 {
		return true
	}
	return false
}
