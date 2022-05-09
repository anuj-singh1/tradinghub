package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"tradingdata/internal/config"
	db "tradingdata/internal/db/sqlc"
	"tradingdata/internal/helper"
	"tradingdata/internal/log"
)

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Status OK",
	})
}

type AuthCodeResponse struct {
	Status      string      `json:"status"`
	Code        int         `json:"code"`
	Message     interface{} `json:"message,omitempty"`
	AccessToken string      `json:"access_token"`
}

func getAuthCodeUrl(c *gin.Context) {
	globalInstance, _ := c.MustGet(config.GIN_ENV_GLOBAL_INSTANCE).(config.GlobalInstance)
	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&state=success",
		config.FyersBaseUrl+config.GenerateAuthCodePath, globalInstance.Config.ClientId, globalInstance.Config.RedirectUrl)
	log.GetLogger().Debugln(url)
	c.JSON(http.StatusOK, map[string]string{
		"url": url,
	})
}

func login(c *gin.Context) {
	globalinstance, _ := c.MustGet(config.GIN_ENV_GLOBAL_INSTANCE).(config.GlobalInstance)
	value, isexist := c.GetQuery("s")
	authCode := ""
	var response = AuthCodeResponse{}
	if isexist && value == "ok" {
		authCode, isexist = c.GetQuery("auth_code")
		if isexist {
			appID := fmt.Sprintf("%s:%s", globalinstance.Config.ClientId, globalinstance.Config.SecretKey)
			log.GetLogger().Debugln(appID)
			appIDHash := fmt.Sprintf("%x", helper.NewSHA256([]byte(appID))[:])
			var params = make(map[string]string)
			params["grant_type"] = "authorization_code"
			params["appIdHash"] = appIDHash
			params["code"] = authCode
			bodyBytes, err := helper.PostApiExecutor(config.FyersBaseUrl+config.ValidateCodePath, params)
			err = json.Unmarshal(bodyBytes, &response)
			log.GetLogger().Debugln(response, err)
			if err == nil && response.Code == 200 {
				createToken := db.CreateTokenParams{
					AccessToken: response.AccessToken,
				}
				_, err = globalinstance.TokenDb.CreateToken(c, createToken)
				if err != nil {
					c.JSON(http.StatusInternalServerError, err)
					return
				}
			}
		}
	}
	c.JSON(http.StatusOK, response)
}
