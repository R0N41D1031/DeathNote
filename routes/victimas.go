package routes

import (
	"github.com/R0N41D1031/Crud-go/controladores"
	"github.com/gin-gonic/gin"
)

func VictimasRoutes(r *gin.Engine) {

	r.GET("/victimas", controladores.ListarVictimasConMuertes)
	r.POST("/victimas", controladores.CrearVictima)

}
