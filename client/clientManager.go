package client

import(
	//"encoding/json"
	"fmt"
	time2 "time"
)

var cm *ClientManager


func GetInstance() *ClientManager{  // 获取单例
	if cm == nil{
		cm = &ClientManager{}
	}
	return cm
}


type Message struct {  // 信息结构体
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   []byte `json:"content,omitempty"`
}


type ClientManager struct {
	Rooms 			[]Room		// 对战房间
	Client2Room		map[*Client]*Room
	Clients    		map[*Client]int	// 所有玩家
	MatchingClients []*Client 	// 匹配中的玩家
	SpecialClients []*Client	// 特别匹配玩家
	Broadcast  		chan []byte      // 广播通道
	Register   		chan *Client     // 注册新客户端表通道
	UnRegister 		chan *Client     // 注销通道
	Test 			chan bool
}


/**
管理类初始化, 监听各种客户端行为
 */
func (cm *ClientManager) Start() {

	for {
		select {
			case conn := <- cm.Register:
				cm.Clients[conn] = 0
				fmt.Println("客户端登录:", conn.Id)
			case conn := <- cm.UnRegister:
				fmt.Println("客户端下线", conn.Id)
				cm.movePlayer(conn)
			case test := <- cm.Test:
				fmt.Println(test)
		}
	}
}

func (cm *ClientManager)Matching()  {  // 对匹配中的用户进行配对
	time := time2.NewTimer(5 * time2.Second)
	for{
		select{
			case <- time.C:
				time.Reset(5 * time2.Second)
				if len(cm.MatchingClients) >= 2{
					fmt.Println("MatchingClients")
					cm.buildRoom(cm.MatchingClients[0], cm.MatchingClients[1])
					cm.MatchingClients = cm.MatchingClients[2:]  // 清除前两个匹配用户
				}
				if len(cm.SpecialClients) >= 2{
					fmt.Println("SpecialClients")
					cm.buildRoom(cm.SpecialClients[0], cm.SpecialClients[1])
					cm.SpecialClients = []*Client{}  // 清空特殊匹配列表
				}
				break
		}
	}
}

func (cm *ClientManager) Init(){
	cm.SpecialClients = make([]*Client, 0)
	cm.MatchingClients = make([]*Client, 0)
	cm.Clients = make(map[*Client]int)
	cm.Register = make(chan *Client)
	cm.UnRegister = make(chan *Client)
	cm.Broadcast = make(chan []byte)
	cm.Test = make(chan bool)
	cm.Rooms = make([]Room, 0)
	cm.Client2Room = make(map[*Client]*Room)
}

func(cm *ClientManager) buildRoom(p1 *Client, p2 *Client){  // 建立对战房间
	fmt.Println("buildRoom")
	room := Room{Clients: [2]*Client{p1,p2}}
	cm.Client2Room[p1] = &room
	cm.Client2Room[p2] = &room

	cm.Rooms = append(cm.Rooms, room)
	go room.Start()
}

func (cm *ClientManager) movePlayer(conn *Client){  // 设置客户端下线
	cm.Clients[conn] = -1  // 设置客户端下线
	for index,client  := range cm.MatchingClients{
		if client == conn{
			cm.MatchingClients = append(cm.MatchingClients[:index],cm.MatchingClients[index+1:]...)
		}
	}
	for index,client  := range cm.SpecialClients{
		if client == conn{
			cm.SpecialClients = append(cm.SpecialClients[:index],cm.SpecialClients[index+1:]...)
		}
	}
	if cm.Client2Room[conn] != nil{
		cm.Client2Room[conn].breakRoom()
	}
}