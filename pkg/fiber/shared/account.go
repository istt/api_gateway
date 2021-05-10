package shared

// UserDTO store information about user
type UserDTO struct {
	Id          string   `json:"id"`
	Login       string   `json:"login"`
	Email       string   `json:"email"`
	FirstName   string   `json:"firstName"`
	LastName    string   `json:"lastName"`
	ImageUrl    string   `json:"imageUrl"`
	LangKey     string   `json:"langKey"`
	Activated   bool     `json:"activated"`
	Authorities []string `json:"authorities"`
}

// ManagedUserDTO store information about user from admin point of view
type ManagedUserDTO struct {
	UserDTO
	Password         string `json:"password"`
	CreatedBy        string `json:"createdBy"`
	CreatedDate      string `json:"createdDate"`
	LastModifiedBy   string `json:"lastModifiedBy"`
	LastModifiedDate string `json:"lastModifiedDate"`
}

// LoginVM is a value model for Login request
type LoginVM struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe"`
}

// PasswordChangeDTO is a model for password change request
type PasswordChangeDTO struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}
