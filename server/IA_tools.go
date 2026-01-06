package main

func opponentColor(color uint8) uint8 {
	if color == 1 {
		return 2
	}
	return 1
}


func checkAlignement(table *s_table, x int, y int, color uint8) int {
	count := 0

	for i := -1 ; i <= 1; i += 2 {
		count += checkOneDirection(table, x, y, color, 1 * i, 0)
		count += checkOneDirection(table, x, y, color, 0, 1 * i)
		count += checkOneDirection(table, x, y, color, 1 * i, 1 * i)
		count += checkOneDirection(table, x, y, color, 1 * i, -1 * i)
	}
	return count
}

func checkOneDirection(table *s_table, x int, y int, color uint8, dx int, dy int) int {
	size := table.size
	count := 0

	for i := 1; i <= 4; i++ {
		nx := x + i * dx
		ny := y + i * dy
		if inbounds(size, nx, ny) && table.cells[ny * size + nx] == color {
			count++
		} else if count > 0 {
			if table.cells[ny * size + nx] != 0 {
				return pow10[count + 1]
			}
			return pow10[count]
		} else {
			break
		}
	}
	return count
}

func setAvailableMoves(table s_table, color uint8) s_table {
	size := table.size
	availableMovesTable := s_table{
		size:      size,
		captured_b: 0,
		captured_w: 0,
	}
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if table.cells[y * size + x] == 0 && 
				check_close(&table, x, y, color) && 
				!illegalMove(&table, x, y, color) {
				availableMovesTable.cells[y*size+x] = 1
			}
		}
	}
	return availableMovesTable
}

func updateAvailableMoves(table s_table, color uint8, x int, y int) s_table {
	size := table.size
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			nx := x + i
			ny := y + j
			if inbounds(size, nx, ny) && table.cells[ny*size+nx] == 0 {
				if check_close(&table, nx, ny, color) {
					table.cells[ny * size + nx] = 1
				}
			}
		}
	}
	table.cells[y*size+x] = 0
	return table
}


func updateAvailableMovesAfterCapture(table s_table, color uint8, capturedStones []s_StonesPos) s_table {
	size := table.size
	for _, pos := range capturedStones {
		x := pos.x
		y := pos.y
		for i := -1; i <= 1; i++ {
			for j := -1; j <= 1; j++ {
				nx := x + i
				ny := y + j
				if inbounds(size, nx, ny) && table.cells[ny*size+nx] == 0 {
					if check_close(&table, nx, ny, color) {
						table.cells[ny * size + nx] = 1
					}
				}
			}
		}
	}
	return table
}

func getAvailableMoves(table s_table, color uint8) []s_StonesPos {
	size := table.size
	availableMoves := []s_StonesPos{}

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if table.cells[y * size + x] == 1{
				availableMoves = append(availableMoves, s_StonesPos{x: x, y: y})
			}
		}
	}
	return availableMoves
}

func check_close(table *s_table, x int, y int, color uint8) bool {
	size := table.size

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if inbounds(size, x + i, y + j) {
				if table.cells[(y + j) * size + (x + i)] != 0 {
					return true
				}
			}
		}
	}
	return false
}
