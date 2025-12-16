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

type s_WinPoint struct {
	x_start int
	y_start int
	x_end int
	y_end int
}


func putStone(table *s_table, x int, y int, color string) {
	table.cells[y*19+x] = color
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


// return (bool, (x_start, y_start), (x_end, y_end))
func verifWinPoint(table *s_table, x int, y int, color string) s_WinPoint {
	size := table.	size
	count_x := 0
	count_y := 0
	count_d1 := 0
	count_d2 := 0

	for i := -4; i <= 4; i++ {
		if x+i >= 0 && x+i < size {
			if table.cells[y*size+(x+i)] == color {
				count_x++
				if count_x >= 5 {
					return s_WinPoint{x_start: x + i - 4, y_start: y, x_end: x + i, y_end: y}
				}
			} else {
				count_x = 0
			}
		}

		if y+i >= 0 && y+i < size {
			if table.cells[(y+i)*size+x] == color {
				count_y++
				if count_y >= 5 {
					return s_WinPoint{x_start: x, y_start: y + i - 4, x_end: x, y_end: y + i}
				}
			} else {
				count_y = 0
			}
		}

		if x+i >= 0 && x+i < size && y+i >= 0 && y+i < size {
			if table.cells[(y+i)*size+(x+i)] == color {
				count_d1++
				if count_d1 >= 5 {
					return s_WinPoint{x_start: x + i - 4, y_start: y + i - 4, x_end: x + i, y_end: y + i}
				}
			} else {
				count_d1 = 0
			}
		}

		if x+i >= 0 && x+i < size && y-i >= 0 && y-i < size {
			if table.cells[(y-i)*size+(x+i)] == color {
				count_d2++
				if count_d2 >= 5 {
					return s_WinPoint{x_start: x + i - 4, y_start: y - i + 4, x_end: x + i, y_end: y - i}
				}
			} else {
				count_d2 = 0
			}
		}
	}

	return s_WinPoint{x_start: -1, y_start: -1, x_end: -1, y_end: -1}
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

	directions := [][2]int{
		{1, 0},  // horizontal
		{0, 1},  // vertical
		{1, 1},  // diagonal \
		{1, -1}, // diagonal /
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


			if next_x >= 0 && next_x < size && next_y >= 0 && next_y < size &&
				mid_x >= 0 && mid_x < size && mid_y >= 0 && mid_y < size &&
				end_x >= 0 && end_x < size && end_y >= 0 && end_y < size {
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
