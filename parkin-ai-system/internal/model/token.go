package model

type AccessTokenOutput struct {
	Uid         uint64 `json:"uid"`
	AccessToken string `json:"ac_token"`
	ExpTime     int64  `json:"exp_time"`
}
