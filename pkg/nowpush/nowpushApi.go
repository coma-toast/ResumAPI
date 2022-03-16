package nowpush

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Service is the service
type Service interface {
	GetUser()
	SendMessage(message_type string, note string, link string) (MessageResponse, error)
}

func (c Client) GetUser() (User, error) {
	var err error
	var results User
	var url = "getUser"
	_, err = c.call(http.MethodGet, url, nil, nil, &results)
	if err != nil {
		return User{}, err
	}

	return results, err
}

func (c Client) SendMessage(message_type string, note string, link string) (MessageResponse, error) {
	var err error
	var results MessageResponse
	var url = "sendMessage"
	data := Message{
		MessageType: message_type,
		Note:        note,
		DeviceType:  "POD Tool",
		URL:         link,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return MessageResponse{}, err
	}
	_, err = c.call(http.MethodPost, url, jsonData, nil, &results)
	if err != nil {
		return MessageResponse{}, err
	}
	log.WithFields(log.Fields{
		"MessageResponse": results,
	}).Info("NowPush message sent")

	return results, err
}
