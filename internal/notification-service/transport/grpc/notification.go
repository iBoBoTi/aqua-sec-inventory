package grpc

import (
	"context"
	"time"

	"github.com/iBoBoTi/aqua-sec-inventory/internal/notification-service/domain"
	"github.com/iBoBoTi/aqua-sec-inventory/internal/notification-service/usecase"
	pb "github.com/iBoBoTi/aqua-sec-inventory/proto/notification"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotificationGRPCService struct {
	pb.UnimplementedNotificationServiceServer
	notificationUC usecase.NotificationUsecase
}

func NewNotificationGRPCService(notificationUC usecase.NotificationUsecase) *NotificationGRPCService {
	return &NotificationGRPCService{
		notificationUC: notificationUC,
	}
}

// GetAllNotifications retrieves all notifications for a given user_id.
func (s *NotificationGRPCService) GetAllNotifications(
	ctx context.Context,
	req *pb.GetAllNotificationsRequest,
) (*pb.GetAllNotificationsResponse, error) {

	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	notifs, err := s.notificationUC.GetAllNotifications(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get notifications: %v", err)
	}

	var pbNotifs []*pb.Notification
	for _, n := range notifs {
		pbNotifs = append(pbNotifs, convertToPBNotification(n))
	}

	return &pb.GetAllNotificationsResponse{
		Notifications: pbNotifs,
	}, nil
}

// ClearSingleNotification deletes one notification by notification_id.
func (s *NotificationGRPCService) ClearSingleNotification(
	ctx context.Context,
	req *pb.ClearSingleNotificationRequest,
) (*pb.ClearSingleNotificationResponse, error) {

	if req.NotificationId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid notification_id")
	}

	err := s.notificationUC.ClearNotification(req.NotificationId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to clear notification: %v", err)
	}

	return &pb.ClearSingleNotificationResponse{Message: "Notification cleared"}, nil
}

// ClearAllNotifications deletes all notifications for a specific user_id.
func (s *NotificationGRPCService) ClearAllNotifications(
	ctx context.Context,
	req *pb.ClearAllNotificationsRequest,
) (*pb.ClearAllNotificationsResponse, error) {

	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	err := s.notificationUC.ClearAllNotifications(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to clear all notifications: %v", err)
	}

	return &pb.ClearAllNotificationsResponse{Message: "All notifications cleared"}, nil
}

// map domain.Notification to proto Notification.
func convertToPBNotification(n domain.Notification) *pb.Notification {
	return &pb.Notification{
		Id:        n.ID,
		UserId:    n.UserID,
		Message:   n.Message,
		CreatedAt: n.CreatedAt.Format(time.RFC3339),
	}
}
