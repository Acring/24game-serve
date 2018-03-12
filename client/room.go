package client

import (
	"math/rand"
	"time"
	"fmt"
	"dot24_server/frame"
)
const(
	LIMITTIME time.Duration = 60
)
type Room struct{
	round 	int				// 对战轮数
	Clients [2]*Client		// 两个玩家
	Score 	map[*Client]int	// 比分
	limitTime 	*time.Timer			// 每轮的时间
	question 	[4]int		// 问题, 四个整数
	Result  chan *Client   //  某个用户得出答案
	BreakRoom   chan *Client   // 用户退出房间
}

func (r *Room) Start(){
	r.init()
	r.getQuestion()
	r.matchedFeedBack()  // 返回匹配成功的反馈

	for {
		select{
			case <- r.limitTime.C:  // 超时
				r.Timeout()
				break
			case winner := <- r.Result:  // 得出答案
				r.limitTime.Reset(LIMITTIME * time.Second)  // 重置计时器
				r.getMatchResult(winner)
				break
			case <- r.BreakRoom:
				r.breakRoom()
				break
		}
	}
}


func (r *Room)getQuestion(){  // 获取新题目给两个玩家
	r.question = questionList[rand.Intn(len(questionList)+1)]
}
func (r *Room) getMatchResult(winner *Client){  // 某用户的给出答案
	r.round += 1
	r.getQuestion()

	question := frame.Question{
		Info: r.question,
		Round: 	r.round,
		Time: 60,
	}

	winnerResult := frame.MatchResult{
		Round: r.round -1,
		Win: 1,
		Question: question,
	}

	winnerFb := frame.Frame{
		Type: "matchresult",
		Code: 1000,
		Data: winnerResult,
	}

	winner.Send <- &winnerFb

	loserResult := frame.MatchResult{
		Round: r.round -1,
		Win: 0,
		Question: question,
	}

	loserFb := frame.Frame{
		Type: "matchresult",
		Code: 1000,
		Data: loserResult,
	}

	for _,client := range r.Clients{
		if client != winner{
			client.Send <- &loserFb
		}
	}

}
func (r *Room)Timeout()  {  // 该次对局超时
	fmt.Println("对局超时")
	r.getQuestion()
	r.round += 1
	r.TimeoutFeedBack()
	r.limitTime.Reset(60* time.Second)
}

func (r *Room) breakRoom(){  // 房间解散
	fmt.Println("房间解散")
	r.limitTime.Stop()
	r.Clients[0].Status = "waiting"
	r.Clients[1].Status = "waiting"
	r.breakRoomFeedBack()
}

func (r *Room)init(){  // 参数初始化
	r.Score = make(map[*Client]int)
	r.round = 1
	r.Score[r.Clients[0]] = 0
	r.Score[r.Clients[1]] = 0
	r.limitTime = time.NewTimer(60 * time.Second)
	r.BreakRoom = make(chan *Client)
	r.Result = make(chan *Client)
}

/**
发送房间解散反馈
 */
func(r *Room) breakRoomFeedBack(){
	fb := &frame.Frame{
		Type: "matchOver",
	}

	r.Clients[0].Send <- fb
	r.Clients[1].Send <- fb
}
/**
发送该回合超时的反馈
 */
func(r *Room) TimeoutFeedBack(){

	data := &frame.MatchResult{
		Round: r.round -1,
		Win: -1,
		Question: frame.Question{
			Round: r.round,
			Info: r.question,
			Time: 60,
		},
	}
	fb := &frame.Frame{
		Type: "matchresult",
		Code: 1000,
		Data: data,
	}

	r.Clients[0].Send <- fb
	r.Clients[1].Send <- fb
}
func (r *Room) matchedFeedBack(){  // 匹配成功的反馈
	p1 := r.Clients[0]
	p2 := r.Clients[1]

	dataP1 := frame.MatchedData{
		OpponentName: p2.Username,
		OpponentId: p2.Id,
		Question: frame.Question{Round:1, Info: r.question, Time:60},
	}

	fbP1 := &frame.Frame{
		Type: "matched",
		Code: 1000,
		Data: dataP1,
	}

	p1.Send <- fbP1

	dataP2 := frame.MatchedData{
		OpponentName: p1.Username,
		OpponentId: p1.Id,
		Question: frame.Question{Round:1, Info: r.question, Time:60},
	}
	fbP2 := &frame.Frame{
		Type: "matched",
		Code: 1000,
		Data: dataP2,
	}

	p2.Send <- fbP2
}