package models

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn int16 `json:"expires_in"`
}
