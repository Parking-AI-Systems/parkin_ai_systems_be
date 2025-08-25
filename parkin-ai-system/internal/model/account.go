package model

type SignInInput struct {
	Account  string `v:"bail|required|length:6,18" json:"account"`
	Password string `v:"bail|required|length:8,256" json:"password"`
}
type SignInOutput struct {
	AccessTokenOutput
	RefreshToken string `json:"rf_token"`
}
