package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const gobanWidth = 19
const gobanSize = gobanWidth * gobanWidth

type MoveRequest struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	Color int `json:"color"` // 1 = Noir, 2 = Blanc
}


func updateGoban(goban []byte,x int, y int, color int) bool {

	index := gobanWidth * y + x

	fmt.Println("Index calculé :", index)

	if index > -1 && index < gobanSize && (color == 1 || color == 2) && goban[index] == 0 {
		goban[index] = byte(color)
		return true
	}

	return false
}

func printGoban(goban []byte) {
	fmt.Println("Goban actuel :")
	for i := 0; i < gobanWidth; i++ {
		for j := 0; j < gobanWidth; j++ {
			index := i*gobanWidth + j
			fmt.Printf("%d ", goban[index])
		}
		fmt.Println()
	}
}

func convertGobanTo2D(goban []byte) [][]int {
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

func main() {

	goban := make([]byte, gobanSize)

	printGoban(goban)

	r := gin.Default()


	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})


	r.GET("/board", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"board": convertGobanTo2D(goban),
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
