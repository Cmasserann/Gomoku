package main

import (
	"github.com/gin-gonic/gin"
)

const gobanWidth = 19
const gobanSize = gobanWidth * gobanWidth

var goban = s_table{ size: gobanWidth, captured_b: 0, captured_w: 0 }

func main() {
	router := gin.Default()

	setRouter(router)
	// Output: olleh
}
