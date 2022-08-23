package notifications

import (
	"fmt"

	"github.com/coma-toast/notifapi/pkg/client"
	"github.com/coma-toast/notifapi/pkg/notification"
)

type NotificationAPI struct {
	Target string
}

func (n *NotificationAPI) SendMessage(title, body string) error {
	client := client.Client{Target: n.Target}
	message := notification.Message{
		Interests: []string{"hello"},
		Title:     title,
		Body:      body,
		Source:    "ResumAPI", // what application are you sending this from
	}

	response, err := client.SendMessage(message)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(response.Status)

	return nil
}
