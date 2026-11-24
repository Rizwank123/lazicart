package grpc_handler

import (
	"context"

	"github.com/Rizwank123/LahziCart/auth-service/internal/models"
	authpb "github.com/Rizwank123/LahziCart/shared/proto"
)

type AuthGRPCServer struct {
	authpb.UnimplementedAuthServiceServer
	AuthService models.UserAuthService
}

func NewAuthGRPCServer(as models.UserAuthService) *AuthGRPCServer {
	return &AuthGRPCServer{AuthService: as}
}

func (s *AuthGRPCServer) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	valid, userId, jti, role, err := s.AuthService.ValidateToken(req.Token)
	if err != nil {
		return nil, err
	}
	return &authpb.ValidateTokenResponse{Valid: valid, UserId: userId, Jti: jti, Role: role}, nil
}
