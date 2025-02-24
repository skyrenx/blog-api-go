package entities

//db table is "users"
type User struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"` // TODO ensure password is not retrievable.
	// Malicious users with access to encrypted passwords can attemp to decrypt the password offline.
	Enabled  bool   `json:"enabled" db:"enabled"`
}
