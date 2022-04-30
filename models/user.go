package models

// swagger:parameters User signin
type User struct {
	Password string `json:"password"`
	Username string `json:"username"`
}
