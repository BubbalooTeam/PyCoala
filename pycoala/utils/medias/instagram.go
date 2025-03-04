package medias

import (
	"encoding/json"
	"fmt"
	"log"
	"pycoala/pycoala/utils/helpers"
	"regexp"
	"strings"

	"github.com/mymmrac/telego/telegoutil"
)

type InstagramData struct {
	ShortcodeMedia ShortcodeMedia `json:"shortcode_media"`
	Data           struct {
		XDTShortcodeMedia ShortcodeMedia `json:"xdt_shortcode_media"`
	} `json:"data,omitempty"`
}

type ShortcodeMedia struct {
	Typename              string                `json:"__typename"`
	ID                    string                `json:"id"`
	Shortcode             string                `json:"shortcode"`
	Dimensions            Dimensions            `json:"dimensions"`
	DisplayResources      []DisplayResources    `json:"display_resources"`
	IsVideo               bool                  `json:"is_video"`
	Title                 string                `json:"title"`
	VideoURL              string                `json:"video_url"`
	EdgeMediaToCaption    EdgeMediaToCaption    `json:"edge_media_to_caption"`
	EdgeSidecarToChildren EdgeSidecarToChildren `json:"edge_sidecar_to_children"`
}

type Dimensions struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

type DisplayResources struct {
	ConfigWidth  int    `json:"config_width"`
	ConfigHeight int    `json:"config_height"`
	Src          string `json:"src"`
}

type EdgeMediaToCaption struct {
	Edges []struct {
		Node struct {
			Text string `json:"text"`
		} `json:"node"`
	} `json:"edges"`
}
type EdgeSidecarToChildren struct {
	Edges []struct {
		Node struct {
			Typename         string             `json:"__typename"`
			ID               string             `json:"id"`
			Shortcode        string             `json:"shortcode"`
			CommenterCount   int                `json:"commenter_count"`
			Dimensions       Dimensions         `json:"dimensions"`
			DisplayResources []DisplayResources `json:"display_resources"`
			IsVideo          bool               `json:"is_video"`
			VideoURL         string             `json:"video_url,omitempty"`
			DisplayURL       string             `json:"display_url"`
		} `json:"node"`
	} `json:"edges"`
}

func (dm *DownloadMedia) Instagram(url string) {
	var instagramData InstagramData
	var postID string

	if matches := (regexp.MustCompile((`(?:reel(?:s?)|p)/([A-Za-z0-9_-]+)`))).FindStringSubmatch(url); len(matches) == 2 {
		postID = matches[1]
	} else {
		return
	}

	headers := map[string]string{
		"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"accept-language":           "en-US,en;q=0.9",
		"cache-control":             "max-age=0",
		"connection":                "close",
		"sec-fetch-mode":            "navigate",
		"upgrade-insecure-requests": "1",
		"referer":                   "https://www.instagram.com/",
		"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
		"viewport-width":            "1280",
	}

	body := helpers.RequestGET(fmt.Sprintf("https://www.instagram.com/p/%v/embed/captioned/", postID), helpers.RequestGETParams{Headers: headers}).Body()
	match := (regexp.MustCompile(`\\\"gql_data\\\":([\s\S]*)\}\"\}`)).FindSubmatch(body)
	if len(match) == 2 {
		s := strings.ReplaceAll(string(match[1]), `\"`, `"`)
		s = strings.ReplaceAll(s, `\\/`, `/`)
		s = strings.ReplaceAll(s, `\\`, `\`)

		err := json.Unmarshal([]byte(s), &instagramData)
		if err != nil {
			log.Println(err)
		}
		dm.Caption = instagramData.ShortcodeMedia.EdgeMediaToCaption.Edges[0].Node.Text

		switch instagramData.ShortcodeMedia.Typename {
		case "GraphVideo":
			dm.MediaItems = append(dm.MediaItems, telegoutil.MediaVideo(
				telegoutil.FileFromURL(instagramData.ShortcodeMedia.VideoURL)).WithWidth(instagramData.ShortcodeMedia.Dimensions.Width).WithHeight(instagramData.ShortcodeMedia.Dimensions.Height),
			)
		case "GraphSidecar":
			for _, results := range instagramData.ShortcodeMedia.EdgeSidecarToChildren.Edges {
				if !results.Node.IsVideo {
					dm.MediaItems = append(dm.MediaItems, telegoutil.MediaPhoto(
						telegoutil.FileFromURL(results.Node.DisplayResources[len(results.Node.DisplayResources)-1].Src)),
					)
				} else {
					dm.MediaItems = append(dm.MediaItems, telegoutil.MediaVideo(
						telegoutil.FileFromURL(results.Node.VideoURL)).WithWidth(results.Node.DisplayResources[len(results.Node.DisplayResources)-1].ConfigWidth).WithHeight(results.Node.DisplayResources[len(results.Node.DisplayResources)-1].ConfigHeight),
					)
				}
			}
		}
	}

	mediaTypeData := regexp.MustCompile(`(?s)data-media-type="(.*?)"`).FindAllStringSubmatch(string(body), -1)
	if len(mediaTypeData) > 0 && len(mediaTypeData[0]) > 1 && mediaTypeData[0][1] == "GraphImage" {
		// Get the main media
		re := regexp.MustCompile(`class="Content(.*?)src="(.*?)"`)
		mainMediaData := re.FindAllStringSubmatch(string(body), -1)
		mainMediaURL := (strings.ReplaceAll(mainMediaData[0][2], "amp;", ""))

		// Get the caption
		re = regexp.MustCompile(`(?s)class="Caption"(.*?)class="CaptionUsername"(.*?)<\/a>(.*?)<div`)
		captionData := re.FindAllStringSubmatch(string(body), -1)
		if len(captionData) > 0 && len(captionData[0]) > 2 {
			re = regexp.MustCompile(`<[^>]*>`)
			dm.Caption = strings.TrimSpace(re.ReplaceAllString(captionData[0][3], ""))
		}

		dm.MediaItems = append(dm.MediaItems, telegoutil.MediaPhoto(
			telegoutil.FileFromURL(mainMediaURL)),
		)
	}

	if len(dm.MediaItems) == 0 {
		headers := map[string]string{
			`User-Agent`:         `Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0`,
			`Accept`:             `*/*`,
			`Accept-Language`:    `pt-BR,pt;q=0.8,en-US;q=0.5,en;q=0.3`,
			`Content-Type`:       `application/x-www-form-urlencoded`,
			`X-FB-Friendly-Name`: `PolarisPostActionLoadPostQueryQuery`,
			`X-CSRFToken`:        `-m5n6c-w1Z9RmrGqkoGTMq`,
			`X-IG-App-ID`:        `936619743392459`,
			`X-FB-LSD`:           `AVp2LurCmJw`,
			`X-ASBD-ID`:          `129477`,
			`DNT`:                `1`,
			`Sec-Fetch-Dest`:     `empty`,
			`Sec-Fetch-Mode`:     `cors`,
			`Sec-Fetch-Site`:     `same-origin`,
		}

		params := []string{
			`av=0`,
			`__d=www`,
			`__user=0`,
			`__a=1`,
			`__req=3`,
			`__hs=19734.HYP:instagram_web_pkg.2.1..0.0`,
			`dpr=1`,
			`__ccg=UNKNOWN`,
			`__rev=1010782723`,
			`__s=qg5qgx:efei15:ng6310`,
			`__hsi=7323030086241513400`,
			`__dyn=7xeUjG1mxu1syUbFp60DU98nwgU29zEdEc8co2qwJw5ux609vCwjE1xoswIwuo2awlU-cw5Mx62G3i1ywOwv89k2C1Fwc60AEC7U2czXwae4UaEW2G1NwwwNwKwHw8Xxm16wUxO1px-0iS2S3qazo7u1xwIwbS1LwTwKG1pg661pwr86C1mwrd6goK68jxe6V8`,
			`__csr=gps8cIy8WTDAqjWDrpda9SoLHhaVeVEgvhaJzVQ8hF-qEPBV8O4EhGmciDBQh1mVuF9V9d2FHGicAVu8GAmfZiHzk9IxlhV94aKC5oOq6Uhx-Ku4Kaw04Jrx64-0oCdw0MXw1lm0EE2Ixcjg2Fg1JEko0N8U421tw62wq8989EMw1QpV60CE02BIw`,
			`__comet_req=7`,
			`lsd=AVp2LurCmJw`,
			`jazoest=2989`,
			`__spin_r=1010782723`,
			`__spin_b=trunk`,
			`__spin_t=1705025808`,
			`fb_api_caller_class=RelayModern`,
			`fb_api_req_friendly_name=PolarisPostActionLoadPostQueryQuery`,
			fmt.Sprintf(`variables={"shortcode": "%v","fetch_comment_count":2,"fetch_related_profile_media_count":0,"parent_comment_count":0,"child_comment_count":0,"fetch_like_count":10,"fetch_tagged_user_count":null,"fetch_preview_comment_count":2,"has_threaded_comments":true,"hoisted_comment_id":null,"hoisted_reply_id":null}`, postID),
			`server_timestamps=true`,
			`doc_id=10015901848480474`,
		}

		body := helpers.RequestPOST("https://www.instagram.com/api/graphql", helpers.RequestPOSTParams{Headers: headers, BodyString: params}).Body()
		err := json.Unmarshal(body, &instagramData)
		if err != nil {
			log.Println(err)
		}
		result := instagramData.Data.XDTShortcodeMedia
		dm.Caption = result.EdgeMediaToCaption.Edges[0].Node.Text
		if strings.Contains(result.Typename, "Video") {
			dm.MediaItems = append(dm.MediaItems, telegoutil.MediaVideo(
				telegoutil.FileFromURL(result.VideoURL)).WithWidth(result.Dimensions.Width).WithHeight(result.Dimensions.Height),
			)
		}
	}
}
