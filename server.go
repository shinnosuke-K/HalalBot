package main

import (
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

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
					log.Print(message)
					if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
						log.Print(err)
					}
				case *linebot.ImageMessage:
					image, err := bot.GetMessageContent(message.ID).Do()
					defer image.Content.Close()
					if err != nil {
						log.Print(err)
					}
					imgFile, err := os.OpenFile("static/img/"+message.ID+"jpg", os.O_CREATE, 0600)
					buf := make([]byte, 1024)
					for {
						n, err := image.Content.Read(buf)
						if n == 0 || err != nil {
							break
						}
						imgFile.Write(buf[:n])
					}
					log.Print("saved images!")
				}
			}
		}
	})

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}

}
