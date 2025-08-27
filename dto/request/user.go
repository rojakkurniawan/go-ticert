package request

type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email,max=100"`
	Password  string `json:"password" validate:"required,min=8,max=50"`
	FirstName string `json:"first_name" validate:"required,max=50"`
	LastName  string `json:"last_name" validate:"required,max=50"`
}

type UpdateUserRequest struct {
	FirstName *string `json:"first_name" validate:"omitempty,max=50"`
	LastName  *string `json:"last_name" validate:"omitempty,max=50"`
	Role      *string `json:"role" validate:"omitempty,oneof=user admin"`
}

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email,max=100"`
	Password  string `json:"password" validate:"required,min=8,max=50"`
	FirstName string `json:"first_name" validate:"required,max=50"`
	LastName  string `json:"last_name" validate:"required,max=50"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=8,max=50"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=50"`
}

type UpdateEmailRequest struct {
	Email string `json:"email" validate:"required,email,max=100"`
}
