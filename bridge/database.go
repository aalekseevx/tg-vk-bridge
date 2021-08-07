package bridge

import (
	"crypto/rand"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
	"time"
)

var db *gorm.DB

func generateToken() (string, error) {
	b := make([]byte, 16)
	var err error
	attempts := 10
	n := 0
	for i := 0; i < attempts; i++ {
		n, err = rand.Read(b)
		if n > 0 {
			break
		}
		time.Sleep(time.Second)
	}
	return fmt.Sprintf("%x", b), err
}

func addChannel(ChatID int64, creatorID int) error {
	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&Channel{
		ChatID: ChatID,
		CreatorID: creatorID,
	}).Error
}

func addCreator(ChatID int) error {
	token, err := generateToken()
	if err != nil {
		return err
	}
	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&Creator{
		ChatID: ChatID,
		Token: token,
	}).Error
}

func getCreatorIDByToken(token string) (int, error) {
	var creator Creator
	err := db.Select("id").Where("token", token).Take(&creator).Error
	return creator.ChatID, err
}

func getChannel(ChatID int64) (channel Channel, err error) {
	err = db.Take(&channel, ChatID).Error
	return
}

func getChannelsByCreatorToken(creatorToken string) (channels []Channel, err error) {
	creatorId, err := getCreatorIDByToken(creatorToken)
	if err == nil {
		err = db.Select("chat_id", "name").Where("creator_id", creatorId).Find(&channels).Error
	}
	return
}

func isOwnerToken(creatorToken string, channelId int64) (bool, error) {
	channel, err := getChannel(channelId)
	if err != nil {
		return false, err
	}
	creatorIdFromDb, err := getCreatorIDByToken(creatorToken)
	if err != nil {
		return false, err
	}
	return channel.CreatorID == creatorIdFromDb, nil
}

func matchChannel(channel Channel) error {
	return db.Model(&Channel{ChatID: channel.ChatID}).Updates(channel).Error
}

func initDatabase() {
	var err error
	db, err = gorm.Open(sqlite.Open(os.Getenv("DATABASE_FILENAME")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&Creator{}, &Channel{})
	if err != nil {
		panic("failed to migrate schema")
	}
}
