package dto

type Oauth2Token struct {
	AccessToken      string  `json:"access_token"`
	RefreshToken     string  `json:"refresh_token"`
	ExpiresIn        int64   `json:"expires_in"`
	ExpiresRefreshIn int64   `json:"expires_refresh_in"`
	TokenType        string  `json:"token_type"`
	UserMap          UserMap `json:"-"`
}
