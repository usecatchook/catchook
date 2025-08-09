package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/theotruvelot/catchook/internal/domain/user"
	"github.com/theotruvelot/catchook/pkg/logger"
	"github.com/theotruvelot/catchook/pkg/response"
	"github.com/theotruvelot/catchook/storage/postgres/generated"
)

type userRepository struct {
	db      *pgxpool.Pool
	queries *generated.Queries
	logger  logger.Logger
}

func NewUserRepository(db *pgxpool.Pool, logger logger.Logger) user.Repository {
	return &userRepository{
		db:      db,
		queries: generated.New(db),
		logger:  logger,
	}
}

func (r *userRepository) Create(ctx context.Context, user *user.User) error {
	r.logger.Debug(ctx, "Creating user in database",
		logger.String("email", user.Email),
	)

	result, err := r.queries.CreateUser(ctx,
		user.Email,
		user.Role,
		user.Password,
		user.FirstName,
		user.LastName,
		user.IsActive,
	)

	if err != nil {
		r.logger.Error(ctx, "Failed to create user in database",
			logger.String("email", user.Email),
			logger.Error(err),
		)
		return fmt.Errorf("failed to create user: %w", err)
	}

	user.ID = result.ID.String()
	user.CreatedAt = result.CreatedAt.Time
	user.UpdatedAt = result.UpdatedAt.Time

	r.logger.Debug(ctx, "User created successfully in database",
		logger.String("user_id", user.ID),
		logger.String("email", user.Email),
	)

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	r.logger.Debug(ctx, "Getting user by ID from database",
		logger.String("user_id", id),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		r.logger.Error(ctx, "Invalid UUID format",
			logger.String("user_id", id),
			logger.Error(err),
		)
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	result, err := r.queries.GetUserByID(ctx, uid)
	if err != nil {
		if err.Error() == "no rows in result set" {
			r.logger.Debug(ctx, "User not found in database",
				logger.String("user_id", id),
			)
			return nil, fmt.Errorf("user not found")
		}
		r.logger.Error(ctx, "Failed to get user from database",
			logger.String("user_id", id),
			logger.Error(err),
		)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	r.logger.Debug(ctx, "User found in database",
		logger.String("user_id", id),
		logger.String("email", result.Email),
	)

	return &user.User{
		ID:        result.ID.String(),
		Email:     result.Email,
		Role:      result.Role,
		Password:  result.PasswordHash,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		IsActive:  result.IsActive,
		CreatedAt: result.CreatedAt.Time,
		UpdatedAt: result.UpdatedAt.Time,
	}, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	r.logger.Debug(ctx, "Getting user by email from database",
		logger.String("email", email),
	)

	result, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err.Error() == "no rows in result set" {
			r.logger.Debug(ctx, "User not found in database",
				logger.String("email", email),
			)
			return nil, fmt.Errorf("user not found")
		}
		r.logger.Error(ctx, "Failed to get user from database",
			logger.String("email", email),
			logger.Error(err),
		)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	r.logger.Debug(ctx, "User found in database",
		logger.String("user_id", result.ID.String()),
		logger.String("email", email),
	)

	return &user.User{
		ID:        result.ID.String(),
		Email:     result.Email,
		Role:      result.Role,
		Password:  result.PasswordHash,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		IsActive:  result.IsActive,
		CreatedAt: result.CreatedAt.Time,
		UpdatedAt: result.UpdatedAt.Time,
	}, nil
}

func (r *userRepository) Update(ctx context.Context, user *user.User) error {
	r.logger.Debug(ctx, "Updating user in database",
		logger.String("user_id", user.ID),
		logger.String("email", user.Email),
	)

	uid, err := uuid.Parse(user.ID)
	if err != nil {
		r.logger.Error(ctx, "Invalid UUID format for update",
			logger.String("user_id", user.ID),
			logger.Error(err),
		)
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	result, err := r.queries.UpdateUser(ctx,
		uid,
		user.Role,
		user.FirstName,
		user.LastName,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			r.logger.Debug(ctx, "User not found for update",
				logger.String("user_id", user.ID),
			)
			return fmt.Errorf("user not found")
		}
		r.logger.Error(ctx, "Failed to update user in database",
			logger.String("user_id", user.ID),
			logger.Error(err),
		)
		return fmt.Errorf("failed to update user: %w", err)
	}

	user.UpdatedAt = result.UpdatedAt.Time

	r.logger.Debug(ctx, "User updated successfully in database",
		logger.String("user_id", user.ID),
		logger.String("email", user.Email),
	)

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	r.logger.Debug(ctx, "Deleting user from database",
		logger.String("user_id", id),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		r.logger.Error(ctx, "Invalid UUID format for deletion",
			logger.String("user_id", id),
			logger.Error(err),
		)
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	err = r.queries.DeleteUser(ctx, uid)
	if err != nil {
		if err.Error() == "no rows in result set" {
			r.logger.Debug(ctx, "User not found for deletion",
				logger.String("user_id", id),
			)
			return fmt.Errorf("user not found")
		}
		r.logger.Error(ctx, "Failed to delete user from database",
			logger.String("user_id", id),
			logger.Error(err),
		)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	r.logger.Debug(ctx, "User deleted successfully from database",
		logger.String("user_id", id),
	)

	return nil
}

func (r *userRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	r.logger.Debug(ctx, "Checking if email exists in database",
		logger.String("email", email),
	)

	result, err := r.queries.CheckEmailExists(ctx, email)
	if err != nil {
		r.logger.Error(ctx, "Failed to check email existence in database",
			logger.String("email", email),
			logger.Error(err),
		)
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	r.logger.Debug(ctx, "Email existence check completed",
		logger.String("email", email),
		zap.Bool("exists", result),
	)

	return result, nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID string, hashedPassword string) error {
	r.logger.Debug(ctx, "Updating user password in database",
		logger.String("user_id", userID),
	)

	uid, err := uuid.Parse(userID)
	if err != nil {
		r.logger.Error(ctx, "Invalid UUID format for password update",
			logger.String("user_id", userID),
			logger.Error(err),
		)
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	_, err = r.queries.UpdateUserPassword(ctx, uid, hashedPassword)
	if err != nil {
		if err.Error() == "no rows in result set" {
			r.logger.Debug(ctx, "User not found for password update",
				logger.String("user_id", userID),
			)
			return fmt.Errorf("user not found")
		}
		r.logger.Error(ctx, "Failed to update user password in database",
			logger.String("user_id", userID),
			logger.Error(err),
		)
		return fmt.Errorf("failed to update user password: %w", err)
	}

	r.logger.Debug(ctx, "User password updated successfully in database",
		logger.String("user_id", userID),
	)

	return nil
}

func (r *userRepository) CountUsers(ctx context.Context) (int64, error) {
	r.logger.Debug(ctx, "Counting users in database")

	count, err := r.queries.CountUsers(ctx)
	if err != nil {
		r.logger.Error(ctx, "Failed to count users in database",
			logger.Error(err),
		)
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	r.logger.Debug(ctx, "Users counted successfully",
		logger.Int("count", int(count)),
	)

	return count, nil
}

func (r *userRepository) List(ctx context.Context, page, limit int) ([]*user.User, *response.Pagination, error) {
	r.logger.Debug(ctx, "Listing users from database",
		logger.Int("page", page),
		logger.Int("limit", limit))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	total, err := r.CountUsers(ctx)
	if err != nil {
		r.logger.Error(ctx, "Failed to count users for pagination", logger.Error(err))
		return nil, nil, fmt.Errorf("failed to count users: %w", err)
	}

	results, err := r.queries.ListUsers(ctx, int32(limit), int32(offset))
	if err != nil {
		r.logger.Error(ctx, "Failed to list users from database", logger.Error(err))
		return nil, nil, fmt.Errorf("failed to list users: %w", err)
	}

	users := make([]*user.User, len(results))
	for i, result := range results {
		users[i] = &user.User{
			ID:        result.ID.String(),
			Email:     result.Email,
			Role:      result.Role,
			Password:  result.PasswordHash,
			FirstName: result.FirstName,
			LastName:  result.LastName,
			IsActive:  result.IsActive,
			CreatedAt: result.CreatedAt.Time,
			UpdatedAt: result.UpdatedAt.Time,
		}
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages < 1 {
		totalPages = 1
	}

	hasNext := page < totalPages
	hasPrev := page > 1

	pagination := &response.Pagination{
		CurrentPage: page,
		TotalPages:  totalPages,
		Total:       int(total),
		Limit:       limit,
		HasNext:     hasNext,
		HasPrev:     hasPrev,
	}

	r.logger.Debug(ctx, "Users listed successfully from database",
		logger.Int("count", len(users)),
		logger.Int("page", page),
		logger.Int("limit", limit),
		logger.Int("total", int(total)),
		logger.Int("total_pages", totalPages))

	return users, pagination, nil
}
