package handler

import (
	"context"
	"log"

	"AP-1/inventoryService/internal/entity"
	"AP-1/inventoryService/internal/usecase"
	pb "AP-1/pb/inventoryService"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InventoryServiceServer struct {
	pb.UnimplementedInventoryServiceServer
	usecase usecase.ProductUsecase
}

func NewInventoryServiceServer(u usecase.ProductUsecase) *InventoryServiceServer {
	return &InventoryServiceServer{usecase: u}
}

func (s *InventoryServiceServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	log.Printf("CreateProduct request: %+v", req)
	product := entity.Product{
		Name:        req.Product.Name,
		Category:    req.Product.Category,
		Description: req.Product.Description,
		Price:       req.Product.Price,
		Stock:       int(req.Product.Stock),
	}

	err := s.usecase.CreateProduct(ctx, &product)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	responseProduct := &pb.Product{
		Id:          product.ID,
		Name:        product.Name,
		Category:    product.Category,
		Description: product.Description,
		Price:       product.Price,
		Stock:       int32(product.Stock),
	}
	return &pb.CreateProductResponse{
		Product: responseProduct,
		Message: "Product created successfully",
	}, nil
}

func (s *InventoryServiceServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	log.Printf("GetProduct request: %+v", req)
	product, err := s.usecase.GetProductByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	responseProduct := &pb.Product{
		Id:          product.ID,
		Name:        product.Name,
		Category:    product.Category,
		Description: product.Description,
		Price:       product.Price,
		Stock:       int32(product.Stock),
	}
	return &pb.GetProductResponse{
		Product: responseProduct,
	}, nil
}

func (s *InventoryServiceServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	log.Printf("ListProducts request: %+v", req)
	products, err := s.usecase.ListProducts(ctx, req.Skip, req.Limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	var responseProducts []*pb.Product
	for _, product := range products {
		responseProducts = append(responseProducts, &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Category:    product.Category,
			Description: product.Description,
			Price:       product.Price,
			Stock:       int32(product.Stock),
		})
	}
	return &pb.ListProductsResponse{
		Products: responseProducts,
	}, nil
}

func (s *InventoryServiceServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	log.Printf("UpdateProduct request: %+v", req)
	product := entity.Product{
		ID:          req.Product.Id,
		Name:        req.Product.Name,
		Category:    req.Product.Category,
		Description: req.Product.Description,
		Price:       req.Product.Price,
		Stock:       int(req.Product.Stock),
	}

	err := s.usecase.UpdateProduct(ctx, &product)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	responseProduct := &pb.Product{
		Id:          product.ID,
		Name:        product.Name,
		Category:    product.Category,
		Description: product.Description,
		Price:       product.Price,
		Stock:       int32(product.Stock),
	}
	return &pb.UpdateProductResponse{
		Product: responseProduct,
		Message: "Product updated successfully",
	}, nil
}

func (s *InventoryServiceServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	log.Printf("DeleteProduct request: %+v", req)
	err := s.usecase.DeleteProduct(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.DeleteProductResponse{
		Message: "Product deleted successfully",
	}, nil
}
