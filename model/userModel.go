package model

//	LoginModel represents a json for registering
type LoginModel struct {
	Sid      string `json:"sid"      binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserModel represents a registered user.
type UserModel struct {
	Id        uint32 `json:"id"`
	Sid       string `json:"sid"`
	Username  string `json:"username"`
	Avatar    string `json:"avatar"`
	IsBlocked uint8  `json:"is_blocked"`
}

// UserInfo represents a user's info
type UserInfo struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

// AuthResponse represents a JSON web token.
type AuthResponse struct {
	Token string `json:"token"`
	IsNew uint8  `json:"is_new"`
}
