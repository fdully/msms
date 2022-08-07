package mtssms

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

var (
	MaxMsgLen = 1000
	MinMsgLen = 5
)

func CreateRequestWithSingleMessage(phones []string, message, signuture string) *MTSRequest {
	m := &MTSRequest{}
	m.Options.From.SmsAddress = signuture

	mm := Messages{}

	mm.Content = Content{
		ShortText: message,
	}

	for _, v := range phones {
		to := To{
			Msisdn: v,
		}
		mm.To = append(mm.To, to)
	}

	m.Messages = append(m.Messages, mm)

	return m
}

func CreateInfoRequest(id string) *MTSInfoRequest {
	return &MTSInfoRequest{
		IntIds: []string{id},
	}
}

func CheckMessage(msg string) error {

	msgLen := utf8.RuneCountInString(msg)

	if msgLen > MaxMsgLen || msgLen < MinMsgLen {
		return fmt.Errorf("wrong message length %d", msgLen)
	}

	return nil
}

func CheckSingleTelephone(tel string) (string, error) {
	r := strings.NewReplacer("-", "", "(", "", ")", "", " ", "", "+", "")
	tel = r.Replace(tel)

	telLen := utf8.RuneCountInString(tel)

	if telLen < 11 || telLen > 18 {
		return "", fmt.Errorf("length of telephone is wrong %s", tel)
	}

	return tel, nil
}
