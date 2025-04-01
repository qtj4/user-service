package handler

import (
	"context"
	"time"

	"github.com/yourusername/user-service/internal/service"
	"github.com/yourusername/user-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserGRPCServer struct {
	proto.UnimplementedUserServiceServer
	userService service.UserService
}

func NewUserGRPCServer(userService service.UserService) *UserGRPCServer {
	return &UserGRPCServer{userService: userService}
}

func (s *UserGRPCServer) GetUserById(ctx context.Context, req *proto.GetUserByIdRequest) (*proto.GetUserByIdResponse, error) {
	user, err := s.userService.GetByID(ctx, int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return &proto.GetUserByIdResponse{
		User: &proto.User{
			Id:        int32(user.ID),
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *UserGRPCServer) GetUserOrders(ctx context.Context, req *proto.GetUserOrdersRequest) (*proto.GetUserOrdersResponse, error) {
	orders, err := s.userService.GetUserOrders(ctx, int(req.UserId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user orders: %v", err)
	}
	var protoOrders []*proto.Order
	for _, order := range orders {
		o := order.(map[string]interface{})
		protoOrders = append(protoOrders, &proto.Order{
			OrderId: int32(o["order_id"].(int)),
			BookId:  int32(o["book_id"].(int)),
			Status:  o["status"].(string),
		})
	}
	return &proto.GetUserOrdersResponse{Orders: protoOrders}, nil
}