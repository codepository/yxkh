package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/codepository/yxkh/config"
)

// RemoteClient 调用远程api客户端
type RemoteClient struct {
	Client     *http.Client
	UserAPIURL string
}

// APIClient 调用远程api客户端
var APIClient *RemoteClient

// StartRemoteClient 启动调用远程api客户端
func StartRemoteClient() {
	APIClient = &RemoteClient{}
	APIClient.Client = &http.Client{}
	APIClient.UserAPIURL = config.Config.UserAPIURL
}

// StopRemoteClient 关闭调用远程api客户端
func StopRemoteClient() {
	APIClient.Client.CloseIdleConnections()
}

// Post Post
func (c *RemoteClient) Post(url string, params interface{}) ([]byte, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	contentType := "application/json;charset=utf-8"
	body := bytes.NewBuffer(b)
	resp, err := c.Client.Post(url, contentType, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	statusCode := resp.StatusCode
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// log.Println("statusCode:", statusCode)
	// log.Println("content:", string(content))
	if statusCode != 200 {
		return nil, errors.New(string(content))
	}
	return content, nil
}

// Get Get
func (c *RemoteClient) Get() {

}
