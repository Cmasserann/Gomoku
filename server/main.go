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

type DebugRequest struct {
	SplitedGoban [][]int `json:"board"`
}


func updateGoban(goban []byte,x int, y int, color int) bool {

	index := gobanWidth * y + x

	fmt.Println("Index calculé :", index)

	if index > -1 && index < gobanSize && color > -1 && color < 3 {
		goban[index] = byte(color)
		return true
	}

	return false
}

func setGoban(goban []byte, splitedGoban [][]int) bool {

	if len(splitedGoban) != gobanWidth {
		return false
	}

	for i := 0; i < gobanWidth; i++ {
		if len(splitedGoban[i]) != gobanWidth {
			return false
		}
		for j := 0; j < gobanWidth; j++ {
			index := i*gobanWidth + j
			goban[index] = byte(splitedGoban[i][j])
		}
	}

	return true
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

	r.POST("/debug", func(c *gin.Context) {
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

		setGoban(goban, newGoban)

		c.JSON(http.StatusOK, gin.H{
			"message": "Goban mis à jour",
			"board":   convertGobanTo2D(goban),
		})
	})


	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
