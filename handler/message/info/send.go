package info

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
)

type Link struct {
	Name     string
	Content  string
	ImageURL string
}

// 公式情報のメッセージを送信します
func SendPublicInfo(s *discordgo.Session, m *discordgo.MessageCreate) error {
	for _, embed := range infoEmbed() {
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
		if err != nil {
			return errors.NewError("メッセージを送信できません", err)
		}
	}

	// コマンドメッセージを削除
	if err := s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
		return errors.NewError("コマンドメッセージを削除できません", err)
	}

	return nil
}

// 公式情報のメッセージを更新します
func UpdatePublicInfo(s *discordgo.Session, m *discordgo.MessageCreate) error {
	const (
		InfoMessageChannelID   = "1116472032738152588"
		InfoMessageID_Link     = "1116525464752754798"
		InfoMessageID_Greeting = "1120581611860271227"
	)

	messageIDs := []string{
		InfoMessageID_Link,
		InfoMessageID_Greeting,
	}

	for i, embed := range infoEmbed() {
		_, err := s.ChannelMessageEditEmbed(
			InfoMessageChannelID,
			messageIDs[i],
			embed,
		)
		if err != nil {
			return errors.NewError("メッセージを更新できません", err)
		}
	}

	// 完了メッセージを送信
	if _, err := s.ChannelMessageSend(m.ChannelID, "更新が完了しました"); err != nil {
		return errors.NewError("完了メッセージを送信できません", err)
	}

	return nil
}

// 公式情報の送信内容です
func infoEmbed() []*discordgo.MessageEmbed {
	description1 := `
**🔗｜公式リンク**

**[OpenSea]** TOKYO ANARCHY
https://opensea.io/collection/tokyoanarchy

**[Twitter]** しつぎょう✱おとうさん
https://twitter.com/shitsugyou_otou
`

	description2 := `
**💬｜あいさつ集**

- 朝のあいさつ「おはーきー！」
`
	return []*discordgo.MessageEmbed{
		{
			Description: description1,
			Color:       internal.ColorYellow,
		},
		{
			Description: description2,
			Color:       internal.ColorYellow,
		},
	}
}
