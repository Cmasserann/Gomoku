package main

import (
	"fmt"
	"sync"
)

type s_ScorePos struct {
	pos   s_StonesPos
	score int
}

var depth = 8

var pow10 = []int{
	1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000,
}

func opponentColor(color uint8) uint8 {
	if color == 1 {
		return 2
	}
	return 1
}

func IAMain(table s_table, color uint8) {
	var wg sync.WaitGroup
	
	result := make(chan s_ScorePos, 100)
	availableMovesTable := setAvailableMoves(table, color)
	moves := getAvailableMoves(availableMovesTable, color)
	
	for _, move := range moves {
		wg.Add(1)
		go func(m s_StonesPos) {
			defer wg.Done()
			newTable := table
			putStone(&newTable, m.x, m.y, color)
			capture := capture(&newTable, m.x, m.y, color, color)

			score := alphaBeta(int(depth) - 1, newTable, updateAvailableMoves(availableMovesTable, color, m.x, m.y), -1000000000, 1000000000, false, color)
			if len(capture) > 0 {
				score += 500
			}
			result <- s_ScorePos{pos: m, score: score}
		}(move)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	listPos := []s_ScorePos{}
	for pos := range result {
		listPos = append(listPos, pos)
	}
	fmt.Println("IA possible positions:", listPos)
	
	// max
	maxScore := -1000000
	pos := s_StonesPos{x: -1, y: -1}
	for _, p := range listPos {
		if p.score > maxScore {
			maxScore = p.score
			pos = p.pos
		} 
	}
	fmt.Println("IA plays at position:", pos, "with score:", maxScore)
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
	return table
}


func alphaBeta(depth int, table s_table, availableMovesTable s_table, alpha, beta int, IsMaximizing bool ,color uint8) int {
	if depth == 0 {
		return evaluateTable(&table, color)
	}

	opponent := opponentColor(color)
	if verifWinPoint(&table, 0, 0, opponent) {
		return -10000
	}
	if verifWinPoint(&table, 0, 0, color) {
		return 10000
	}

	moves := getAvailableMoves(availableMovesTable, color)
	if len(moves) == 0 {
		return 0
	}
	
	size := table.size
	if IsMaximizing {
		maxEval := -1000000000
		for _, move := range moves {
			x := move.x
			y := move.y
			if table.cells[y * size + x] == 0 && 
				check_close(&table, x, y, color) &&
				!illegalMove(&table, x, y, color) {

					newTable := table
					if putStone(&newTable, x, y, color) == false {
						continue
					}

					if len(capture(&newTable, x, y, color, color)) > 0 {
						alpha += 500
					}

					eval := alphaBeta(depth - 1, newTable, updateAvailableMoves(availableMovesTable, color, x, y), alpha, beta, false, color)

					if eval > maxEval {
						maxEval = eval
					}
					if eval > alpha {
						alpha = eval
					}

					if beta <= alpha {
						break
					}
			}
		}
		return maxEval
	} else {
		minEval := 1000000000
		opponent := opponentColor(color)
		for _, move := range moves {
			x := move.x
			y := move.y
			if table.cells[y * size + x] == 0 && 
				check_close(&table, x, y, opponent) &&
				!illegalMove(&table, x, y, opponent) {

					newTable := table
					if putStone(&newTable, x, y, opponent) == false {
						continue
					}
					
					localAlpha := alpha
					localBeta := beta
					if len(capture(&newTable, x, y, opponent, opponent)) > 0 {
						localBeta -= 500
					}

					eval := alphaBeta(depth-1, newTable, updateAvailableMoves(availableMovesTable, opponent, x, y), localAlpha, localBeta, true, color)

					if eval < minEval {
						minEval = eval
					}
					if eval < beta {
						beta = eval
					}

					if beta <= alpha {
						break
					}
			}
		}
		return minEval
	}
}


func pruning(table s_table, color uint8, pos s_StonesPos) int{
	score := 0
	score += check_win(&table, pos.x, pos.y, color)
	score += check_lose(&table, pos.x, pos.y, color)
	score += check_capture_possible(&table, color)
	// check alignment
	return score
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

func check_win(table *s_table, x int, y int, color uint8) int {
	score := 0

	if verifWinPoint(table, x, y, color) {
		score = 10000
	}
	if getCapturedStones(table, color) == 4 && len(verifCapturePossible(table, color)) > 0 {
		score += 10000
	}
	return score
}

func check_lose(table *s_table, x int, y int, color uint8) int {
	score := 0

	opponent := opponentColor(color)

	if verifWinPoint(table, x, y, opponent) {
		score = -9000
	}
	if getCapturedStones(table, opponent) == 4 && len(verifCapturePossible(table, opponent)) > 0 {
		score += -9000
	}
	return score
}

func check_capture_possible(table *s_table, color uint8) int {
	score := 0

	if len(verifCapturePossible(table, color)) > 0 {
		score = 1000
	}
	if len(verifCapturePossible(table, opponentColor(color))) > 0 {
		score = -1000
	}
	return score
}

func evaluateTable(table *s_table, color uint8) int {
	score := 0
	size := table.size

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if table.cells[y * size + x] == color {
				score += check_alignments(table, x, y, color)
			} else if table.cells[y * size + x] == opponentColor(color) {
				score -= check_alignments(table, x, y, opponentColor(color))
			}
		}
	}
	return score
}

func check_alignments(table *s_table, x int, y int, color uint8) int {
	score := 0
	size := table.size

	count := 0
	for _, dir := range directions {
		dx := dir[0]
		dy := dir[1]	
		count = 0
		for i := -4; i <= 4; i++ {
			nx := x + i * dx
			ny := y + i * dy
			if inbounds(size, nx, ny) && table.cells[ny * size + nx] == color {
				count++
			} else {
				if count > 0 {
					score += pow10[count]
					// score += math.Pow(10, float64(count))
					// score -= float64(blockedSides(table, x, y, dx, dy, count))
					count = 0
				}
			}
		}
	}

	return score
}
// 10 ** align_length - (1 if blocked on one side else 0)

