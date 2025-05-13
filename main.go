package main

import (
	"net/http"

	"github.com/R0N41D1031/Crud-go/db"
	"github.com/R0N41D1031/Crud-go/models"
	"github.com/R0N41D1031/Crud-go/routes"
	"github.com/gorilla/mux"
)

func main() {

	db.DBconexion()

	db.DB.AutoMigrate(models.Victimas{})
	db.DB.AutoMigrate(models.Muertes{})

	r := mux.NewRouter()

	r.HandleFunc("/", routes.HomeHandler)

	r.HandleFunc("/victimas", routes.GetVictimasHandler).Methods("GET")
	r.HandleFunc("/victimas/{id}", routes.GetVictimaHandler).Methods("GET")
	r.HandleFunc("/victimas", routes.PostVictimaHandler).Methods("POST")
	r.HandleFunc("/victimas/{id}", routes.DeleteVictimaHandler).Methods("DELETE")

	http.ListenAndServe(":8000", r)
}
