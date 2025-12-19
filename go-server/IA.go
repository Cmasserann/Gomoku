package main

import (
	"fmt"
	"sync"
	"math"
)

type s_ScorePos struct {
	pos   s_StonesPos
	score float64
}

var depth =float64(8)

func opponentColor(color string) string {
	if color == "b" {
		return "w"
	}
	return "b"
}

func IAMain(table s_table, color string) {
	var wg sync.WaitGroup
	
	result := make(chan s_ScorePos, 100)
	moves := getAvailableMoves(table, color)
	
	for _, move := range moves {
		wg.Add(1)
		go func(m s_StonesPos) {
			defer wg.Done()
			newTable := table
			newTable.cells[m.y * newTable.size + m.x] = color
			score := alphaBeta(int(depth)-1, newTable, math.Inf(-1), math.Inf(1), false, color)
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
	maxScore := math.Inf(-1)
	pos := s_StonesPos{x: -1, y: -1}
	for _, p := range listPos {
		if p.score > maxScore {
			maxScore = p.score
			pos = p.pos
		} 
	}
	fmt.Println("IA plays at position:", pos, "with score:", maxScore)
}

func getAvailableMoves(table s_table, color string) []s_StonesPos {
	size := table.size
	availableMoves := []s_StonesPos{}

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if table.cells[y*size+x] == "" && 
				check_close(&table, x, y, color) &&
				!illegalMove(&table, x, y, color) {
				availableMoves = append(availableMoves, s_StonesPos{x: x, y: y})
			}
		}
	}
	return availableMoves
}

func alphaBeta(depth int, table s_table, alpha, beta float64, IsMaximizing bool ,color string) float64 {
	if depth == 0 {
		return evaluateTable(&table, color)
	}

	opponent := opponentColor(color)
	if verifWinPoint(&table, 0, 0, opponent) {
		return float64(-10000)
	}
	if verifWinPoint(&table, 0, 0, color) {
		return float64(10000)
	}

	moves := getAvailableMoves(table, color)
	if len(moves) == 0 {
		return float64(0)
	}

	
	size := table.size
	if IsMaximizing {
		maxEval := math.Inf(-1)
		for _, move := range moves {
			x := move.x
			y := move.y
			if table.cells[y * size + x] == "" && 
				check_close(&table, x, y, color) &&
				!illegalMove(&table, x, y, color) {

					newTable := table
					if putStone(&newTable, x, y, color) == false {
						continue
					}

					if capture(&newTable, x, y, color, color).x != -1 {
						alpha += 0.5
					}

					eval := alphaBeta(depth - 1, newTable, alpha, beta, false, color)

					maxEval = math.Max(maxEval, eval)
					alpha = math.Max(alpha, eval)

					if beta <= alpha {
						break
					}
			}
		}
		return maxEval
	} else {
		minEval := math.Inf(1)
		opponent := opponentColor(color)
		for _, move := range moves {
			x := move.x
			y := move.y
			if table.cells[y * size + x] == "" && 
				check_close(&table, x, y, opponent) &&
				!illegalMove(&table, x, y, opponent) {

					newTable := table
					if putStone(&newTable, x, y, opponent) == false {
						continue
					}
					
					if capture(&newTable, x, y, opponent, opponent).x != -1 {
						beta -= 0.5
					}

					eval := alphaBeta(depth-1, newTable, alpha, beta, true, color)

					minEval = math.Min(minEval, eval)
					beta = math.Min(beta, eval)

					if beta <= alpha {
						break
					}
			}
		}
		return minEval
	}
}


func pruning(table s_table, color string, pos s_StonesPos) float64{
	score := float64(0)
	score += check_win(&table, pos.x, pos.y, color)
	score += check_lose(&table, pos.x, pos.y, color)
	score += check_capture_possible(&table, color)
	// check alignment
	return score
}

func check_close(table *s_table, x int, y int, color string) bool {
	size := table.size

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if inbounds(size, x + i, y + j) {
				if table.cells[(y + j) * size + (x + i)] != "" {
					return true
				}
			}
		}
	}
	return false
}

func check_win(table *s_table, x int, y int, color string) float64 {
	score := float64(0)

	if verifWinPoint(table, x, y, color) {
		score = 10000
	}
	if getCapturedStones(table, color) == 4 && verifCapturePossible(table, color).x != -1 {
		score += 10000
	}
	return score
}

func check_lose(table *s_table, x int, y int, color string) float64 {
	score := float64(0)

	opponent := opponentColor(color)

	if verifWinPoint(table, x, y, opponent) {
		score = -9000
	}
	if getCapturedStones(table, opponent) == 4 && verifCapturePossible(table, opponent).x != -1 {
		score += -9000
	}
	return score
}

func check_capture_possible(table *s_table, color string) float64 {
	score := float64(0)

	if verifCapturePossible(table, color).x != -1 {
		score = 1000
	}
	if verifCapturePossible(table, opponentColor(color)).x != -1 {
		score = -1000
	}
	return score
}

func evaluateTable(table *s_table, color string) float64 {
	score := float64(0)
	size := table.size

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if table.cells[y * size + x] == color {
				score += check_allignments(table, x, y, color)
			} else if table.cells[y * size + x] == opponentColor(color) {
				score -= check_allignments(table, x, y, opponentColor(color))
			}
		}
	}
	return score
}

func check_allignments(table *s_table, x int, y int, color string) float64 {
	score := float64(0)
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
					score += math.Pow(10, float64(count))
					// score -= float64(blockedSides(table, x, y, dx, dy, count))
					count = 0
				}
			}
		}
	}

	return score
}
// 10 ** align_length - (1 if blocked on one side else 0)

