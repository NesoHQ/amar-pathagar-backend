package authhandler

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	User         *UserDTO `json:"user"`
}

type UserDTO struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	Role         string `json:"role"`
	SuccessScore int    `json:"success_score"`
}
