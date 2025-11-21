package main

import "fmt"


type s_table struct {
	size int
	h_cells [19*19]string
	v_cells [19*19]string
	rd_cells [19*19]string
	dd_cells [19*19]string
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
	table.h_cells[y*19+x] = color
	table.v_cells[x*19+y] = color
	table.dd_cells[(18 - x)*19 + y] = color
	table.rd_cells[(18 - y)*19 + x] = color
}

func printTable(table *s_table) {
	size := table.size

	fmt.Println("Horizontal:")
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			cell := table.h_cells[y*size+x]
			if cell == "" {
				fmt.Print(". ")
			} else {
				fmt.Printf("%s ", cell)
			}
		}	
		fmt.Println()
	}

	fmt.Println("\nVertical:")
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			cell := table.v_cells[x*size+y]
			if cell == "" {
				fmt.Print(". ")
			} else {
				fmt.Printf("%s ", cell)
			}
		}
		fmt.Println()
	}

	//TODO: Fix diagonal 
	fmt.Println("\nDiagonal Descending:")
	// for i := 0; i < size*2-1; i++ {
	// 	for y := 0; y < size; y++ {
	// 		x := i - y
	// 		if x >= 0 && x < size {
	// 			cell := table.dd_cells[x*size+y]
	// 			if cell == "" {
	// 				fmt.Print(". ")
	// 			} else {
	// 				fmt.Printf("%s ", cell)
	// 			}
	// 		}
	// 	}
	// 	fmt.Println()
	// }

	// fmt.Println("\nDiagonal Rising:")
	// for i := 0; i < size*2-1; i++ {
	// 	for y := 0; y < size; y++ {
	// 		x := i - (size - 1 - y)
	// 		if x >= 0 && x < size {
	// 			cell := table.rd_cells[x*size+y]
	// 			if cell == "" {
	// 				fmt.Print(". ")
	// 			} else {
	// 				fmt.Printf("%s ", cell)
	// 			}
	// 		}
	// 	}
	// 	fmt.Println()
	// }	
}

func tableToDict(table *s_table) s_Stones {
	result := s_Stones{}
	size := table.size

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			cell := table.h_cells[y*size+x]
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
	size := len(table.h_cells)

	for playerIndex := 0; playerIndex < len(player); playerIndex++ {
		count_h := 0
		count_v := 0

		color := player[playerIndex]
		for i := 0; i < size; i++ {
			if i % table.size == 0 {
				count_h = 0
				count_v = 0
			}

			// Horizontal
			if table.h_cells[i] == color {
				count_h++
				if count_h >= 5 { return color }
			} else { count_h = 0 }

			// Vertical
			if table.v_cells[i] == color {
				count_v++
				if count_v >= 5 { return color }
			} else { count_v = 0 }
		}

		count_rd := 0
		count_dd := 0
		checkfor_d := 10 + 5
		count_d := 1
		for i := 10; i < size - 10; i++ {
			if i == checkfor_d {
				fmt.Println("count_dd:", count_dd, " count_rd:", count_rd)
				count_dd = 0
				count_rd = 0
				checkfor_d += 5 + count_d
				count_d++
			}

			// Diagonal Descending
			if table.dd_cells[i] == color {
				count_dd++
				if count_dd >= 5 { return color }
			} else { count_dd = 0 }

			// Diagonal Rising
			if table.rd_cells[i] == color {
				count_rd++
				if count_rd >= 5 { return color }
			} else { count_rd = 0 }		
		}

		fmt.Println("--------")
	}
	return "n"
}

func verifWinPoint(table *s_table, x int, y int, color string) bool {
	size := table.size
	count_x := 0
	count_y := 0

	for i := -4; i <= 4; i++ {
		// Vérification horizontale
		if x+i >= 0 && x+i < size {
			if table.h_cells[y*size+(x+i)] == color {
				count_x++
				if count_x >= 5 {
					return true
				}
			} else {
				count_x = 0
			}
		}

		// Vérification verticale
		if y+i >= 0 && y+i < size {
			if table.v_cells[x*size+(y+i)] == color{
				count_y++
				if count_y >= 5 {
					return true
				}
			} else {
				count_y = 0
			}
		}
	}

	return false
}