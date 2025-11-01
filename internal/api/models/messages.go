package models

type Message struct {
	Id                   string `json:"id"`
	Content              string `json:"content"`
	RecipientPhoneNumber string `json:"recipient_phone_number"`
	Status               string `json:"status"`
	CreatedAt            string `json:"created_at"`
}

type MessageResponse struct {
	Count    int       `json:"count"`
	Messages []Message `json:"messages"`
}
