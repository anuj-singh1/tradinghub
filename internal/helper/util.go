package helper

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"tradingdata/internal/log"
)

func NewSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func PostApiExecutor(baseURL string, params map[string]string) ([]byte, error) {
	log.GetLogger().
		WithFields(logrus.Fields{"baseurl": baseURL}).
		Infoln("PostApiExecutor---")
	client := &http.Client{}
	jsonBody, err := json.Marshal(params)
	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.GetLogger().Errorf(err.Error())
		return []byte(`{}`), errors.New("pre request failure")
	}
	req.Header.Add("Content-Type", "application/json")
	log.GetLogger().Debugln("Request [POST] url::%+v", req)
	resp, err := client.Do(req)
	if err != nil {
		log.GetLogger().Errorf(err.Error())
		return []byte(`{}`), errors.New("post request failure")
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.GetLogger().Errorf(err.Error())
		}
	}(resp.Body)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.GetLogger().Errorf(err.Error())
		return []byte(`{}`), errors.New("response body unreadable")
	}
	if resp.StatusCode != 200 {
		return []byte(`{}`), errors.New(string(bodyBytes))
	}
	return bodyBytes, nil
}

func GetApiExecutor(baseURL string, params map[string]string, headers map[string]string) ([]byte, error) {
	logger := log.GetLogger()
	logger.Debugln(baseURL, params)
	client := &http.Client{}
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return []byte(`{}`), fmt.Errorf("failed to create request: %+v", err)
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	query := req.URL.Query()
	for key, value := range params {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()
	logger.Debugln("Request [GET] url::", req.URL.String())
	resp, err := client.Do(req)
	if err != nil {
		return []byte(`{}`), fmt.Errorf("failed to create client: %+v", err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logger.Errorf("error closing response body: %+v", err)
		}
	}(resp.Body)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(`{}`), fmt.Errorf("failed to read response body: %+v", err)
	}
	if resp.StatusCode != 200 {
		return []byte(`{}`), fmt.Errorf("failed response code: %+v\nbody %+v", resp.StatusCode, string(bodyBytes))
	}
	return bodyBytes, nil
}
