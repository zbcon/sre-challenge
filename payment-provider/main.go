package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := setupRouter()

	err := router.Run(":8082")
	if err != nil {
		fmt.Printf("could not start server: %v", err)
	}
}

func setupRouter() *gin.Engine {
	r := gin.New()
	r.POST("payments/pay", pay)
	return r
}

func pay(c *gin.Context) {
	if c.Request.Body == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "empty body request",
		})
		return
	}

	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "could not retrieve payment body request",
		})
		return
	}

	var req payRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to process payment request parameters",
		})
		return
	}

	if req.Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "could not retrieve payment id",
		})
		return
	}
	if req.Value == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "could not retrieve payment value",
		})
		return
	}
	if req.Currency == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "could not retrieve payment currency",
		})
		return
	}

	fmt.Printf("Payment: %v\n", string(body))

	c.JSON(http.StatusOK, &payResponse{
		Id: req.Id,
	})
}

type payResponse struct {
	Id string `json:"id"`
}

type payRequest struct {
	Id       string  `json:"id"`
	Value    float32 `json:"value"`
	Currency string  `json:"currency"`
}
