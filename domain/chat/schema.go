// schema - 2024/12/16
// Author: wangzx
// Description:

package chat

import "encoding/xml"

// 定义结构体
type Message struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
	MsgId        int64    `xml:"MsgId"`
	MsgDataId    string   `xml:"MsgDataId"`
	Idx          string   `xml:"Idx"`
}
