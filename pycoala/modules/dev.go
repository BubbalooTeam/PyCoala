package modules

import (
	"fmt"
	"log"
	"time"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
)

func PingModule(bot *telego.Bot, update telego.Update) {
	startTime := time.Now()
	msg, err := bot.SendMessage(&telego.SendMessageParams{
		ChatID:    telegoutil.ID(update.Message.Chat.ID),
		Text:      "<b>Pong!</b>",
		ParseMode: "HTML",
	},
	)
	if err != nil {
		log.Fatal("Telego is instable... Please try again!")
		return
	}
	endTime := time.Now()

	finalTime := endTime.Sub(startTime).Abs().Milliseconds()

	bot.EditMessageText(&telego.EditMessageTextParams{
		ChatID:    telegoutil.ID(update.Message.Chat.ID),
		MessageID: msg.MessageID,
		Text:      fmt.Sprintf("<b>Pong!</b> <code>%d ms</code>", finalTime),
		ParseMode: "HTML",
	})
}
