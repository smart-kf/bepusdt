package dao

type Message struct {
	MsgType     string `json:"msgType"`     // text || image || video
	MsgId       string `json:"msgId"`       // 消息id
	GuestName   string `json:"guestName"`   // 客户名称
	GuestAvatar string `json:"guestNvatar"` // 客户头像
	KfName      string `json:"kfName"`      // 客服名称
	KfAvatar    string `json:"kfAvatar"`    // 客服头像
	Content     string `json:"content"`     // 具体消息内容
	Ip          string `json:"ip"`          // 客户IP
	IsFromKf    bool   `json:"isFromKf"`    // 是否是来自客服发送
}
