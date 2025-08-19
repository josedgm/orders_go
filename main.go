package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/josedgm/orders_go/orders"
)

func main() {
	router := gin.Default()

	router.POST("/orders", func(c *gin.Context) {
		var incomingOrders []orders.IncomingOrder

		// Bind JSON input to slice of IncomingOrder
		if err := c.ShouldBindJSON(&incomingOrders); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input JSON: " + err.Error()})
			return
		}

		// Transform orders using the module logic
		customerItemsList, errors := orders.TransformOrders(incomingOrders)
		if len(errors) > 0 {
			for _, err := range errors {
				log.Println("Validation error:", err)
			}
		}

		// Return the transformed customer items as JSON response
		c.JSON(http.StatusOK, customerItemsList)

	})

	router.Run(":8080")
}
