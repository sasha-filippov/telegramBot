package telegramClient

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

// As message can be empty => using pointer for nil values
type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}
type From struct {
	UserName string `json:"username"`
}
type Chat struct {
	ID int `json:"id"`
}
