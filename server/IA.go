package main

import (
	"fmt"
)

type s_ScorePos struct {
	pos   s_StonesPos
	score int
}

var depth = 10

var pow10 = []int{
	1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000, 10000000000,
}

var endDepth = 0

func getIAMove(table s_table, color uint8) s_StonesPos {
	endDepth = 0
	availableMovesTable := setAvailableMoves(table, color)

	resultRecurse := RecursiveSearch(depth, table, availableMovesTable, true, color)

	fmt.Println("IA possible positions:", resultRecurse)
	return resultRecurse.pos
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
		
		putStone(&newTable, move.x, move.y, color)

		capture := capture(&newTable, move.x, move.y, color, color)

		if (getCapturedStones(&newTable, color) >= 5) {
			endDepth = depth
			scoreMove = make([]s_ScorePos, 0)
			scoreMove = append(scoreMove, s_ScorePos{pos: move, score: 10000000000})
			break
			// return s_ScorePos{pos: move, score: 10000000000}
		}

		if (verifWinPoint(&newTable, move.x, move.y, color)) {
			endDepth = depth
			captured := verifCapturePossible(&newTable, opponentColor(color))
			if len(captured) + getCapturedStones(&newTable, opponentColor(color)) >= 5 {
				continue
			}
			scoreMove = make([]s_ScorePos, 0)
			scoreMove = append(scoreMove, s_ScorePos{pos: move, score: 10000000000})
			break
			// return s_ScorePos{pos: move, score: 10000000000}
		}


		if len(capture) > 0 {
			score += 500
		}
		score += checkAlignement(&newTable, move.x, move.y, color)
		score += checkAlignement(&newTable, move.x, move.y, opponentColor(color))
		scoreMove = append(scoreMove, s_ScorePos{pos: move, score: score})
	}

	// fmt.Println("Score moves at depth", depth, ":", scoreMove)
	maxScore := -10000000000
	var bestMove s_ScorePos
	for _, sm := range scoreMove {
		if (sm.score > maxScore) {
			maxScore = sm.score
			bestMove = sm
		}
	}
	// fmt.Println("Best move at depth", depth, ":", bestMove)

	bestMoves := []s_ScorePos{}
	uppurQuartile := maxScore * 75 / 100
	for _, sm := range scoreMove {
		if (sm.score >= uppurQuartile) {
			newTable := table
			putStone(&newTable, bestMove.pos.x, bestMove.pos.y, color)
			captured := capture(&newTable, bestMove.pos.x, bestMove.pos.y, color, color)
			if len(captured) > 0 {
				newTable = updateAvailableMovesAfterCapture(newTable, color, captured)
			}

			availableMovesTable = updateAvailableMoves(availableMovesTable, color, bestMove.pos.x, bestMove.pos.y)
			result := RecursiveSearch(depth - 1, newTable, availableMovesTable, !AIMove, opponentColor(color))
			// fmt.Println("Result of move", sm, "at depth", depth, ":", result)
			if (AIMove) {
				result.score = result.score + bestMove.score * pow10[depth]
			} else {
				result.score = result.score - bestMove.score * pow10[depth]
			}
			result.pos = bestMove.pos
			bestMoves = append(bestMoves, result)
		}
	}
	// fmt.Println("Best Moves at depth", depth, ":", bestMoves)
	maxScore = -1000000000
	for _, bm := range bestMoves {
		if (AIMove && bm.score > maxScore) || (!AIMove && bm.score < maxScore) {
			maxScore = bm.score
			bestMove = bm
		}
	}

	return bestMove
}