package models

type Message struct {
	ToUser string `json:"touser"`
	TemplateId string `json:"template_id"`
	Url string `json:"url"`
	Data Data `json:"data"`
}

type Data struct {
	First Raw `json:"first"`
	Subject Raw `json:"subject"`
	Sender Raw `json:"sender"`
	Remark Raw `json:"remark"`
}

type Raw struct {
	Value string `json:"value"`
	Color string `json:"color"`
}


