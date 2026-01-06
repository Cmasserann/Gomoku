package main

import "testing"

func TestIAMainNoThread(t *testing.T) {
	test_table := s_table{ size: 19, captured_b: 0, captured_w: 0 }
	putStone(&test_table, 5, 5, 1)
	result := IAMainNoThread(test_table, 2)
	if result.x == -1 && result.y == -1 {
		t.Error("Expected IA to play a move, got no move")
	}
	if illegalMove(&test_table, result.x, result.y, 2) {
		t.Error("Expected IA to play a legal move, got illegal move at:", result)
	}
}

func TestIAWinCapture(t *testing.T) {
	test_table := s_table{ size: 19, captured_b: 0, captured_w: 4 }
	putStone(&test_table, 5, 5, 1)
	putStone(&test_table, 6, 5, 1)
	putStone(&test_table, 7, 5, 2)
	result := IAMainNoThread(test_table, 2)
	if result.x == -1 && result.y == -1 {
		t.Error("Expected IA to play a move, got no move")
	}
	if illegalMove(&test_table, result.x, result.y, 2) {
		t.Error("Expected IA to play a legal move, got illegal move at:", result)
	}
	if !(result.x == 4 && result.y == 5) {
		t.Error("Expected IA to play at (4,5) to win by capture, got:", result)
	}
}

func TestIAWinAlign(t *testing.T) {
	test_table := s_table{ size: 19, captured_b: 0, captured_w: 0 }
	putStone(&test_table, 5, 5, 2)
	putStone(&test_table, 6, 5, 2)
	putStone(&test_table, 7, 5, 2)
	putStone(&test_table, 8, 5, 2)
	result := IAMainNoThread(test_table, 2)
	if result.x == -1 && result.y == -1 {
		t.Error("Expected IA to play a move, got no move")
	}
	if illegalMove(&test_table, result.x, result.y, 2) {
		t.Error("Expected IA to play a legal move, got illegal move at:", result)
	}
	if !(result.x == 9 && result.y == 5 || result.x == 4 && result.y == 5) {
		t.Error("Expected IA to play at (9,5) or (4,5) to win by alignment, got:", result)
	}
}

func TestIAWinAfterLoseAlign(t *testing.T) {
	test_table := s_table{ size: 19, captured_b: 4, captured_w: 0 }
	putStone(&test_table, 5, 4, 2)
	putStone(&test_table, 5, 5, 1)
	putStone(&test_table, 5, 6, 1)
	putStone(&test_table, 6, 5, 1)
	putStone(&test_table, 7, 5, 1)
	putStone(&test_table, 8, 5, 1)
	result := IAMainNoThread(test_table, 2)
	if result.x == -1 && result.y == -1 {
		t.Error("Expected IA to play a move, got no move")
	}
	if illegalMove(&test_table, result.x, result.y, 2) {
		t.Error("Expected IA to play a legal move, got illegal move at:", result)
	}
	if !(result.x == 5 && result.y == 7) {
		t.Error("Expected IA to play at (5,7) to block opponent win, got:", result)
	}
}

func TestIADontLoseByCapture(t *testing.T) {
	test_table := s_table{ size: 19, captured_b: 4, captured_w: 0 }
	putStone(&test_table, 4, 5, 2)
	putStone(&test_table, 5, 5, 2)
	putStone(&test_table, 6, 5, 1)
	result := IAMainNoThread(test_table, 2)
	if result.x == -1 && result.y == -1 {
		t.Error("Expected IA to play a move, got no move")
	}
	if illegalMove(&test_table, result.x, result.y, 2) {
		t.Error("Expected IA to play a legal move, got illegal move at:", result)
	}
	if !(result.x == 3 && result.y == 5) {
		t.Error("Expected IA to play at (3,5) to avoid capture loss, got:", result)
	}
}

func TestIADontLoseByCapture2(t *testing.T) {
	test_table := s_table{ size: 19, captured_b: 4, captured_w: 0 }
	putStone(&test_table, 5, 6, 2)
	putStone(&test_table, 5, 7, 2)
	putStone(&test_table, 5, 8, 2)
	putStone(&test_table, 5, 9, 2)
	putStone(&test_table, 4, 5, 2)
	putStone(&test_table, 6, 5, 1)

	result := IAMainNoThread(test_table, 2)
	if result.x == -1 && result.y == -1 {
		t.Error("Expected IA to play a move, got no move")
	}
	if illegalMove(&test_table, result.x, result.y, 2) {
		t.Error("Expected IA to play a legal move, got illegal move at:", result)
	}
	if result.x == 5 && result.y == 5 {
		t.Error("Expected IA to not play at (5,5) to avoid capture loss, got:", result)
	}
	
	putStone(&test_table, 4, 10, 2)
	putStone(&test_table, 6, 10, 1)
	result = IAMainNoThread(test_table, 2)
	if result.x == 5 && result.y == 10 {
		t.Error("Expected IA to not play at (5,10) to avoid capture loss, got:", result)
	}
}