package modules

import (
	"encoding/json"
	"fmt"

	bothttp "pycoala/pycoala/utils/helpers"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
)

const CodeName_List = "https://raw.githubusercontent.com/androidtrackers/certified-android-devices/master/by_device.json"

func DeviceInfo(bot *telego.Bot, update telego.Update) {
	ChatID := telegoutil.ID(update.Message.Chat.ID)
	args := update.Message.Text
	params := bothttp.RequestGETParams{}
	cdevices := bothttp.RequestGET(CodeName_List, params)
	// Get a json from response.
	body_cdevices := cdevices.Body()
	var cdevices_Data map[string]interface{}
	json.Unmarshal(body_cdevices, &cdevices_Data)
	// Initializate a checks.
	if len(args) < 7 {
		bot.SendMessage(&telego.SendMessageParams{
			ChatID:    ChatID,
			Text:      "<b>Device not specified!</b>",
			ParseMode: "HTML",
		},
		)
		return
	}

	device := args[7:]

	bot.SendMessage(&telego.SendMessageParams{
		ChatID:    ChatID,
		Text:      fmt.Sprintf("<b>Search Device:</b> %s", device),
		ParseMode: "HTML",
	},
	)
}
