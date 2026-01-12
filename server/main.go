package main

import (
	"github.com/gin-gonic/gin"
)

const gobanWidth = 19
const gobanSize = gobanWidth * gobanWidth


func main() {
    router := gin.Default()

    server := &GameServer{
        goban: s_table{
            size:       gobanWidth,
            captured_b: 0,
            captured_w: 0,
        },
		AIMode:			false,
        localMode:	    false,
        gameStarted:	false,
        turn:			1,
        isBusy:			false,
    }

    setRouter(router, server)

    router.Run(":8080")
}