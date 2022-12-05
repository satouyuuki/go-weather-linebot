package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func ToJstFromTimestamp(ts int) time.Time {
	i, err := strconv.ParseInt(fmt.Sprint(ts), 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)

	return tm.In(jst)
}

func ParseRequest(channelSecret string, r events.LambdaFunctionURLRequest) ([]*linebot.Event, error) {
	if !validateSignature(channelSecret, r.Headers["x-line-signature"], []byte(r.Body)) {
		return nil, linebot.ErrInvalidSignature
	}

	request := &struct {
		Events []*linebot.Event `json:"events"`
	}{}
	if err := json.Unmarshal([]byte(r.Body), request); err != nil {
		return nil, err
	}
	return request.Events, nil

}

func validateSignature(channelSecret, signature string, body []byte) bool {
	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		log.Fatal(err)
		return false
	}
	hash := hmac.New(sha256.New, []byte(channelSecret))

	_, err = hash.Write(body)
	if err != nil {
		log.Fatal(err)
		return false
	}

	return hmac.Equal(decoded, hash.Sum(nil))
}
