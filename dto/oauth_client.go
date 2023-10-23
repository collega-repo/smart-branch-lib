package dto

type Oauth2Client struct {
	ClientId         string `json:"clientId"`
	ClientSecret     string `json:"clientSecret"`
	ExpiresIn        int64  `json:"expiresIn"`
	ExpiresRefreshIn int64  `json:"expiresRefreshIn"`
	Domain           string `json:"domain"`
	ApplicationId    string `json:"applicationId"`
}
