package routes

import (
	"github.com/R0N41D1031/Crud-go/controladores"
	"github.com/gin-gonic/gin"
)

func MuertesRoutes(r *gin.Engine) {

	r.POST("/victimas/muerte", controladores.CrearMuerte)
	r.POST("/muerte/iniciar", controladores.IniciarMuerte)
	r.POST("/muerte/finalizar", controladores.FinalizarMuerte)
}
