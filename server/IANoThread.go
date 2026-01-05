package main

import (
	"fmt"
)

var endDepth = 0

func IAMainNoThread(table s_table, color uint8) {
	// result := []s_ScorePos{}
	availableMovesTable := setAvailableMoves(table, color)

	resultRecurse := RecursiveSearch(depth, table, availableMovesTable, true, color)

	fmt.Println("IA possible positions:", resultRecurse)
}

func RecursiveSearch(depth int, table s_table, availableMovesTable s_table, AIMove bool, color uint8) s_ScorePos {
	moves := getAvailableMoves(availableMovesTable, color)
	
	if (endDepth != 0 && depth <= endDepth) || depth == 0 || len(moves) == 0 {
		return s_ScorePos{pos: s_StonesPos{x: -1, y: -1}, score: 0}
	}
	
	scoreMove := make([]s_ScorePos, len(moves))


	for _, move := range moves {
		score := 0
		newTable := table
		if (illegalMove(&newTable, move.x, move.y, color)) {
			continue
		}
		
		if (verifWinPoint(&newTable, move.x, move.y, color)) {
			endDepth = depth
			if (AIMove) {
				return s_ScorePos{pos: move, score: 10000000000}
			} else {
				return s_ScorePos{pos: move, score: -10000000000}
			}
		}

		putStone(&newTable, move.x, move.y, color)
		capture := capture(&newTable, move.x, move.y, color, color)

		if (getCapturedStones(&newTable, opponentColor(color)) >= 5) {
			endDepth = depth
			if (AIMove) {
				return s_ScorePos{pos: move, score: 10000000000}
			} else {
				return s_ScorePos{pos: move, score: -10000000000}
			}
		}

		if len(capture) > 0 {
			score += 500
		}
		score += checkAlignement(&newTable, move.x, move.y, color)
		score += checkAlignement(&newTable, move.x, move.y, opponentColor(color))
		scoreMove = append(scoreMove, s_ScorePos{pos: move, score: score})
	}


	maxScore := -10000000000
	var bestMove s_ScorePos
	for _, sm := range scoreMove {
		if (sm.score > maxScore) {
			maxScore = sm.score
			bestMove = sm
		}
	}

	newTable := table
	putStone(&newTable, bestMove.pos.x, bestMove.pos.y, color)
	availableMovesTable = updateAvailableMoves(availableMovesTable, color, bestMove.pos.x, bestMove.pos.y)
	result := RecursiveSearch(depth - 1, newTable, availableMovesTable, !AIMove, opponentColor(color))
	if (AIMove) {
		result.score = result.score + bestMove.score * pow10[depth]
	} else {
		result.score = result.score - bestMove.score * pow10[depth]
	}
	result.pos = bestMove.pos
	fmt.Println("Best Move at depth", depth, ":", bestMove)
	fmt.Println("score:", result.score)
	return result
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
			return pow10[count]
		} else {
			break
		}
	}
	return count
}