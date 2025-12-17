package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 1. On définit le Goban en global pour qu'il soit accessible partout
var goban [19][19]byte

// 2. Structure pour recevoir les données du JSON (Client -> Serveur)
type MoveRequest struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	Color int `json:"color"` // 1 = Noir, 2 = Blanc
}

func printGoban() {
	// J'ai retiré l'argument car on utilise la variable globale maintenant
	for _, row := range goban {
		for _, cell := range row {
			switch cell {
			case 0:
				fmt.Print(". ")
			case 1:
				fmt.Print("X ")
			case 2:
				fmt.Print("O ")
			}
		}
		fmt.Println() // J'ai remplacé log.Println par fmt.Println pour un affichage plus propre dans le terminal
	}
	fmt.Println("-------------------")
}

// C'est ICI que tu mettras ta logique plus tard (règles, captures, etc.)
// Pour l'instant, elle retourne juste true
func updateGoban(x int, y int, color int) bool {
	// TODO: Coder la logique ici
	// Exemple temporaire pour tester :
	if x >= 0 && x < 19 && y >= 0 && y < 19 {
		goban[x][y] = byte(color)
		return true
	}
	return false
}

func main() {
	// Initialisation basique pour l'exemple
	goban[3][3] = 1
	goban[3][4] = 2

	r := gin.Default()

	// Route de test
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// --- LA REQUÊTE POUR POSER UNE PIERRE ---
	r.POST("/move", func(c *gin.Context) {
		var req MoveRequest

		// 1. On décode le JSON reçu (Bind)
		// Si le JSON n'est pas bon (ex: il manque "x"), ça renvoie une erreur 400
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides (besoin de x, y, color)"})
			return
		}

		// 2. On appelle TA fonction pour mettre à jour le plateau
		validMove := updateGoban(req.X, req.Y, req.Color)

		if !validMove {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Coup invalide"})
			return
		}

		// Debug: on affiche dans le terminal du serveur
		fmt.Printf("Coup reçu : Joueur %d en (%d, %d)\n", req.Color, req.X, req.Y)
		printGoban()

		// 3. RÉPONSE : On renvoie le plateau complet mis à jour au client
		// Gin convertit automatiquement le tableau [19][19]byte en JSON
		c.JSON(http.StatusOK, gin.H{
			"message": "Coup accepté",
			"board":   goban, // Le client recevra le plateau actualisé ici
		})
	})

	// Route pour récupérer le plateau actuel (utile pour le client qui attend)
	r.GET("/board", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"board": goban,
		})
	})

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}