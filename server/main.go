package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const gobanSize = 361 // 19x19 goban

type MoveRequest struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	Color int `json:"color"` // 1 = Noir, 2 = Blanc
}


func updateGoban(goban []byte,x int, y int, color int) bool {

	index := x*gobanSize + y

	if index > -1 && index < gobanSize {
		goban[index] = byte(color)
		return true
	}

	return false
}

func printGoban(goban []byte) {
	fmt.Println("Goban actuel :")
	for i := 0; i < 19; i++ {
		for j := 0; j < 19; j++ {
			index := i*19 + j
			fmt.Printf("%d ", goban[index])
		}
		fmt.Println()
	}
}

func convertGobanTo2D(goban []byte) [][]int {
	board2D := make([][]int, 19)
	for i := range board2D {
		board2D[i] = make([]int, 19)
	}
	for i := 0; i < 19; i++ {
		for j := 0; j < 19; j++ {
			index := i*19 + j
			board2D[i][j] = int(goban[index])
		}
	}
	return board2D
}

func main() {

	goban := make([]byte, gobanSize)

	printGoban(goban)

	r := gin.Default()


	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})


	r.GET("/board", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"board": goban,
		})
	})


	r.POST("/move", func(c *gin.Context) {
		var req MoveRequest


		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides (besoin de x, y, color)"})
			return
		}


		validMove := updateGoban(goban, req.X, req.Y, req.Color)

		if !validMove {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Coup invalide"})
			return
		}


		fmt.Printf("Coup reçu : Joueur %d en (%d, %d)\n", req.Color, req.X, req.Y)


		c.JSON(http.StatusOK, gin.H{
			"message": "Coup accepté",
			"board":   convertGobanTo2D(goban),
		})
	})


	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
