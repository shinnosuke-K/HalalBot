package ocr

import (
	"io"
	"io/ioutil"
	"log"

	errorHand "github.com/HalalBot/error"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/rekognition"

	"github.com/aws/aws-sdk-go/aws/session"
)

func DoOCR(imageContent io.ReadCloser) {

	bytes, err := ioutil.ReadAll(imageContent)
	errorHand.HandleError(err)

	sess := session.Must(session.NewSession())

	svc := rekognition.New(sess, aws.NewConfig().WithRegion("ap-northeast-1"))

	params := &rekognition.DetectTextInput{
		Image: &rekognition.Image{
			Bytes: bytes,
		},
	}

	resp, err := svc.DetectText(params)
	errorHand.HandleError(err)

	log.Print(resp)
}
