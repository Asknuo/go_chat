package utlis

const (
	SUCCESS               = 200
	UpdatePasswordSuccess = 201   //密码成功
	NotExistInentifier    = 202   //未绑定
	ERROR                 = 500   //失败
	InvalidParams         = 400   //请求参数错误
	ErrorDatabase         = 40001 //数据库操作错误

	WebsocketSuccessMessage = 50001 //解析content内容信息
	WebsocketSuccess        = 50002 //请求历史纪录操作成功
	WebsocketEnd            = 50003 //请求没有更多历史记录
	WebsocketOnlineReply    = 50004 //在线应答
	WebsocketOfflineReply   = 50005 //离线回答
	WebsocketLimit          = 50006 //请求受到限制
)

var MsgFlags = map[int]string{
	SUCCESS:                 "ok",
	UpdatePasswordSuccess:   "修改密码成功",
	NotExistInentifier:      "该第三方账号未绑定",
	ERROR:                   "失败",
	InvalidParams:           "请求参数错误",
	ErrorDatabase:           "数据库操作出错,请重试",
	WebsocketSuccessMessage: "解析content内容信息",
	WebsocketSuccess:        "发送信息，请求历史纪录操作成功",
	WebsocketEnd:            "请求历史纪录，但没有更多记录了",
	WebsocketOnlineReply:    "针对回复信息在线应答成功",
	WebsocketOfflineReply:   "针对回复信息离线回答成功",
	WebsocketLimit:          "请求受到限制",
}

// GetMsg 获取状态码对应信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
