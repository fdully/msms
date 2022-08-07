package mtssms

import "time"

// type MTSRequest struct {
// 	Messages []struct {
// 		Content struct {
// 			ShortText string `json:"short_text"`
// 		} `json:"content"`
// 		To []struct {
// 			Msisdn    string `json:"msisdn"`
// 			MessageID string `json:"message_id"`
// 		} `json:"to"`
// 	} `json:"messages"`
// 	Options struct {
// 		From struct {
// 			SmsAddress string `json:"sms_address"`
// 		} `json:"from"`
// 	} `json:"options"`
// }

type MTSResponse struct {
	Messages []struct {
		Msisdn     string `json:"msisdn"`
		Email      string `json:"email"`
		MessageID  string `json:"message_id"`
		InternalID string `json:"internal_id"`
	} `json:"messages"`
}

type MTSRequest struct {
	Messages []Messages `json:"messages"`
	Options  Options    `json:"options"`
}
type Content struct {
	ShortText string `json:"short_text"`
}
type To struct {
	Msisdn    string `json:"msisdn"`
	MessageID string `json:"message_id"`
}
type Messages struct {
	Content Content `json:"content"`
	To      []To    `json:"to"`
}
type From struct {
	SmsAddress string `json:"sms_address"`
}
type Options struct {
	From From `json:"from"`
}

type MTSInfoRequest struct {
	IntIds []string `json:"int_ids"`
}

type MTSInfoResponse struct {
	EventsInfo []struct {
		EventsInfo []struct {
			Channel        int         `json:"channel"`
			Destination    string      `json:"destination"`
			EventAt        time.Time   `json:"event_at"`
			InternalErrors interface{} `json:"internal_errors"`
			InternalID     string      `json:"internal_id"`
			MessageID      string      `json:"message_id"`
			Naming         string      `json:"naming"`
			ReceivedAt     time.Time   `json:"received_at"`
			Status         int         `json:"status"`
			TotalParts     int         `json:"total_parts"`
		} `json:"events_info"`
		Key string `json:"key"`
	} `json:"events_info"`
}
