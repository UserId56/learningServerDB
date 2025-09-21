package models

type AuthResponse struct {
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
	UserID         int64  `json:"user_id"`
	UserProfilesID int64  `json:"user_profiles_id"`
}

type RefreshTokenData struct {
	RefreshToken string `json:"refresh_token"`
}
