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


type MoveRequest struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	Token string `json:"token"`
}


type DebugRequest struct {
	SplitedGoban [][]int `json:"board"`
	CaptiredB   int     `json:"captured_b"`
	CaptiredW   int     `json:"captured_w"`
	ResetGoban  bool    `json:"reset_board"`
}


type SetGameRequest struct {
	AIMode		bool	`json:"ai_mode"`
	LocalMode	bool	`json:"local_mode"`
}


type GameServer struct {
	mu          sync.Mutex
	goban       s_table
	AIMode		bool
	localMode	bool
	gameStarted bool
	isBusy      bool
	turn		int
	playerOne	string
	playerTwo	string
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
		"message":    "Invitation accepted",
		"token": gs.playerTwo,
	})
}


func (gs *GameServer) handleStatus(c *gin.Context) {
	gs.mu.Lock()
	if gs.gameStarted {
		c.JSON(http.StatusOK, gin.H{"goban status": "Locked"})
	} else {
		c.JSON(http.StatusOK, gin.H{"goban status": "Free"})
	}
	gs.mu.Unlock()
}


func (gs *GameServer) handleGetBoard(c *gin.Context) {
	gs.mu.Lock()
	c.JSON(http.StatusOK, gin.H{
		"board": convertGobanTo2D(&gs.goban.cells),
		"captured_b": gs.goban.captured_b,
		"captured_w": gs.goban.captured_w,
		"turn": gs.turn,
	})
	gs.mu.Unlock()
}


func (gs *GameServer) handleMove(c *gin.Context) {

	var req MoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	gs.mu.Lock()
	
	if gs.gameStarted == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No game in progress"})
		return
	}
	
	gs.mu.Unlock()
	
	player := gs.checktoken(c, req.Token)

	gs.mu.Lock()
	defer gs.mu.Unlock()

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
			"status": "Victoire",
			"board":  convertGobanTo2D(&gs.goban.cells),
		})
		gs.gameStarted = false
		gs.turn = 1
		return
	}
	
	fmt.Printf("Coup reçu : Joueur %d en (%d, %d)\n", player, req.X, req.Y)
	gs.turn += 1

	if gs.AIMode {
		time, turn := timedAIMove(&gs.goban, 2)
		
		if turn == 1 {
			c.JSON(http.StatusOK, gin.H{
				"status": "Victoire de l'IA",
				"board":  convertGobanTo2D(&gs.goban.cells),
				"time μs": time,
			})
			gs.goban = s_table{
				size:       gobanWidth,
				captured_b: 0,
				captured_w: 0,
			}
			gs.turn = 1
			gs.gameStarted = false
			gs.playerOne = ""
			gs.playerTwo = ""
			return
		}

		gs.turn += 1
		c.JSON(http.StatusOK, gin.H{
			"board":  convertGobanTo2D(&gs.goban.cells),
			"turn":	gs.turn,
			"time μs": time,
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Une partie est déjà en cours"})
		return
	}

	var req SetGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Donnée invalide"})
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
	gs.playerTwo = generateInvitationToken()
	gs.turn = 1

	c.JSON(http.StatusOK, gin.H{
		"ai_mode":     gs.AIMode,
		"player_one":  gs.playerOne,
		"player_two":  gs.playerTwo,
	})
}


func (gs *GameServer) handleAISuggest(c *gin.Context) {

	if gs.checktoken(c, c.Query("token")) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide"})
		return
	}

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
		gs.turn = 1
		gs.gameStarted = false
		gs.playerOne = ""
		gs.playerTwo = ""
	} else {
		setGoban(&gs.goban, req.SplitedGoban)
		gs.goban.captured_b = req.CaptiredB
		gs.goban.captured_w = req.CaptiredW
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


func (gs *GameServer) checktoken(c *gin.Context, token string, ) int {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	
	if len(token) != 16 {
		return 0
	} else if token == gs.playerOne {
		return 1
	} else if token == gs.playerTwo{
		return 2
	}
	return 0
}


func setRouter(router *gin.Engine, gs *GameServer) {
	router.GET("/status", gs.handleStatus)
	router.GET("/board", gs.handleGetBoard)
	router.GET("/ai-suggest", gs.handleAISuggest)
	router.POST("/create", gs.handleSetGame)
	router.POST("/join", gs.handleInvitation)
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