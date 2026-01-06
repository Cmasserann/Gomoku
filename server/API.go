package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MoveRequest struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	Color int `json:"color"` // 1 = Noir, 2 = Blanc
}

type DebugRequest struct {
	SplitedGoban [][]int `json:"board"`
}

func setRouter(router *gin.Engine) {

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})


	router.GET("/board", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"board": convertGobanTo2D(&goban.cells),
		})
	})


	router.POST("/move", func(c *gin.Context) {
		var req MoveRequest


		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides (besoin de x, y, color)"})
			return
		}


		validMove := playTurn(&goban, req.X, req.Y, uint8(req.Color))

		if validMove == -1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Coup invalide"})
			return
		}


		fmt.Printf("Coup reçu : Joueur %d en (%d, %d)\n", req.Color, req.X, req.Y)

		if validMove == 1 {
			c.JSON(http.StatusOK, gin.H{
				"message": "Coup accepté - Victoire !",
				"board":   convertGobanTo2D(&goban.cells),
			})
			goban = s_table{size: gobanWidth, captured_b: 0, captured_w: 0}
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Coup accepté",
			"board":   convertGobanTo2D(&goban.cells),
		})
	})

	router.POST("/debug", func(c *gin.Context) {
		var req DebugRequest


		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Donnée invalide (besoin de splited_goban)"})
			return
		}

		newGoban := req.SplitedGoban

		fmt.Println("Nouveau goban reçu :")
		for i := 0; i < gobanWidth; i++ {
			for j := 0; j < gobanWidth; j++ {
				fmt.Printf("%d ", newGoban[i][j])
			}
			fmt.Println()
		}

		setGoban(&goban, newGoban)

		c.JSON(http.StatusOK, gin.H{
			"message": "Goban mis à jour",
			"board":   convertGobanTo2D(&goban.cells),
		})
	})


	if err := router.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
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