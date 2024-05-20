package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error

type User struct {
	ID      uint    `json:"id" gorm:"primary_key"`
	Balance float64 `json:"balance"`
}

type Transaction struct {
	ID              uint    `json:"id" gorm:"primary_key"`
	UserID          uint    `json:"user_id"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
}

func main() {
	dsn := "host=localhost user=user dbname=user_balance_db sslmode=disable password=password"
	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&User{}, &Transaction{})

	router := gin.Default()

	router.POST("/add-funds", addFunds)
	router.POST("/reserve-funds", reserveFunds)
	router.POST("/recognize-revenue", recognizeRevenue)
	router.GET("/balance/:id", getBalance)

	router.Run(":8080")
}

func addFunds(c *gin.Context) {
	var input struct {
		UserID uint    `json:"user_id"`
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	db.FirstOrCreate(&user, User{ID: input.UserID})
	user.Balance += input.Amount
	db.Save(&user)

	transaction := Transaction{UserID: user.ID, Amount: input.Amount, TransactionType: "credit"}
	db.Create(&transaction)

	c.JSON(http.StatusOK, user)
}

func reserveFunds(c *gin.Context) {
	var input struct {
		UserID    uint    `json:"user_id"`
		Amount    float64 `json:"amount"`
		ServiceID uint    `json:"service_id"`
		OrderID   uint    `json:"order_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	if db.First(&user, input.UserID).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if user.Balance < input.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient funds"})
		return
	}

	user.Balance -= input.Amount
	db.Save(&user)

	transaction := Transaction{UserID: user.ID, Amount: input.Amount, TransactionType: "reserve"}
	db.Create(&transaction)

	c.JSON(http.StatusOK, user)
}

func recognizeRevenue(c *gin.Context) {
	var input struct {
		UserID    uint    `json:"user_id"`
		Amount    float64 `json:"amount"`
		ServiceID uint    `json:"service_id"`
		OrderID   uint    `json:"order_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction := Transaction{UserID: input.UserID, Amount: -input.Amount, TransactionType: "debit"}
	db.Create(&transaction)

	c.JSON(http.StatusOK, gin.H{"message": "revenue recognized"})
}

func getBalance(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var user User
	if db.First(&user, userID).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
