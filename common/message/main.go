package message

import (
	"fmt"

	"github.com/9688101/hx-admin/global"
)

const (
	ByAll           = "all"
	ByEmail         = "email"
	ByMessagePusher = "message_pusher"
)

func Notify(by string, title string, description string, content string) error {
	if by == ByEmail {
		return SendEmail(title, global.RootUserEmail, content)
	}
	if by == ByMessagePusher {
		return SendMessage(title, description, content)
	}
	return fmt.Errorf("unknown notify method: %s", by)
}
