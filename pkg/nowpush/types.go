package nowpush

import "time"

type User struct {
	Message interface{} `json:"message"`
	Data    struct {
		AuthToken []string  `json:"auth_token"`
		IsPro     bool      `json:"is_pro"`
		ID        string    `json:"_id"`
		CreateAt  time.Time `json:"create_at"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		MyID      string    `json:"my_id"`
		ImageURL  string    `json:"image_url"`
		V         int       `json:"__v"`
	} `json:"data"`
	IsError bool `json:"isError"`
}
type Message struct {
	MessageType string `json:"message_type"`
	Note        string `json:"note"`
	DeviceType  string `json:"device_type"`
	URL         string `json:"url"`
}

type MessageResponse struct {
	Msg  string `json:"msg"`
	Data struct {
		ID          string    `json:"_id"`
		MsgAt       time.Time `json:"msg_at"`
		ChatRoom    string    `json:"chat_room"`
		DeviceType  string    `json:"device_type"`
		MessageType string    `json:"message_type"`
		Message     string    `json:"message"`
		Note        string    `json:"note"`
		IsEncrypted bool      `json:"is_encrypted"`
	} `json:"data"`
	IsError bool `json:"isError"`
}
