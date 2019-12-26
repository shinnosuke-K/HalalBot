package ocr

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	errorHand "github.com/HalalBot/error"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/rekognition"

	"github.com/aws/aws-sdk-go/aws/session"
)

func DoOCR(imageContent io.ReadCloser) []string {

	bytes, err := ioutil.ReadAll(imageContent)
	errorHand.HandleError(err)

	sess := session.Must(session.NewSession())

	svc := rekognition.New(sess, aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")))

	params := &rekognition.DetectTextInput{
		Image: &rekognition.Image{
			Bytes: bytes,
		},
	}

	resp, err := svc.DetectText(params)
	errorHand.HandleError(err)
	var ocrText []string
	for _, text := range resp.TextDetections {
		ocrText = append(ocrText, *text.DetectedText)
		log.Print(*text.DetectedText)
	}

	return ocrText
}
