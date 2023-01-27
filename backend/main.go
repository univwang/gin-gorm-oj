package main

import (
	"gin_oj/router"
)

func main() {

	r := router.Router()
	r.Run(":3000")
}
