package ocr

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strings"

	errorHand "github.com/HalalBot/error"
)

type ResultsWords struct {
	BoundingBox [][]float64
	Text        string
	Confidence  float64
	LineBreak   bool
}

type Results struct {
	Words []ResultsWords `json:"words"`
}

var ngFood = []string{"ワイン", "みりん", "日本酒", "ビール", "ラム酒", "料理酒", "豚", "ポーク", "ゼラチン", "ラード"}

func foodContain(text string) bool {
	for _, ng := range ngFood {
		if strings.Contains(text, ng) {
			return true
		}
	}
	return false
}

func PosOCR() {
	url := "https://ocr-devday19.linebrain.ai/v1/recognition"
	entrance := "detection"
	fieldName := "image"
	imageName := "./static/img/img.png"
	image, err := os.Open(imageName)
	errorHand.HandleError(err)

	log.Print(reflect.TypeOf(image))

	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)

	fw, err := mw.CreateFormField(fieldName)
	//fw, err := mw.CreateFormFile(fieldName, "")
	_, err = io.Copy(fw, image)
	errorHand.HandleError(err)

	err = mw.WriteField("entrance", entrance)
	errorHand.HandleError(err)

	err = mw.WriteField("segments", "false")
	errorHand.HandleError(err)

	contentType := mw.FormDataContentType()

	err = mw.Close()
	errorHand.HandleError(err)

	req, err := http.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("X-ClovaOCR-Service-ID", os.Getenv("ClovaOCR_ID"))

	resp, err := http.DefaultClient.Do(req)
	errorHand.HandleError(err)
	defer resp.Body.Close()

	jsonBody, err := ioutil.ReadAll(resp.Body)

	var result Results
	err = json.Unmarshal(jsonBody, &result)
	errorHand.HandleError(err)

	log.Print(result.Words)

	for num, word := range result.Words {
		if foodContain(word.Text) {
			log.Print("NG")
			break
		} else if num == len(result.Words)-1 {
			log.Print("OK")
		}
	}
}
