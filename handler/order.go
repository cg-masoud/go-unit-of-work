package handler

import (
	"net/http"
	"strconv"
	"unit-of-work/db"
	"unit-of-work/service"

	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
)

func DeleteOrderHandler(dbConn *sql.DB, orderService *service.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
			return
		}

		// Create Unit of Work
		uow, err := db.NewUnitOfWork(dbConn)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Start transaction
		uow.WithTX()
		if _, err := uow.GetQuerier(c.Request.Context()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx := uow.Context(c)

		// Call the service to delete the order
		if err := orderService.DeleteOrder(ctx, orderID); err != nil {
			if rollbackErr := uow.Rollback(); rollbackErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error rolling back: %v", rollbackErr)})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Commit the transaction
		if err := uow.Commit(); err != nil {
			if rollbackErr := uow.Rollback(); rollbackErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error rolling back: %v", rollbackErr)})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
	}
}
