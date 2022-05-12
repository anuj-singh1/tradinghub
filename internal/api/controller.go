package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
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

func getAuthCodeUrl(c *gin.Context) {
	globalInstance, _ := c.MustGet(config.GIN_ENV_GLOBAL_INSTANCE).(config.GlobalInstance)
	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&state=success",
		config.FyersBaseUrl+config.GenerateAuthCodePath, globalInstance.Config.ClientId, globalInstance.Config.RedirectUrl)
	log.GetLogger().Debugln(url)
	c.String(http.StatusOK, "url: %s", url)
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

func streamingApi(c *gin.Context) {
	globalInstance, _ := c.MustGet(config.GIN_ENV_GLOBAL_INSTANCE).(config.GlobalInstance)
	chanStream := make(chan QuotesResponse)
	stock, err := c.GetQueryArray("stock[]")
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "stock parameter is required",
		})
		return
	}
	interval, err := c.GetQuery("interval")
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "interval parameter is required",
		})
		return
	}
	duration, _err := strconv.Atoi(interval)
	if _err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "interval value is invalid",
		})
		return
	}
	if duration > 60 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "interval should be less than 60",
		})
		return
	}
	go func() {
		defer close(chanStream)
		for i := 0; i < 60/duration; i++ {
			chanStream <- getQuotes(stock, globalInstance)
			time.Sleep(time.Second * time.Duration(duration))
		}
	}()
	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-chanStream; ok {
			if msg.Status == "error" {
				c.SSEvent(time.Now().String(), msg.Message)
				return false
			}
			c.SSEvent(time.Now().String(), msg.Data)
			return true
		}
		return false
	})
}

func getData(c *gin.Context) {
	globalInstance, _ := c.MustGet(config.GIN_ENV_GLOBAL_INSTANCE).(config.GlobalInstance)
	stock, err := c.GetQueryArray("stock")
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "stock parameter is required",
		})
	}
	response := getQuotes(stock, globalInstance)
	if response.Status != "ok" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error fetching data " + response.Message,
		})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func getQuotes(stock []string, instance config.GlobalInstance) QuotesResponse {
	logger := log.GetLogger()
	response := QuotesResponse{Status: "error", Message: ""}
	params := make(map[string]string)
	headers := make(map[string]string)
	params["symbols"] = strings.Join(stock, ",")
	token, _ := instance.TokenDb.GetLastToken(context.Background())
	headers["Authorization"] = instance.Config.ClientId + ":" + token.AccessToken
	headers["Content-Type"] = "application/json"
	bodyBytes, err := helper.GetApiExecutor(config.FyersDataUrl+config.QuotesPath, params, headers)
	if err != nil {
		logger.Errorln(err)
		response.Message = err.Error()
		return response
	}
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		logger.Errorf("error in masrshaling body: %+v", err)
		response.Message = err.Error()
		return response
	}
	return response
}
