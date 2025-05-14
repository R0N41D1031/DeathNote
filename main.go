package main

import (
	"github.com/R0N41D1031/Crud-go/db"
	"github.com/R0N41D1031/Crud-go/models"
	"github.com/R0N41D1031/Crud-go/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	db.DBconexion()
	db.DB.AutoMigrate(models.Victima{})
	db.DB.AutoMigrate(models.Muerte{})

	r := gin.Default()

	routes.VictimasRoutes(r)
	routes.MuertesRoutes(r)

	r.Run()
}
