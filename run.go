package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zzp198/Ginga/lib"
)

func main() {
	r := gin.Default()

	fmt.Println(testlib.Add(1, 2))

	r.Run("0.0.0.0:80")
}
