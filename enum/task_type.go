package enum


const (
	TASK_TYPE_TEMPLATE = "TYPE_TEMPLATE"
	TASK_TYPE_TEXT_CUSTOM = "TYPE_TEXT_CUSTOM"
	TASK_TYPE_IMAGE = "TYPE_IMAGE"
	TASK_TYPE_WORD = "TYPE_WORD"
	TASK_TYPE_IMAGE_WORD = "TYPE_IMAGE_WORD"
)


var URL_MAP map[string]string = map[string]string {
	TASK_TYPE_TEMPLATE:"https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s",
	TASK_TYPE_TEXT_CUSTOM:"https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=%s"}


