package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/HalalBot/ocr"

	errorHand "github.com/HalalBot/error"

	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)

	errorHand.HandleError(err)

	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
						log.Print(err)
					}
				case *linebot.ImageMessage:
					image, err := bot.GetMessageContent(message.ID).Do()
					defer image.Content.Close()
					errorHand.HandleError(err)

					imageFile, err := os.Create("./static/img/img.png")
					errorHand.HandleError(err)

					_, err = io.Copy(imageFile, image.Content)

					files, _ := ioutil.ReadDir("./static/img/")
					for _, f := range files {
						log.Print(f.Name())
					}

					ocr.PosOCR()
					log.Print("2")

					//originalURL := "https://pbs.twimg.com/media/ELWG8dcU8AAG1Hi.jpg:small " //"https://halal-bot.herokuapp.com/static/img/sample.jpeg"
					//previewURL := "https://pbs.twimg.com/media/ELWG8dcU8AAG1Hi.jpg:small"   //"https://halal-bot.herokuapp.com/static/img/sample.jpeg"
					//
					//if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage(originalURL, previewURL)).Do(); err != nil {
					//	log.Print(err)
					//}
				}
			}
		}
	})
	err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	errorHand.HandleError(err)

}
