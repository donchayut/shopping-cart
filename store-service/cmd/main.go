package main

import (
	"fmt"
	"log"
	"store-service/cmd/api"
	"store-service/internal/order"
	"store-service/internal/product"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	connection, err := sqlx.Connect("mysql", "sealteam:sckshuhari@(store-database:3306)/toy")
	if err != nil {
		log.Fatalln("cannot connect to tearup", err)
	}
	productRepository := product.ProductRepositoryMySQL{
		DBConnection: connection,
	}
	orderRepository := order.OrderRepositoryMySQL{
		DBConnection: connection,
	}
	orderService := order.OrderService{
		ProductRepository: &productRepository,
		OrderRepository:   &orderRepository,
	}
	storeAPI := api.StoreAPI{
		OrderService: &orderService,
	}

	route := gin.Default()
	route.POST("/api/v1/order", storeAPI.SubmitOrderHandler)

	route.GET("/api/v1/health", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": GetUserNameFromDB(connection),
		})
	})
	log.Fatal(route.Run(":8000"))
}

func GetUserNameFromDB(connection *sqlx.DB) User {
	user := User{}
	err := connection.Get(&user, "SELECT id,name FROM user WHERE ID=1")
	if err != nil {
		fmt.Printf("Get user name from tearup get error : %s", err.Error())
		return User{}
	}
	return user
}

type User struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
