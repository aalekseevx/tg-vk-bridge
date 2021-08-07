package bridge

type Post struct {
	PersonalChatID int64 `json:"personal_chat_id"`
	ConfirmationMessageID int `json:"confirmation_message_id"`
	ForwardedMessagesIds []int `json:"forwarded_messages_ids"`
	Text string `json:"text"`
	TgAttachments []string `json:"tg_attachments"`
	VkAttachments []string `json:"vk_attachments"`
}

type Creator struct {
	ChatID   int `gorm:"primaryKey"`
	Token    string
}

type Channel struct {
	ChatID        int64 `gorm:"primaryKey" json:"chat_id"`
	Name          string `json:"name"`
	CreatorID     int
	Creator       Creator
	GroupID       int
	Delay         int
	VkAccessToken string
}
