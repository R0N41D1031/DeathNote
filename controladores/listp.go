package controladores

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/R0N41D1031/Crud-go/db"
	"github.com/R0N41D1031/Crud-go/models"
	"github.com/gin-gonic/gin"
)

type MuerteTemporal struct {
	Nombre    string
	CreadoEn  time.Time
	Cancelado bool
	Mutex     sync.Mutex
}

var sesiones = make(map[string]*MuerteTemporal)
var mu sync.Mutex

func IniciarMuerte(c *gin.Context) {
	var data struct {
		Nombre string `json:"nombre"`
	}

	if err := c.BindJSON(&data); err != nil || data.Nombre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nombre es requerido"})
		return
	}

	idSesion := data.Nombre + time.Now().Format("150405") // ID único por nombre+hora

	muerte := &MuerteTemporal{
		Nombre:   data.Nombre,
		CreadoEn: time.Now(),
	}

	mu.Lock()
	sesiones[idSesion] = muerte
	mu.Unlock()

	// Inicia el temporizador de 40 segundos
	go func(id string) {
		time.Sleep(40 * time.Second)
		muerte.Mutex.Lock()
		defer muerte.Mutex.Unlock()
		if !muerte.Cancelado {
			// Registrar la muerte por defecto si no se recibe causa
			v := models.Victima{Nombre: muerte.Nombre}
			db.DB.Create(&v)
			m := models.Muerte{Causa: "Desconocida", VictimaID: v.ID}
			db.DB.Create(&m)
			delete(sesiones, id)
		}
	}(idSesion)

	c.JSON(http.StatusOK, gin.H{
		"mensaje":   "Temporizador iniciado. Tienes 40 segundos para ingresar la causa.",
		"sesion_id": idSesion,
	})
}

func FinalizarMuerte(c *gin.Context) {
	var data struct {
		SesionID    string `json:"sesion_id"`
		Causa       string `json:"causa"`
		Descripcion string `json:"descripcion"`
	}

	if err := c.BindJSON(&data); err != nil || data.SesionID == "" || data.Causa == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sesion ID y causa son requeridos"})
		return
	}

	mu.Lock()
	sesion, existe := sesiones[data.SesionID]
	mu.Unlock()

	if !existe {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sesión no encontrada o tiempo expirado"})
		return
	}

	sesion.Mutex.Lock()
	defer sesion.Mutex.Unlock()

	// Marcar como completado y guardar en DB
	sesion.Cancelado = true
	v := models.Victima{Nombre: sesion.Nombre}
	db.DB.Create(&v)
	m := models.Muerte{Causa: data.Causa, Descripcion: data.Descripcion, VictimaID: v.ID}
	db.DB.Create(&m)

	delete(sesiones, data.SesionID)

	c.JSON(http.StatusOK, gin.H{"mensaje": "Muerte registrada correctamente"})
}

var tempVictimas = make(map[uint]time.Time)

func CrearVictima(c *gin.Context) {
	var nueva models.Victima

	if err := c.ShouldBindJSON(&nueva); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	result := db.DB.Create(&nueva)
	if result.Error != nil || nueva.ID == 0 {
		log.Println("[ERROR] No se pudo guardar la víctima:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo guardar la víctima"})
		return
	}

	log.Printf("[TEMPORAL] Guardando víctima con ID %d", nueva.ID)
	tempVictimas[nueva.ID] = time.Now()
	log.Printf("[TEMPORAL] Estado actual: %+v", tempVictimas)

	c.JSON(http.StatusOK, gin.H{
		"message": "Nombre registrado. Tienes 40 segundos para ingresar la causa de la muerte.",
		"id":      nueva.ID,
	})
}

func CrearMuerte(c *gin.Context) {
	var muerte models.Muerte

	if err := c.ShouldBindJSON(&muerte); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	log.Printf("[CAUSA] Recibido ID: %d", muerte.VictimaID)
	log.Printf("[CAUSA] Estado actual antes de buscar: %+v", tempVictimas)

	registro, ok := tempVictimas[muerte.VictimaID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "No se encontró una víctima pendiente"})
		return
	}

	if time.Since(registro) > 40*time.Second {
		delete(tempVictimas, muerte.VictimaID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Se ha superado el tiempo de 40 segundos"})
		return
	}

	result := db.DB.Create(&muerte)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	delete(tempVictimas, muerte.VictimaID)

	c.JSON(http.StatusOK, gin.H{"message": "Causa de muerte registrada correctamente"})
}

func ListarVictimasConMuertes(c *gin.Context) {
	var victimas []models.Victima

	result := db.DB.Preload("Muerte").Find(&victimas)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener las víctimas"})
		return
	}
	c.JSON(http.StatusOK, victimas)
}
