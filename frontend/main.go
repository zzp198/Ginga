package frontend

import "github.com/gin-gonic/gin"

func Server(addr ...string) error {

	r := gin.New()

	return r.Run(addr...)
}
