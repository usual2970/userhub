package domain

type AuthLoginReq struct {
	PlatformId int               `json:"platformId" binding:"required"`
	Data       map[string]string `json:"data" binding:"required"`
}

type AuthLoginResp struct {
	AccessToken string `json:"accessToken"`
	ExpiresAt   int64  `json:"expiresAt"`
}

type AuthSmsCodeReq struct {
	NationCode string `json:"nationCode"`
	Tel        string `json:"tel" binding:"required"`
}
