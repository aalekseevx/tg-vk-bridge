package bridge

import (
	//"encoding/json"
	//"fmt"
	"github.com/joeyave/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	//"strings"
	//"time"
)

var bot *tgbotapi.BotAPI

func setWebhook(webhook string) {
	_, err := bot.SetWebhook(tgbotapi.NewWebhookWithCert("https://www.google.com:8443/"+bot.Token, "cert.pem"))
	if err != nil {
		log.Fatal(err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	updates := bot.ListenForWebhook("/" + bot.Token)
}

func setupBot() tgbotapi.UpdatesChannel {
	var err error

	err = tgbotapi.SetLogger(log.WithField("service", "bot"))
	if err != nil {
		log.Panic(err)
	}
	bot, err = tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug, err = strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		log.Panic(err)
	}

	log.WithFields(logrus.Fields{
		"account": bot.Self.UserName,
		"debug":   bot.Debug,
	}).Info("Authorized.")

	return bot.GetUpdatesChan(tgbotapi.NewUpdate(0))
}

func getChatCreator(config tgbotapi.ChatConfig) (*tgbotapi.User, error) {
	administrators, err := bot.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{ChatConfig: config})
	if err != nil {
		return new(tgbotapi.User), err
	}
	var creator *tgbotapi.User = nil
	for _, member := range administrators {
		if member.Status == "creator" {
			creator = member.User
		}
	}
	return creator, nil
}

func handleChannelPost(post *tgbotapi.Message) {
	//creator, err := getChatCreator(post.Chat.ChatConfig())
	//if creator == nil {
	//	log.Fatal("Creator is not available")
	//}
	//
	//channel := getChannel(post.Chat.ID)
	//
	//var attachment string
	//var exists int64
	//
	//if post.Photo != nil {
	//	var biggestFileID string
	//	biggestFileSize := 0
	//	for _, photo := range post.Photo {
	//		if photo.FileSize > biggestFileSize {
	//			biggestFileSize = photo.FileSize
	//			biggestFileID = photo.FileID
	//		}
	//	}
	//
	//	if len(post.MediaGroupID) > 0 {
	//		exists, err = rdsAttachmentFileIDS.Exists(ctx, post.MediaGroupID).Result()
	//		if err != nil {
	//			return
	//		}
	//
	//		rdsAttachmentFileIDS.RPush(ctx, post.MediaGroupID, biggestFileID)
	//	}
	//
	//	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: biggestFileID})
	//	if err != nil {
	//		return
	//	}
	//
	//	url := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", os.Getenv("TOKEN"), file.FilePath)
	//	path := fmt.Sprintf("tmp/%s.jpg", file.FileID)
	//	err = DownloadFile(path, url)
	//	if err != nil {
	//		return
	//	}
	//
	//	photo, err := uploadPhotoToVk(channel.VkAccessToken, channel.GroupID, path)
	//	if err != nil {
	//		return
	//	}
	//
	//	attachment = fmt.Sprintf("photo%d_%d", photo.OwnerID, photo.ID)
	//}
	//
	//var attachmentString string
	//
	//channelId := post.Chat.ID
	//postId := post.MessageID
	//
	//var forward tgbotapi.Message
	//var postText string
	//if len(post.MediaGroupID) > 0 {
	//	rdsVkAttachments.RPush(ctx, post.MediaGroupID, attachment)
	//	rdsCaption.Append(ctx, post.MediaGroupID, post.Caption)
	//	var attachments []string
	//
	//	if exists == 1 {
	//		return
	//	}
	//
	//	timer := time.NewTimer(time.Second)
	//	<-timer.C
	//	attachments, err = rdsVkAttachments.LRange(ctx, post.MediaGroupID, 0, -1).Result()
	//	if err != nil {
	//		return
	//	}
	//	fileIDS, err := rdsAttachmentFileIDS.LRange(ctx, post.MediaGroupID, 0, -1). Result()
	//	if err != nil {
	//		return
	//	}
	//	attachmentString = strings.Join(attachments, ",")
	//
	//
	//	photos := make([]interface{}, len(fileIDS))
	//	for i, fileID := range fileIDS {
	//		photo := tgbotapi.InputMediaPhoto{BaseInputMedia: tgbotapi.BaseInputMedia{
	//			Type:  "photo",
	//			Media: tgbotapi.FileID(fileID),
	//		}}
	//		if i == 0 {
	//			photo.Caption, err = rdsCaption.Get(ctx, post.MediaGroupID).Result()
	//			postText = photo.Caption
	//		}
	//		photos[i] = photo
	//	}
	//
	//	forwards, err := bot.SendMediaGroup(tgbotapi.MediaGroupConfig{
	//		ChatID: int64(creator.ID),
	//		Media:  photos,
	//	})
	//	if err != nil {
	//		return
	//	}
	//	forward = forwards[0]
	//} else {
	//	postText = post.Text
	//	attachmentString = attachment
	//	forward, err = bot.Send(tgbotapi.NewForward(int64(creator.ID), channelId, postId))
	//	if err != nil {
	//		return
	//	}
	//}
	//
	//token := generateToken()
	//msg := tgbotapi.NewMessage(int64(creator.ID), "Confirm cross-posting")
	//msg.ReplyToMessageID = forward.MessageID
	//postUrl := fmt.Sprintf("http://social.club/confirm?message_token=%s&group=%d&attachments=%s", token, channel.GroupID, attachmentString)
	//
	//msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
	//	InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{{
	//		tgbotapi.InlineKeyboardButton{
	//			Text: "Yes",
	//			URL:  &postUrl,
	//		},
	//	}},
	//}
	//
	//confirmation, err := bot.Send(msg)
	//if err != nil {
	//	return
	//}
	//
	//messageJson, err := json.Marshal(Message{
	//	ChatID:    confirmation.Chat.ID,
	//	MessageID: confirmation.MessageID,
	//	Text:      postText,
	//})
	//err = rdsMessages.Set(ctx, token, messageJson, 0).Err()
}

func handleMessage(message *tgbotapi.Message) {

}

func handleNewChannel(update tgbotapi.Update) {

}

func handleUpdate(update tgbotapi.Update) {
	switch {
	case update.ChannelPost != nil:
		handleChannelPost(update.ChannelPost)
	case update.Message != nil:
		{
			switch {
			default:
				handleMessage(update.Message)
			}
		}
	}
}

func sendFinaliseNotification(post Post, vkResponse string) (err error) {
	var resultMsg tgbotapi.MessageConfig
	if len(vkResponse) == 0 {
		// Delete all confirmation messages on success
		if _, err = bot.Send(tgbotapi.NewDeleteMessage(post.PersonalChatID, post.ConfirmationMessageID)); err != nil {
			return err
		}

		for _, id := range post.ForwardedMessagesIds {
			if _, err = bot.Send(tgbotapi.NewDeleteMessage(post.PersonalChatID, id)); err != nil {
				return err
			}
		}
		resultMsg = tgbotapi.NewMessage(post.PersonalChatID, "✅ Done!")
	} else {
		resultMsg = tgbotapi.NewMessage(post.PersonalChatID, "❌ VK API error. " + vkResponse)
	}
	_, err = bot.Send(resultMsg)
	return err
}

func mainBot() {
	for update := range setupBot() {
		go handleUpdate(update)
	}
}
