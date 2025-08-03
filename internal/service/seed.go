package service

import (
	"context"

	"github.com/theotruvelot/catchook/internal/domain/user"
	"github.com/theotruvelot/catchook/pkg/logger"
)

type SetupService interface {
	IsFirstTimeSetup(ctx context.Context) (bool, error)
}

type setupService struct {
	userRepo user.Repository
	logger   logger.Logger
}

func NewSetupService(userRepo user.Repository, logger logger.Logger) SetupService {
	return &setupService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// IsFirstTimeSetup checks if this is the first time the application is being used
func (s *setupService) IsFirstTimeSetup(ctx context.Context) (bool, error) {
	count, err := s.userRepo.CountUsers(ctx)
	if err != nil {
		s.logger.Error(ctx, "Failed to count users during setup check", logger.Error(err))
		return false, err
	}

	isFirstTime := count == 0
	s.logger.Debug(ctx, "Setup check completed",
		logger.Any("is_first_time", isFirstTime),
		logger.Int("user_count", int(count)),
	)

	return isFirstTime, nil
}
