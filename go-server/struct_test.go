package main

import "testing"

// test putStone
func TestPutStone(t *testing.T) {
	test_table := s_table{ size: 5, captured_b: 0, captured_w: 0 }
	putStone(&test_table, 2, 2, "b")
	putStone(&test_table, 3, 3, "w")
	expectedCells := [5*5]string{}
	expectedCells[2*5+2] = "b"
	expectedCells[3*5+3] = "w"
	for i, cell := range expectedCells {
		if test_table.cells[i] != cell {
			t.Errorf("Expected cell %d to be %s, got %s", i, cell, test_table.cells[i])
		}
	}

	if putStone(&test_table, 5, 5, "b") {
		t.Error("Expected putStone to return false for out of bounds")
	}
}

// test printTable (just to ensure it runs without error)
func TestPrintTable(t *testing.T) {
	test_table := s_table{ size: 5, captured_b: 0, captured_w: 0 }
	putStone(&test_table, 0, 0, "b")
	putStone(&test_table, 1, 1, "w")
	printTable(&test_table)
}


// test getCapturedStones
func TestGetCapturedStones(t *testing.T) {
	test_table := s_table{ size: 19, captured_b: 8, captured_w: 5 }
	capturedB := getCapturedStones(&test_table, "b")
	if capturedB != 8 {
		t.Errorf("Expected 8 captured black stones, got %d", capturedB)
	}
	capturedW := getCapturedStones(&test_table, "w")
	if capturedW != 5 {
		t.Errorf("Expected 5 captured white stones, got %d", capturedW)
	}
}

// test tableToDict
func TestTableToDict(t *testing.T) {
	test_table := s_table{ size: 19, captured_b: 0, captured_w: 0 }
	putStone(&test_table, 0, 0, "b")
	putStone(&test_table, 1, 1, "w")
	dict := tableToDict(&test_table)
	if len(dict.b) != 1 || dict.b[0] != (s_StonesPos{x: 0, y: 0}) {
		t.Error("Expected one black stone at (0,0)")
	}
	if len(dict.w) != 1 || dict.w[0] != (s_StonesPos{x: 1, y: 1}) {
		t.Error("Expected one white stone at (1,1)")
	}
}

// test verifWinPoint
func TestVerifWinPoint(t *testing.T) {
	test_table := s_table{ size: 19, captured_b: 0, captured_w: 0 }
	for i := 0; i < 5; i++ {
		test_table.cells[0*19+i] = "b"
		test_table.cells[i*19+0] = "b"
		test_table.cells[i*19+i] = "b"
		test_table.cells[i*19+(4-i)] = "b"
	}

	if !verifWinPoint(&test_table, 1, 0, "b") {
		t.Error("Expected win by horizontal line for black")
	}
	if !verifWinPoint(&test_table, 0, 1, "b") {
		t.Error("Expected win by vertical line for black")
	}
	if !verifWinPoint(&test_table, 1, 1, "b") {
		t.Error("Expected win by diagonal \\ line for black")
	}
	if !verifWinPoint(&test_table, 3, 1, "b") {
		t.Error("Expected win by diagonal / line for black")
	}
	if verifWinPoint(&test_table, 8, 12, "b") {
		t.Error("Expected no win point for black at (8,12)")
	}
}

// test verifCapturePossible
func TestVerifCapturePossible(t *testing.T) {
	test_table := s_table{ size: 19, captured_b: 0, captured_w: 0 }
	putStone(&test_table, 10, 10, "b")
	putStone(&test_table, 11, 10, "w")
	putStone(&test_table, 12, 10, "w")
	result := verifCapturePossible(&test_table, "b")
	if result != (s_StonesPos{x: 13, y: 10}) {
		t.Error("Expected verifCapturePossible to return position (13,10) for possible capture, got: ", result)
	}
	result = verifCapturePossible(&test_table, "w")
	if result != (s_StonesPos{x: -1, y: -1}) {
		t.Error("Expected verifCapturePossible to return (-1,-1) for no capture")
	}
}

// test capture
func TestCapture(t *testing.T) {
	test_table := s_table{ size: 19, captured_b: 0, captured_w: 0 }
	putStone(&test_table, 10, 10, "b")
	putStone(&test_table, 11, 10, "w")
	putStone(&test_table, 12, 10, "w")
	putStone(&test_table, 13, 10, "b")
	capture(&test_table, 10, 10, "b", "b")
	if test_table.cells[11*19+10] != "" || test_table.cells[12*19+10] != "" {
		t.Error("Expected stones at (11,10) and (12,10) to be captured")
	}
	if test_table.captured_b != 1 {
		t.Error("Expected captured_b to be 1, got", test_table.captured_b)
	}
	putStone(&test_table, 10, 11, "b")
	putStone(&test_table, 10, 12, "w")
	putStone(&test_table, 10, 9, "w")
	capture(&test_table, 10, 9, "w", "w")
	if test_table.cells[10*19+10] != "" || test_table.cells[10*19+11] != "" {
		t.Error("Expected stones at (10,10) and (10,11) to be captured")
	}
	if test_table.captured_w != 1 {
		t.Error("Expected captured_w to be 1, got", test_table.captured_w)
	}
}

// test freeThree
func TestFreeThree(t *testing.T) {
	test_table := s_table{ size: 19, captured_b: 0, captured_w: 0 }
	putStone(&test_table, 10, 8, "w")
	putStone(&test_table, 10, 9, "w")
	putStone(&test_table, 9, 10, "w")
	putStone(&test_table, 7, 10, "w")
	result := freeThrees(&test_table, 10, 10, "w")
	if result == 0 {
		t.Error("Expected freeThree to not return 0 for white, got: ", result)
	}
	result = freeThrees(&test_table, 10, 10, "b")
	if result != 0 {
		t.Error("Expected freeThree to 0 for black, got: ", result)
	}

	putStone(&test_table, 1, 0, "w")
	putStone(&test_table, 1, 2, "w")
	putStone(&test_table, 0, 0, "w")
	putStone(&test_table, 3, 3, "w")
	putStone(&test_table, 3, 1, "w")
	putStone(&test_table, 4, 1, "w")
	result = freeThrees(&test_table, 1, 1, "w")
	if result == 0 {
		t.Error("Expected freeThree to not return 0 for white at (0,0), got: ", result)
	}
}

// test illegalMove
func TestIllegalMove(t *testing.T) {
	test_table := s_table{ size: 19, captured_b: 0, captured_w: 0 }
	putStone(&test_table, 10, 8, "w")
	putStone(&test_table, 10, 9, "w")
	putStone(&test_table, 9, 10, "w")
	putStone(&test_table, 8, 10, "w")
	if !illegalMove(&test_table, 10, 10, "w") {
		t.Error("Expected move at (10,10) for white to be illegal (double free three)")
	}

	if illegalMove(&test_table, 10, 10, "b") {
		t.Error("Expected move at (10,10) for black to be legal")
	}

	if !illegalMove(&test_table, -1, 0, "b") {
		t.Error("Expected move at (-1,0) for black to be illegal (out of bounds)")
	}

	if !illegalMove(&test_table, 0, 19, "w") {
		t.Error("Expected move at (0,19) for white to be illegal (out of bounds)")
	}

	if !illegalMove(&test_table, 10, 8, "w") {
		t.Error("Expected move at (10,8) for white to be illegal (cell occupied)")
	}

	if !illegalMove(&test_table, 0, 0, "T") {
		t.Error("Expected move at (0,0) for invalid color 'T' to be illegal")
	}
}