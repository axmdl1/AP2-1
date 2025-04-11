package main

import (
	userservice "AP-1/pb/userService"
	"AP-1/userService/internal/handler"
	"AP-1/userService/internal/repository"
	"AP-1/userService/internal/usecase"
	"AP-1/userService/pkg/mongo"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	db := mongo.ConnectMongoDB("mongodb://localhost:27017/userServiceDB", "userServiceDB")

	userRepo := repository.NewUserRepository(db)
	useCase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserServiceServer(useCase)

	lis, err := net.Listen("tcp", ":1004")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	userservice.RegisterUserServiceServer(grpcServer, userHandler)
	log.Println("Server started on port 1004")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	//router
	/*router := gin.Default()
	routes.SetupRoutes(router, userHandler)

	router.LoadHTMLGlob("userService/ui/*")
	router.Static("/ui", "./ui")

	if err := router.Run(":1004"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}*/
}
