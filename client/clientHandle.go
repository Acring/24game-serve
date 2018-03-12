package client

import (
	"dot24_server/frame"
	"fmt"
)

func HandleRawMsg(frame frame.Frame, client *Client){
	switch frame.Type {
	case "match":
		handleMatch(frame.Data.(map[string]interface{})["username"].(string), client)
		break
	case "stopmatch":
		handleStopMatch(client)
		break
	case "matchresult":
		handleMatchResult(client)
		break
	case "matchOver":
		handleMatchOver(client)
		break
	}
}

func handleMatchOver(client *Client){
	GetInstance().Client2Room[client].BreakRoom <- client
}
func handleMatchResult(client *Client){  // 有客户端给出答案
	fmt.Println(client.Username,"给出答案")
	GetInstance().Client2Room[client].Result <- client
}
func handleStopMatch(client *Client){  // 停止匹配
	fmt.Println("停止匹配: ", client.Username)
	if client.Status == "matching"{
		client.Status = "waiting"

		fb := &frame.Frame{Type: "stopmatch", Code: 1000}

		client.Send <- fb
	}
}

func handleMatch(name string, client *Client){  // 分配姓名开始匹配
	client.Username = name

	fmt.Println(client.Id, client.Username, "开始匹配")

	if client.Status == "waiting"{  // 用户从等待状态开始匹配

		client.Status = "matching"
		clientsNum := len(GetInstance().Clients)
		data := frame.MatchingData {OnLineNum: clientsNum}
		feedback := &frame.Frame{Type:"matching", Code: 1000, Data: data}
		client.Send <- feedback  // 发送反馈给客户端
		// 将客户端放入匹配列表
		if client.Username == "芷芊" || client.Username == "刘圳"{
			GetInstance().SpecialClients = append(GetInstance().SpecialClients, client)
		}else{
			GetInstance().MatchingClients = append(GetInstance().MatchingClients, client)
		}

	}
}

