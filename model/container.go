package model

import (
	"github.com/mumushuiding/util"
)

// Container 参数和结果容器
type Container struct {
	Header CHeader `json:"header,omitemtpy"`
	Body   CBody   `json:"body,omitempty"`
}

// CHeader CHeader
type CHeader struct {
	Token string `json:"token,omitempty"`
	Msg   string `json:"msg,omitempty"`
}

// CBody 用于获取前台参数和返回结果
type CBody struct {
	Data       []interface{} `json:"data,omitempty"`
	Total      int           `json:"total,omitempty"`
	StartIndex int           `json:"start_index,omitempty"`
	MaxResults int           `json:"max_results,omitempty"`
	StartDate  string        `json:"start_date,omitempty"`
	EndDate    string        `json:"end_date,omitempty"`
	UserName   string        `json:"username,omitempty"`
	UserID     int           `json:"user_id,omitempty"`
	Method     string        `json:"method,omitempty"`
	Metrics    string        `json:"metrics,omitempty"`
	Fields     []string      `json:"fields,omitempty"`
	Order      string        `json:"order,omitempty"`
}

// ToString ToString
func (c *Container) ToString() string {
	str, _ := util.ToJSONStr(c)
	return str
}
