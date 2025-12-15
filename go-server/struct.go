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


func verifWinPoint(table *s_table, x int, y int, color string) bool {
	size := table.size
	count_x := 0
	count_y := 0
	count_d1 := 0
	count_d2 := 0

	for i := -4; i <= 4; i++ {
		if x+i >= 0 && x+i < size {
			if table.cells[y*size+(x+i)] == color {
				count_x++
				if count_x >= 5 {
					return true
				}
			} else {
				count_x = 0
			}
		}

		if y+i >= 0 && y+i < size {
			if table.cells[(y+i)*size+x] == color {
				count_y++
				if count_y >= 5 {
					return true
				}
			} else {
				count_y = 0
			}
		}

		if x+i >= 0 && x+i < size && y+i >= 0 && y+i < size {
			if table.cells[(y+i)*size+(x+i)] == color {
				count_d1++
				if count_d1 >= 5 {
					return true
				}
			} else {
				count_d1 = 0
			}
		}

		if x+i >= 0 && x+i < size && y-i >= 0 && y-i < size {
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

func capturePossibe(table *s_table, x int, y int, color string) bool {
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

	// test if xy in the end of a capture
	for i := -1; i <= 1; i += 2 {
		for _, dir := range directions {
			dx := dir[0] * i
			dy := dir[1] * i

			mid_x := x + dx
			mid_y := y + dy
			end_x := x + 2 * dx
			end_y := y + 2 * dy
			empty_x := x + dx * i
			empty_y := y + dy * i

			if mid_x >= 0 && mid_x < size && mid_y >= 0 && mid_y < size &&
				end_x >= 0 && end_x < size && end_y >= 0 && end_y < size &&
				empty_x >= 0 && empty_x < size && empty_y >= 0 && empty_y < size {
				if table.cells[mid_y*size+mid_x] == color &&
					table.cells[end_y*size+end_x] == opponent &&
					table.cells[empty_y*size+empty_x] == "" {
					return true
				}
			}
		}
	}

	// test if xy in the middle of a capture
	for i := -1; i <= 1; i += 2 {
		for _, dir := range directions {
			dx := dir[0] * i
			dy := dir[1] * i

			start_x := x - dx
			start_y := y - dy
			end_x := x + dx
			end_y := y + dy
			empty_x := x - 2 * dx
			empty_y := y - 2 * dy

			if start_x >= 0 && start_x < size && start_y >= 0 && start_y < size &&
				end_x >= 0 && end_x < size && end_y >= 0 && end_y < size &&
				empty_x >= 0 && empty_x < size && empty_y >= 0 && empty_y < size {
				if table.cells[start_y*size+start_x] == opponent &&
					table.cells[end_y*size+end_x] == color &&
					table.cells[empty_y*size+empty_x] == "" {
					return true
				}
			}
		}
	}


	return false
}