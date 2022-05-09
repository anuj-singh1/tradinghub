package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"tradingdata/internal/config"
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
	Status			string			`json:"status"`
	Code			int				`json:"code"`
	Message			interface{}		`json:"message,omitempty"`
	AccessToken		string			`json:"access_token"`
}

func getAuthCodeUrl(c *gin.Context){
	globalInstance, _ := c.MustGet(config.GIN_ENV_GLOBAL_INSTANCE).(config.GlobalInstance)
	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&state=success",
		config.FyersBaseUrl+config.GenerateAuthCodePath, globalInstance.Config.ClientId, globalInstance.Config.RedirectUrl)
	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

func login(c *gin.Context){
	globalinstance, _ := c.MustGet(config.GIN_ENV_GLOBAL_INSTANCE).(config.GlobalInstance)
	value, isexist := c.GetQuery("s")
	authCode := ""
	var response = AuthCodeResponse{}
	if isexist && value == "ok" {
		authCode, isexist = c.GetQuery("auth_code")
		if isexist {
			appID := fmt.Sprintf("%s:%s", globalinstance.Config.ClientId, globalinstance.Config.SecretKey)
			log.GetLogger().Infoln(appID)
			appIDHash := fmt.Sprintf("%x", helper.NewSHA256([]byte(appID))[:])
			var params = make(map[string]string)
			params["grant_type"] = "authorization_code"
			params["appIdHash"] = appIDHash
			params["code"] = authCode
			bodyBytes, err := helper.PostApiExecutor(config.FyersBaseUrl + config.ValidateCodePath, params)
			err = json.Unmarshal(bodyBytes, &response)
			log.GetLogger().Infoln(response, err)
		}
	}
	c.JSON(http.StatusOK, response)
}