package main

import (
	"distribuidos/tarea-1/server/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Levantar servidor Gin
	router := gin.Default()

	router.POST("/api/reserva", routes.PostReserva)
	router.GET("/api/reserva", routes.GetReserva)
	router.GET("/api/vuelo", routes.GetVuelo)
	router.POST("/api/vuelo", routes.PostVuelo)

	router.Run("localhost:5000")
}
