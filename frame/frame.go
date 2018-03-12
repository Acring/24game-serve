package frame

type Frame struct{
	Type string `json:"type"`
	Code int	`json:"code"`
	Data interface{} `json:"data"`
}

type MatchingData struct {  // 匹配中信息
	OnLineNum int `json:"onlinenum"`
}

type MatchedData struct{  // 匹配成功信息
	OpponentName string `json:"opponentname"`
	OpponentId string `json:"opponentid"`
	Question Question `json:"question"`
}

type Question struct{  // 题目信息
	Round int `json:"round"`
	Info [4]int `json:"info"`
	Time int `json:"time"`
}

type MatchResult struct {  // 比赛结果
	Round int `json:"round"`// 回合
	Win int `json:"win"`  	// 判定比赛的输赢, win=1代表胜利, 0代表失败, -1代表超时
	Question Question `json:"question"`  // 下一轮的题目
}