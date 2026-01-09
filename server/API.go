package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
	
	"github.com/gin-gonic/gin"
)

type MoveRequest struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	Color int `json:"color"` // 1 = Noir, 2 = Blanc
}

type DebugRequest struct {
	SplitedGoban [][]int `json:"board"`
	CaptiredB   int     `json:"captured_b"`
	CaptiredW   int     `json:"captured_w"`
	ResetGoban  bool    `json:"reset_board"`
}

type GameServer struct {
	mu          sync.Mutex
	goban       s_table
	AIMode		bool
	isBusy      bool
	gameStarted bool
}

func (gs *GameServer) handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (gs *GameServer) handleGetBoard(c *gin.Context) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	
	c.JSON(http.StatusOK, gin.H{
		"board": convertGobanTo2D(&gs.goban.cells),
	})
}

func (gs *GameServer) handleMove(c *gin.Context) {
	var req MoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	gs.mu.Lock()
	defer gs.mu.Unlock()

	turn := playTurn(&gs.goban, req.X, req.Y, uint8(req.Color))
	if turn == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Coup illégal"})
		return
	} else if turn == 1 {
		c.JSON(http.StatusOK, gin.H{
			"status": "Victoire",
			"board":  convertGobanTo2D(&gs.goban.cells),
		})
		return
	}
	
	fmt.Printf("Coup reçu : Joueur %d en (%d, %d)\n", req.Color, req.X, req.Y)
	
	c.JSON(http.StatusOK, gin.H{
		"status": "Coup accepté et IA a répondu",
		"board":  convertGobanTo2D(&gs.goban.cells),
	})
	
	if gs.AIMode {
		aiColor := uint8(2)
		if req.Color == 2 { aiColor = 1 }
		timedAIMove(&gs.goban, aiColor)
	}
}

func (gs *GameServer) handleAISuggest(c *gin.Context) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	time, move := timedAIMoveSuggest(&gs.goban, 1)

	c.JSON(http.StatusOK, gin.H{
		"time μs":	time,
		"x":		move.x,
		"y":		move.y,
	})
}

func (gs *GameServer) handleDebug(c *gin.Context) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	
	var req DebugRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Donnée invalide"})
		return
	}

	if req.ResetGoban {
		gs.goban = s_table{
			size:       gobanWidth,
			captured_b: req.CaptiredB,
			captured_w: req.CaptiredW,
		}
	} else {
		setGoban(&gs.goban, req.SplitedGoban)
		gs.goban.captured_b = req.CaptiredB
		gs.goban.captured_w = req.CaptiredW
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Goban mis à jour",
		"board":   convertGobanTo2D(&gs.goban.cells),
	})
}


func setRouter(router *gin.Engine, gs *GameServer) {
	router.GET("/ping", gs.handlePing)
	router.GET("/board", gs.handleGetBoard)
	router.GET("/ai-suggest", gs.handleAISuggest)
	router.POST("/move", gs.handleMove)
	router.POST("/debug", gs.handleDebug)
}

func setGoban(goban *s_table, newGoban [][]int) {
	for i := 0; i < gobanWidth; i++ {
		for j := 0; j < gobanWidth; j++ {
			index := i*gobanWidth + j
			goban.cells[index] = uint8(newGoban[i][j])
		}
	}
}

func convertGobanTo2D(goban *[gobanSize]uint8) [][]int {

	board2D := make([][]int, gobanWidth)
	for i := range board2D {
		board2D[i] = make([]int, gobanWidth)
	}
	for i := 0; i < gobanWidth; i++ {
		for j := 0; j < gobanWidth; j++ {
			index := i*gobanWidth + j
			board2D[i][j] = int(goban[index])
		}
	}
	return board2D
}

func timedAIMove(goban *s_table, color uint8) int64 {

	start := time.Now()
	move := getAIMove(*goban, color)
	elapsed := time.Since(start)
	playTurn(goban, move.x, move.y, color)

	fmt.Printf("AI move computed in %d μs\n", elapsed.Microseconds())
	return elapsed.Microseconds()
}

func timedAIMoveSuggest(goban *s_table, color uint8) (int64, s_StonesPos) {

	start := time.Now()
	move := getAIMove(*goban, color)
	elapsed := time.Since(start)

	fmt.Printf("AI suggestion computed in %d μs\n", elapsed.Microseconds())
	return	elapsed.Microseconds(), move
}