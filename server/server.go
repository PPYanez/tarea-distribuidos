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
	router.PUT("/api/reserva", routes.UpdateReserva)

	router.GET("/api/vuelo", routes.GetVuelo)
	router.POST("/api/vuelo", routes.PostVuelo)
	router.PUT("/api/vuelo", routes.UpdateStock)

	router.GET("/api/estadisticas", routes.GetEstadisticas)

	router.Run("localhost:5000")
}
