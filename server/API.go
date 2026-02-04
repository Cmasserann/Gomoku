package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
	"crypto/rand"
    "encoding/hex"
	
	"github.com/gin-gonic/gin"
)


type GameServer struct {
	mu          sync.Mutex
	goban       s_table
	AIMode		bool
	localMode	bool
	gameStarted bool
	turn		int
	playerOne	string
	playerTwo	string
}


func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}


func setRouter(router *gin.Engine, gs *GameServer) {

	router.Use(CORSMiddleware())
	router.GET("/status", gs.handleStatus)
	router.GET("/board", gs.handleGetBoard)
	router.POST("/ai-suggest", gs.handleAISuggest)
	router.POST("/create", gs.handleSetGame)
	router.POST("/join", gs.handleInvitation)
	router.POST("/move", gs.handleMove)
	router.POST("/giveUp", gs.handleGiveUp)
	router.POST("/debug", gs.handleDebug)
}


func (gs *GameServer) handleInvitation(c *gin.Context) {

	gs.mu.Lock()
	defer gs.mu.Unlock()
	
	var req struct {
		Token string `json:"token"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}
	
	if req.Token != gs.playerTwo || len(req.Token) != 4 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	
	gs.playerTwo = generateConnectionToken()
	
	c.JSON(http.StatusOK, gin.H{
		"token": gs.playerTwo,
	})
}


func (gs *GameServer) handleStatus(c *gin.Context) {

	gs.mu.Lock()
	c.JSON(http.StatusOK, gin.H{
		"goban_free": !gs.gameStarted,
		"pending_invitation": gs.playerTwo != "" && len(gs.playerTwo) == 4,
	})
	gs.mu.Unlock()
}
	
	
func (gs *GameServer) handleGetBoard(c *gin.Context) {

	gs.mu.Lock()
	c.JSON(http.StatusOK, gin.H{
		"board": convertGobanTo2D(&gs.goban.cells),
		"captured_b": gs.goban.captured_b,
		"captured_w": gs.goban.captured_w,
		"goban_free": !gs.gameStarted,
		"turn": gs.turn,
	})
	gs.mu.Unlock()
}


func (gs *GameServer) handleMove(c *gin.Context) {
	
	var req struct {
		X     int `json:"x"`
		Y     int `json:"y"`
		Token string `json:"token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}
	
	gs.mu.Lock()
	defer gs.mu.Unlock()
	
	if gs.gameStarted == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No game in progress"})
		return
	}
	
	player := gs.checkToken(req.Token)

	
	if player == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		fmt.Printf("Invalid token received: %s\n token needed: %s\n", req.Token ,gs.playerOne)
		return
	}
	
	if gs.localMode == false && player == gs.turn % 2 + 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not your turn"})
		return
	}

	play := playTurn(&gs.goban, req.X, req.Y, uint8((gs.turn + 1) % 2 + 1))
	if play == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Illegal move"})
		return
	} else if play == 1 {
		c.JSON(http.StatusOK, gin.H{
			"winner": player,
			"board":  convertGobanTo2D(&gs.goban.cells),
		})
		gs.gameStarted = false
		gs.playerOne = ""
		gs.playerTwo = ""
		return
	}
	
	gs.turn += 1

	if gs.AIMode {
		time, turn := timedAIMove(&gs.goban, 2)
		
		if turn == 1 {
			c.JSON(http.StatusOK, gin.H{
				"winner": 2,
				"board":  convertGobanTo2D(&gs.goban.cells),
				"time_us": time,
			})
			gs.gameStarted = false
			gs.playerOne = ""
			gs.playerTwo = ""
			return
		}

		gs.turn += 1
		c.JSON(http.StatusOK, gin.H{
			"board":  convertGobanTo2D(&gs.goban.cells),
			"turn":	gs.turn,
			"time": time,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"board":  convertGobanTo2D(&gs.goban.cells),
			"turn":	gs.turn,
		})
	}
	
	
}


func (gs *GameServer) handleSetGame(c *gin.Context) {

	gs.mu.Lock()
	defer gs.mu.Unlock()

	if gs.gameStarted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game already in progress"})
		return
	}

	var req struct {
		AIMode		bool	`json:"ai_mode"`
		LocalMode	bool	`json:"local_mode"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	gs.AIMode = req.AIMode
	if gs.AIMode == true {
		gs.localMode = false
	} else {
		gs.localMode = req.LocalMode
	}
	gs.gameStarted = true
	gs.playerOne = generateConnectionToken()
	if gs.localMode == true || gs.AIMode == true {
		gs.playerTwo = ""
	} else {
		gs.playerTwo = generateInvitationToken()
	}
	gs.turn = 1
	ressetGoban(&gs.goban)

	c.JSON(http.StatusOK, gin.H{
		"ai_mode":     gs.AIMode,
		"local_mode":  gs.localMode,
		"player_one":  gs.playerOne,
		"player_two":  gs.playerTwo,
	})
}


func (gs *GameServer) handleAISuggest(c *gin.Context) {

	var req struct {
		Token string `json:"token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	gs.mu.Lock()
	defer gs.mu.Unlock()

	player := gs.checkToken(req.Token)
	
	if player == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	
	time, move := timedAIMoveSuggest(&gs.goban, 1)

	c.JSON(http.StatusOK, gin.H{
		"time_us":	time,
		"x":		move.x,
		"y":		move.y,
	})
}


func (gs *GameServer) handleDebug(c *gin.Context) {
	
	var req struct {
		SplitedGoban [][]int `json:"board"`
		CaptiredB   int     `json:"captured_b"`
		CaptiredW   int     `json:"captured_w"`
		ResetGoban  bool    `json:"reset_board"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}
	
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if req.ResetGoban {
		gs.goban = s_table{
			size:       gobanWidth,
			captured_b: req.CaptiredB,
			captured_w: req.CaptiredW,
		}
		gs.turn = 1
		gs.gameStarted = false
		gs.playerOne = ""
		gs.playerTwo = ""
	} else {

		if len(req.SplitedGoban) == gobanWidth && len(req.SplitedGoban[0]) == gobanWidth {
			setGoban(&gs.goban, req.SplitedGoban)
		}

		if req.CaptiredB >= 0 {
			gs.goban.captured_b = req.CaptiredB
		}
		
		if req.CaptiredW >= 0 {
			gs.goban.captured_w = req.CaptiredW
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"board":   convertGobanTo2D(&gs.goban.cells),
		"captured_b": gs.goban.captured_b,
		"captured_w": gs.goban.captured_w,
		"turn": gs.turn,
		"token_player_one": gs.playerOne,
		"token_player_two": gs.playerTwo,
	})
}


func (gs *GameServer) checkToken(token string) int {
    if len(token) != 16 { return 0 }
    if token == gs.playerOne { return 1 }
    if token == gs.playerTwo { return 2 }
    return 0
}


func (gs *GameServer) handleGiveUp(c *gin.Context) {

	var req struct {
		Token string `json:"token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	gs.mu.Lock()
	defer gs.mu.Unlock()

	player := gs.checkToken(req.Token)
	
	if player == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	gs.gameStarted = false
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Game over.",
	})
}


func setGoban(goban *s_table, newGoban [][]int) {

	for i := 0; i < gobanWidth; i++ {
		for j := 0; j < gobanWidth; j++ {
			index := i*gobanWidth + j
			goban.cells[index] = uint8(newGoban[i][j])
		}
	}
}


func ressetGoban(goban *s_table) {

	*goban = s_table{
		size:       gobanWidth,
		captured_b: 0,
		captured_w: 0,
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


func timedAIMove(goban *s_table, color uint8) (int64, int) {

	start := time.Now()
	move := getAIMove(*goban, color)
	elapsed := time.Since(start)
	turn := playTurn(goban, move.x, move.y, color)

	fmt.Printf("AI move computed in %d μs\n", elapsed.Microseconds())
	return elapsed.Microseconds(), turn
}


func timedAIMoveSuggest(goban *s_table, color uint8) (int64, s_StonesPos) {

	start := time.Now()
	move := getAIMove(*goban, color)
	elapsed := time.Since(start)

	fmt.Printf("AI suggestion computed in %d μs\n", elapsed.Microseconds())
	return	elapsed.Microseconds(), move
}


func generateConnectionToken() string {

	b := make([]byte, 8)
    rand.Read(b)
	fmt.Printf("Generated token: %x\n", b)
    return hex.EncodeToString(b)
}


func generateInvitationToken() string {

	b := make([]byte, 2)
    rand.Read(b)
	fmt.Printf("Generated token: %x\n", b)
    return hex.EncodeToString(b)
}