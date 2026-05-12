package auth

type LoginResponse struct {
	AccessToken string `json:"access_token,omitempty" example:"eyJhbGciO..."`
}
