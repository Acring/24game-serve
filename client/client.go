package client

import (
	"github.com/gorilla/websocket"
	"encoding/json"
	"fmt"
	"dot24_server/frame"
)



type Client struct {
	Id     string
	Username string
	Socket *websocket.Conn  // 连接
	Send   chan *frame.Frame  // 待发送信息
	Status string // waiting-主页面 matching-匹配中 matched-对战中
}


/**
接收客户端信息
 */
func (c *Client) Read() {

	defer func() {
		c.Socket.Close()
	}()
	newFrame := frame.Frame{}
	for {

		_, data, err := c.Socket.ReadMessage()
		if err != nil{
			fmt.Println("客户端断线")
			GetInstance().UnRegister <- c
			break
		}
		err = json.Unmarshal(data, &newFrame)
		if err != nil{
			fmt.Println("解析客户端数据失败", err.Error())
			GetInstance().UnRegister <- c
			fmt.Println(string(data))
		}

		fmt.Println(newFrame)
		HandleRawMsg(newFrame, c)
	}
}
/**
向客户端发送信息
 */
func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			feedback,err := json.Marshal(message)
			if err != nil{
				fmt.Print("解析反馈失败")
			}
			err = c.Socket.WriteMessage(websocket.TextMessage, feedback)
			if err != nil{
				fmt.Println("信息发送失败:", err.Error())
			}
		}
	}
}
