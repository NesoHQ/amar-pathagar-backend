package repository

import (
	"context"
	"database/sql"

	"github.com/yourusername/online-library/internal/auth"
	"github.com/yourusername/online-library/internal/domain"
	"github.com/yourusername/online-library/internal/user"
	"go.uber.org/zap"
)

type UserRepository struct {
	db  *sql.DB
	log *zap.Logger
}

var _ auth.UserRepo = (*UserRepository)(nil)
var _ user.UserRepo = (*UserRepository)(nil)

func NewUserRepository(db *sql.DB, log *zap.Logger) *UserRepository {
	return &UserRepository{db: db, log: log}
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, full_name, role, success_score, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		u.ID, u.Username, u.Email, u.PasswordHash, u.FullName, u.Role, u.SuccessScore, u.CreatedAt, u.UpdatedAt)
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	u := &domain.User{}
	query := `
		SELECT id, username, email, password_hash, full_name, role, 
		       COALESCE(avatar_url, ''), COALESCE(bio, ''),
		       location_lat, location_lng, COALESCE(location_address, ''),
		       COALESCE(success_score, 100), COALESCE(books_shared, 0),
		       COALESCE(books_received, 0), COALESCE(reviews_received, 0),
		       COALESCE(ideas_posted, 0), COALESCE(total_upvotes, 0),
		       COALESCE(total_downvotes, 0), COALESCE(is_donor, false),
		       created_at, updated_at
		FROM users WHERE id = $1
	`
	var locationLat, locationLng sql.NullFloat64
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName, &u.Role,
		&u.AvatarURL, &u.Bio, &locationLat, &locationLng, &u.LocationAddress,
		&u.SuccessScore, &u.BooksShared, &u.BooksReceived, &u.ReviewsReceived,
		&u.IdeasPosted, &u.TotalUpvotes, &u.TotalDownvotes, &u.IsDonor,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	u.LocationLat = float64Ptr(locationLat)
	u.LocationLng = float64Ptr(locationLng)
	return u, err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	u := &domain.User{}
	query := `
		SELECT id, username, email, password_hash, full_name, role, 
		       COALESCE(avatar_url, ''), COALESCE(bio, ''),
		       location_lat, location_lng, COALESCE(location_address, ''),
		       COALESCE(success_score, 100), COALESCE(books_shared, 0),
		       COALESCE(books_received, 0), COALESCE(reviews_received, 0),
		       COALESCE(ideas_posted, 0), COALESCE(total_upvotes, 0),
		       COALESCE(total_downvotes, 0), COALESCE(is_donor, false),
		       created_at, updated_at
		FROM users WHERE email = $1
	`
	var locationLat, locationLng sql.NullFloat64
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName, &u.Role,
		&u.AvatarURL, &u.Bio, &locationLat, &locationLng, &u.LocationAddress,
		&u.SuccessScore, &u.BooksShared, &u.BooksReceived, &u.ReviewsReceived,
		&u.IdeasPosted, &u.TotalUpvotes, &u.TotalDownvotes, &u.IsDonor,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	u.LocationLat = float64Ptr(locationLat)
	u.LocationLng = float64Ptr(locationLng)
	return u, err
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	u := &domain.User{}
	query := `
		SELECT id, username, email, password_hash, full_name, role, 
		       COALESCE(avatar_url, ''), COALESCE(bio, ''),
		       location_lat, location_lng, COALESCE(location_address, ''),
		       COALESCE(success_score, 100), COALESCE(books_shared, 0),
		       COALESCE(books_received, 0), COALESCE(reviews_received, 0),
		       COALESCE(ideas_posted, 0), COALESCE(total_upvotes, 0),
		       COALESCE(total_downvotes, 0), COALESCE(is_donor, false),
		       created_at, updated_at
		FROM users WHERE username = $1
	`
	var locationLat, locationLng sql.NullFloat64
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName, &u.Role,
		&u.AvatarURL, &u.Bio, &locationLat, &locationLng, &u.LocationAddress,
		&u.SuccessScore, &u.BooksShared, &u.BooksReceived, &u.ReviewsReceived,
		&u.IdeasPosted, &u.TotalUpvotes, &u.TotalDownvotes, &u.IsDonor,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	u.LocationLat = float64Ptr(locationLat)
	u.LocationLng = float64Ptr(locationLng)
	return u, err
}

func (r *UserRepository) Update(ctx context.Context, id string, u *domain.User) error {
	query := `
		UPDATE users SET full_name = $1, bio = $2, avatar_url = $3, 
		       location_lat = $4, location_lng = $5, location_address = $6, updated_at = $7
		WHERE id = $8
	`
	_, err := r.db.ExecContext(ctx, query,
		u.FullName, u.Bio, u.AvatarURL,
		nullFloat64(u.LocationLat), nullFloat64(u.LocationLng), u.LocationAddress,
		u.UpdatedAt, id)
	return err
}

func (r *UserRepository) AddInterests(ctx context.Context, userID string, interests []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, interest := range interests {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO user_interests (user_id, interest, weight)
			VALUES ($1, $2, 1.0)
			ON CONFLICT (user_id, interest) DO UPDATE SET weight = user_interests.weight + 0.1
		`, userID, interest)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *UserRepository) GetTopUsers(ctx context.Context, limit int) ([]*domain.User, error) {
	query := `
		SELECT id, username, email, full_name, role, 
		       COALESCE(avatar_url, ''), COALESCE(success_score, 100),
		       COALESCE(books_shared, 0), COALESCE(books_received, 0),
		       created_at, updated_at
		FROM users
		ORDER BY success_score DESC
		LIMIT $1
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		u := &domain.User{}
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.Role,
			&u.AvatarURL, &u.SuccessScore, &u.BooksShared, &u.BooksReceived,
			&u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
