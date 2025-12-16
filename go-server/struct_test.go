package main

import "testing"

// test putStone
func TestPutStone(t *testing.T) {
	test_table := s_table{ size: 5, captured_b: 0, captured_w: 0 }
	putStone(&test_table, 2, 2, "b")
	putStone(&test_table, 3, 3, "w")
	expectedCells := [19*19]string{}
	expectedCells[2*19+2] = "b"
	expectedCells[3*19+3] = "w"
	for i, cell := range expectedCells {
		if test_table.cells[i] != cell {
			t.Errorf("Expected cell %d to be %s, got %s", i, cell, test_table.cells[i])
		}
	}
}

// test printTable (just to ensure it runs without error)
func TestPrintTable(t *testing.T) {
	test_table := s_table{ size: 5, captured_b: 0, captured_w: 0 }
	putStone(&test_table, 0, 0, "b")
	putStone(&test_table, 4, 4, "w")
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
		putStone(&test_table, i, 0, "b")
		putStone(&test_table, 0, i, "b")
		putStone(&test_table, i, i, "b")
		putStone(&test_table, i, 4-i, "b")
	}

	winPoint := verifWinPoint(&test_table, 1, 0, "b")
	if winPoint.x_start != 0 || winPoint.y_start != 0 ||
		winPoint.x_end != 4 || winPoint.y_end != 0 {
		t.Error("Expected win point from (0,0) to (4,0), got (", winPoint.x_start, ",", winPoint.y_start, ") to (", winPoint.x_end, ",", winPoint.y_end, ")")
	}
	winPoint = verifWinPoint(&test_table, 0, 1, "b")
	if winPoint.x_start != 0 || winPoint.y_start != 0 ||
		winPoint.x_end != 0 || winPoint.y_end != 4 {
		t.Error("Expected win point from (0,0) to (0,4), got (", winPoint.x_start, ",", winPoint.y_start, ") to (", winPoint.x_end, ",", winPoint.y_end, ")")
	}
	winPoint = verifWinPoint(&test_table, 1, 1, "b")
	if winPoint.x_start != 0 || winPoint.y_start != 0 ||
		winPoint.x_end != 4 || winPoint.y_end != 4 {
		t.Error("Expected win point from (0,0) to (4,4), got (", winPoint.x_start, ",", winPoint.y_start, ") to (", winPoint.x_end, ",", winPoint.y_end, ")")
	}
	winPoint = verifWinPoint(&test_table, 3, 1, "b")
	if winPoint.x_start != 0 || winPoint.y_start != 4 ||
		winPoint.x_end != 4 || winPoint.y_end != 0 {
		t.Error("Expected win point from (0,4) to (4,0), got (", winPoint.x_start, ",", winPoint.y_start, ") to (", winPoint.x_end, ",", winPoint.y_end, ")")
	}
	winPoint = verifWinPoint(&test_table, 8, 12, "b")
	if winPoint.x_start != -1 || winPoint.y_start != -1 ||
		winPoint.x_end != -1 || winPoint.y_end != -1 {
		t.Error("Expected no win point, got (", winPoint.x_start, ",", winPoint.y_start, ") to (", winPoint.x_end, ",", winPoint.y_end, ")")
	}

}
// test capturePossibe
func TestCapturePossible(t *testing.T) {
	test_table := s_table{ size: 19, captured_b: 0, captured_w: 0 }
	putStone(&test_table, 10, 10, "b")
	putStone(&test_table, 11, 10, "w")
	putStone(&test_table, 12, 10, "w")
	if !capturePossibe(&test_table, 11, 10, "w") {
		t.Error("Expected capturePossibe to return true for possible capture")
	}
	if capturePossibe(&test_table, 10, 10, "b") {
		t.Error("Expected capturePossibe to return false for no capture")
	}
	if !capturePossibe(&test_table, 12, 10, "w") {
		t.Error("Expected capturePossibe to return true for possible capture")
	}
}
