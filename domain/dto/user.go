package dto

import "jiqi/internal/user/domain"

type UserSetCityReq struct {
	Region   string `json:"region"`
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
}

type UserSetHeadResp struct {
	Headimgurl string `json:"headimgurl"`
}

type UserSetNicknameReq struct {
	Nickname string `json:"nickname" binding:"required"`
}

type UserSetSexReq struct {
	Sex int `json:"sex" binding:"required"`
}

type UserSetSignatureReq struct {
	Signature string `json:"signature" binding:"required,max=140"`
}

type UserSetHeightReq struct {
	Height int `json:"height" binding:"required"`
}

type UserSetWeightReq struct {
	Weight int `json:"weight" binding:"required"`
}

type UserSetBirthdayReq struct {
	Birthday string `json:"birthday" binding:"required"`
}

type UserProfileResp struct {
	Uri        string `json:"uri"`
	Nickname   string `json:"nickname"`
	Signature  string `json:"signature"`
	Sex        int    `json:"sex"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Country    string `json:"country"`
	Region     string `json:"region"`
	Headimgurl string `json:"headimgurl"`
	Height     int    `json:"height"`
	Weight     int    `json:"weight"`
	Birthday   string `json:"birthday"`
}

func TransProfileDto(profile *domain.Profile) *UserProfileResp {
	return &UserProfileResp{
		Uri:        profile.Uri,
		Nickname:   profile.Nickname,
		Signature:  profile.Signature,
		Sex:        profile.Sex,
		Province:   profile.Province,
		City:       profile.City,
		Country:    profile.Country,
		Headimgurl: profile.Headimgurl,
		Height:     profile.Height,
		Weight:     profile.Weight,
		Birthday:   profile.Birthday,
		Region:     profile.Region,
	}
}
