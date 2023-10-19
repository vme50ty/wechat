package main

import (
	"genchat/message/router"
	"genchat/message/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMySQL()
	r := router.Router()
	r.Run(":8081")
}
