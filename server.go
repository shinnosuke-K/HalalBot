package main

import (
	"log"
	"net/http"
	"os"

	"github.com/shinnosuke-K/HalalBot/ocr"

	errorHand "github.com/shinnosuke-K/HalalBot/error"

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

					replyMess := ocr.DoOCR(image.Content)

					if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMess)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})
	err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	errorHand.HandleError(err)
}
