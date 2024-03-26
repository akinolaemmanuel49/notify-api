package models

type AuthCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthPasswordChange struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
