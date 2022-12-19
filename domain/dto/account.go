package dto

type AccountInfoResp struct {
	Tel         string `json:"tel"`
	WechatBound bool   `json:"wechatBound"`
}

type AccountBindTelCodeReq struct {
	Tel        string `json:"tel" binding:"required"`
	NationCode string `json:"nationCode"`
}

type AccountBindTelReq struct {
	Tel        string `json:"tel" binding:"required"`
	Code       string `json:"code" binding:"required"`
	NationCode string `json:"nationCode"`
}

type AccountBindWechatReq struct {
	Code string `json:"string"`
}

type AccountUnBindWechatReq struct {
}

type AccountSwitchTelCodeReq struct {
}

type AccountCheckSwitchCodeReq struct {
}

type AccountSwitchTelReq struct {
}

type AccountSwitchTargetCodeReq struct {
}
