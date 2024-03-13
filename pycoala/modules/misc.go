package modules

import (
	"encoding/json"
	"fmt"
	"log"
	bothttp "pycoala/pycoala/utils/helpers"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
)

var statusEmojis = map[int]string{
	0:  "⛈",
	1:  "⛈",
	2:  "⛈",
	3:  "⛈",
	4:  "⛈",
	5:  "🌨",
	6:  "🌨",
	7:  "🌨",
	8:  "🌨",
	9:  "🌨",
	10: "🌨",
	11: "🌧",
	12: "🌧",
	13: "🌨",
	14: "🌨",
	15: "🌨",
	16: "🌨",
	17: "⛈",
	18: "🌧",
	19: "🌫",
	20: "🌫",
	21: "🌫",
	22: "🌫",
	23: "🌬",
	24: "🌬",
	25: "🌨",
	26: "☁️",
	27: "🌥",
	28: "🌥",
	29: "⛅️",
	30: "⛅️",
	31: "🌙",
	32: "☀️",
	33: "🌤",
	34: "🌤",
	35: "⛈",
	36: "🔥",
	37: "🌩",
	38: "🌩",
	39: "🌧",
	40: "🌧",
	41: "❄️",
	42: "❄️",
	43: "❄️",
	44: "n/a",
	45: "🌧",
	46: "🌨",
	47: "🌩",
}

func getStatusEmoji(statusCode int) string {
	emoji, ok := statusEmojis[statusCode]

	if ok {
		return emoji
	} else {
		return "n/a"
	}
}

func WeatherModule(bot *telego.Bot, update telego.Update) {
	chatID := telegoutil.ID(update.Message.Chat.ID)
	args := update.Message.Text

	if len(args) < 9 {
		bot.SendMessage(&telego.SendMessageParams{
			ChatID:    chatID,
			Text:      "<b>Please provide a location to get the weather.</b>",
			ParseMode: "HTML",
		})
		return
	}

	location := args[8:]

	params_Location := bothttp.RequestGETParams{
		Query: map[string]string{
			"apiKey":   "8de2d8b3a93542c9a2d8b3a935a2c909",
			"format":   "json",
			"language": "en",
			"query":    location,
		},
	}

	req_Location := bothttp.RequestGET("https://api.weather.com/v3/location/search", params_Location)
	body_Location := req_Location.Body()

	var weather_Location_Data map[string]interface{}
	err_Location := json.Unmarshal(body_Location, &weather_Location_Data)
	if err_Location != nil {
		log.Println("Error decoding JSON in WeatherModule:", err_Location)
		bot.SendMessage(&telego.SendMessageParams{
			ChatID:    chatID,
			Text:      "<b>Sorry, there was an error retrieving weather information.</b>",
			ParseMode: "HTML",
		})
		return
	}

	// Access first location information
	if locationData, ok := weather_Location_Data["location"].(map[string]interface{}); ok { // Check for "location" key
		// Extract address, latitude, and longitude from the location data
		if address, ok := locationData["address"]; ok {
			if lat, ok := locationData["latitude"]; ok {
				if lon, ok := locationData["longitude"]; ok {
					// Declare vars =>
					addressFirst := address.([]interface{})[0]
					latFirst := lat.([]interface{})[0]
					lonFirst := lon.([]interface{})[0]

					params_Weather := bothttp.RequestGETParams{
						Query: map[string]string{
							"apiKey":   "8de2d8b3a93542c9a2d8b3a935a2c909",
							"format":   "json",
							"language": "en",
							"geocode":  fmt.Sprintf("%.3f,%.3f", latFirst, lonFirst),
							"units":    "m",
						},
					}

					req_Weather := bothttp.RequestGET("https://api.weather.com/v3/aggcommon/v3-wx-observations-current", params_Weather)

					body_Weather := req_Weather.Body()

					var weatherData map[string]interface{}
					err_Weather := json.Unmarshal(body_Weather, &weatherData)
					if err_Weather != nil {
						log.Println("Error decoding JSON in WeatherModule:", err_Weather)
						bot.SendMessage(&telego.SendMessageParams{
							ChatID:    chatID,
							Text:      "<b>Sorry, there was an error retrieving weather information.</b>",
							ParseMode: "HTML",
						})
						return
					}
					if observations_wx, ok := weatherData["v3-wx-observations-current"].(map[string]interface{}); ok {
						temperature := observations_wx["temperature"]
						feelsLike := observations_wx["temperatureFeelsLike"]
						airHumidity := observations_wx["relativeHumidity"]
						windSpeed := observations_wx["windSpeed"]
						iconCode := observations_wx["iconCode"]
						weatherType := observations_wx["wxPhraseLong"]
						intCode := int(iconCode.(float64))
						bot.SendMessage(&telego.SendMessageParams{
							ChatID: chatID,
							Text: fmt.Sprintf("<b>%s</b>:\n\n<b><i>📍 Location Info:</i></b>\n<b>Latitude:</b> <code>%.3f</code>\n<b>Longitude:</b> <code>%.3f</code>\n\n<b><i>%s Weather — %s:</i></b>\n<b>Temperature:</b> <code>%.0f °C</code>\n<b>Thermal sensation:</b> <code>%.0f °C</code>\n<b>Air umidity:</b> <code>%.0f %%</code>\n<b>Wind:</b> <code>%.0f km/h</code>",
								addressFirst, latFirst, lonFirst, getStatusEmoji(intCode), weatherType, temperature, feelsLike, airHumidity, windSpeed),
							ParseMode: "HTML",
						})
						return
					}
				}
			}
		} else {
			bot.SendMessage(&telego.SendMessageParams{
				ChatID:    chatID,
				Text:      "<b>Sorry, there was an error retrieving weather information.</b>",
				ParseMode: "HTML",
			})
			return
		}
	} else {
		bot.SendMessage(&telego.SendMessageParams{
			ChatID:    chatID,
			Text:      "<b>Location not found!</b>",
			ParseMode: "HTML",
		})
		return
	}
}
