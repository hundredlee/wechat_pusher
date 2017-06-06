package statics

const (
	WECHAT_GET_ACCESS_TOKEN = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	WECHAT_TEMPLATE_SEND = "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s"
)