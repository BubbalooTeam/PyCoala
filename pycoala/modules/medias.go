package modules

import (
	"pycoala/pycoala/utils/medias"
	"regexp"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
)

func MediaDownloader(bot *telego.Bot, message telego.Message) {
	// Extract URL from the message text using regex
	url := regexp.MustCompile(`(?:htt.*?//)?(:?.*)?(?:instagram|twitter|x|tiktok|threads)\.(?:com|net)\/(?:\S*)`).FindStringSubmatch(message.Text)
	if len(url) < 1 {
		bot.SendMessage(telegoutil.Message(
			telegoutil.ID(message.Chat.ID),
			"No URL found",
		))
		return
	}

	dm := medias.NewDownloadMedia()
	mediaItems, caption := dm.Download(url[0])

	// Check if only one photo is present and link preview is enabled, then return
	if len(mediaItems) == 1 && mediaItems[0].MediaType() == "photo" && !message.LinkPreviewOptions.IsDisabled {
		return
	}

	if len(mediaItems) > 0 {
		for _, media := range mediaItems[:1] {
			switch media.MediaType() {
			case "photo":
				if photo, ok := media.(*telego.InputMediaPhoto); ok {
					photo.WithCaption(caption).WithParseMode("HTML")
				}
			case "video":
				if video, ok := media.(*telego.InputMediaVideo); ok {
					video.WithCaption(caption).WithParseMode("HTML")
				}
			}
		}
	}

	bot.SendMediaGroup(telegoutil.MediaGroup(
		telegoutil.ID(message.Chat.ID),
		mediaItems...,
	))
}
