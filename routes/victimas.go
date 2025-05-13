package routes

import (
	"encoding/json"
	"net/http"

	"github.com/R0N41D1031/Crud-go/db"
	"github.com/R0N41D1031/Crud-go/models"
	"github.com/gorilla/mux"
)

func GetVictimasHandler(w http.ResponseWriter, r *http.Request) {
	var victimas []models.Victimas
	db.DB.Find(&victimas)
	json.NewEncoder(w).Encode(&victimas)
}

func GetVictimaHandler(w http.ResponseWriter, r *http.Request) {
	var victim models.Victimas
	parametros := mux.Vars(r)
	db.DB.First(&victim, parametros["id"])

	if victim.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Victima no encontrada"))
		return
	}

	json.NewEncoder(w).Encode(&victim)

}

func PostVictimaHandler(w http.ResponseWriter, r *http.Request) {
	var victima models.Victimas
	json.NewDecoder(r.Body).Decode(&victima)

	victimacreada := db.DB.Create(&victima)
	err := victimacreada.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		w.Write([]byte(err.Error()))
	}

	json.NewEncoder(w).Encode(&victima)
}

func DeleteVictimaHandler(w http.ResponseWriter, r *http.Request) {

	var victim models.Victimas
	parametros := mux.Vars(r)
	db.DB.First(&victim, parametros["id"])

	if victim.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Victima no encontrada"))
		return
	}

	db.DB.Delete(&victim)
	w.WriteHeader(http.StatusOK)
}
