package service

import (
	"context"
	"encoding/json"
	"fmt"
	"parkin-ai-system/internal/config"
	"parkin-ai-system/internal/consts"
	"parkin-ai-system/utility"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/golang-jwt/jwt/v5"
)

type Token interface {
	Gen() (token string, err error)
	Verify(token string) (err error)
}

type AccessToken struct {
	Iss string `json:"iss"`
	Sub string `json:"sub"`
	Exp int64  `json:"exp"`
}

func (ac *AccessToken) Gen() (token string, err error) {
	cf := config.GetConfig()
	expTime := gtime.Now().Add(time.Duration(cf.Auth.AccessTokenExpireMinute) * time.Minute)
	ac.Exp = expTime.Unix()

	dataByte, err := json.Marshal(ac)
	if err != nil {
		return "", err
	}
	var mapClaims jwt.MapClaims
	if err := json.Unmarshal(dataByte, &mapClaims); err != nil {
		return "", err
	}
	return utility.GenJWT(mapClaims)
}

func (ac *AccessToken) Verify(ctx context.Context, token string) (err error) {
	jwtMap, err := utility.ParseJWT(token)
	if err != nil {
		return err
	}
	dataByte, _ := json.Marshal(jwtMap)
	err = json.Unmarshal(dataByte, &ac)
	if err != nil {
		return err
	}
	return
}

type RefreshToken struct {
	Ctx  context.Context
	Uuid string
	Uid  string
	Exp  string
	Nbf  string
}

func (rf *RefreshToken) CheckCtx() (err error) {
	if rf.Ctx == nil {
		return gerror.NewCode(consts.CodeInvalidToken)
	}
	return
}

func (rf *RefreshToken) Gen() (token string, err error) {
	cf := config.GetConfig()
	expTime := gtime.Now().Add(time.Duration(cf.Auth.RefreshTokenExpireMinute) * time.Minute)
	rf.Exp = expTime.String()
	rf.Nbf = fmt.Sprintf("%v", expTime.Unix())
	token = fmt.Sprintf("%v", rf.Uuid)
	return
}

func (rf *RefreshToken) Verify(token string) (err error) {
	if gtime.Now().Unix() >= utility.String2Int64(rf.Nbf) {
		return gerror.NewCode(consts.CodeTokenExpired)
	}
	return
}
