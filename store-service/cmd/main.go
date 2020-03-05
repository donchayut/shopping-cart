package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"store-service/cmd/api"
	"store-service/internal/order"
	"store-service/internal/payment"
	"store-service/internal/product"
	"store-service/internal/shipping"
	"time"
)

func main() {
	connection, err := sqlx.Connect("mysql", "sealteam:sckshuhari@(store-database:3306)/toy")
	if err != nil {
		log.Fatalln("cannot connect to database", err)
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
	bankGateway := payment.BankGateway{
		BankEndpoint: "",
	}
	shippingGateway := shipping.ShippingGateway{
		KerryEndpoint: "",
	}
	paymentService := payment.PaymentService{
		BankGateway:       &bankGateway,
		ShippingGateway:   &shippingGateway,
		OrderRepository:   orderRepository,
		ProductRepository: productRepository,
		Time:              time.Now(),
	}
	storeAPI := api.StoreAPI{
		OrderService: &orderService,
	}
	paymentAPI := api.PaymentAPI{
		PaymentService: paymentService,
	}

	route := gin.Default()
	route.POST("/api/v1/order", storeAPI.SubmitOrderHandler)
	route.POST("/api/v1/confirmPayment", paymentAPI.ConfirmPaymentHandler)

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
