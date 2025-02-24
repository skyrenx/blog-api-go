package dto

type UserWithoutPassword struct {
	Username string `json:"username" db:"username"`
	// Malicious users with access to encrypted passwords can attempt to decrypt the password offline.
	Enabled  bool   `json:"enabled" db:"enabled"`
}