package main

import (
	"database/sql"
	"unit-of-work/handler"
	repo "unit-of-work/repository"
	"unit-of-work/service"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbConn, err := sql.Open("mysql", "root:123456@tcp(172.18.0.2:3306)/go")
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()

	// Create repositories
	orderRepo := repo.NewOrderRepository(dbConn)
	itemRepo := repo.NewItemRepository(dbConn)

	// Create services
	orderService := service.NewOrderService(orderRepo, itemRepo)

	// Set up routes
	r := gin.Default()
	r.DELETE("/orders/:id", handler.DeleteOrderHandler(dbConn, orderService))
	r.Run(":8090")
}
