package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	pb "grpc-project/gen/proto"
	"log"
	"net/http"
	"strconv"
)

var client pb.TestApiClient

func main() {

	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Println(err)
	}

	client = pb.NewTestApiClient(conn)

	g := gin.Default()

	g.GET("/show/:param", echo)

	g.GET("/product/:id", getProduct)

	g.GET("/products", getProducts)

	if err := g.Run(":8081"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}

func echo(ctx *gin.Context) {
	param := ctx.Param("param")
	req := &pb.ResponseRequest{Msg: param}
	if response, err := client.Echo(ctx, req); err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(response.Msg),
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func getProduct(ctx *gin.Context) {
	param := ctx.Param("id")
	idParam, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req := &pb.ProductId{
		Id: int32(idParam),
	}
	if response, err := client.GetProduct(ctx, req); err == nil {
		ctx.JSON(http.StatusOK, response)
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func getProducts(ctx *gin.Context) {
	req := &pb.Empty{}
	if response, err := client.GetProducts(ctx, req); err == nil {
		ctx.JSON(http.StatusOK, response)
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
