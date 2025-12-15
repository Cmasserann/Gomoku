package main

import "fmt"


type s_table struct {
	size int
	cells [19*19]string
}

type s_StonesPos struct {
	x int
	y int
}

type s_Stones struct {
	b []s_StonesPos
	w []s_StonesPos
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

	return result
}

func GameEnded(table *s_table, player []string) string {
	size := len(table.cells)

	for playerIndex := 0; playerIndex < len(player); playerIndex++ {
		count := 0

		color := player[playerIndex]
		for i := 0; i < size; i++ {
			if i % table.size == 0 {
				count = 0
			}

			if table.cells[i] == color {
				count++
				if count >= 5 { return color }
			} else { count = 0 }
		}
	}
	return "n"
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