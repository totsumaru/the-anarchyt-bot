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
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, infoEmbed())
	if err != nil {
		return errors.NewError("メッセージを送信できません", err)
	}

	// コマンドメッセージを削除
	if err = s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
		return errors.NewError("コマンドメッセージを削除できません", err)
	}

	return nil
}

// 公式情報のメッセージを更新します
func UpdatePublicInfo(s *discordgo.Session, m *discordgo.MessageCreate) error {
	const (
		InfoMessageChannelID = "1116472032738152588"
		InfoMessageID        = "1116525464752754798"
	)

	_, err := s.ChannelMessageEditEmbed(InfoMessageChannelID, InfoMessageID, infoEmbed())
	if err != nil {
		return errors.NewError("メッセージを更新できません", err)
	}

	// 完了メッセージを送信
	if _, err = s.ChannelMessageSend(m.ChannelID, "更新が完了しました"); err != nil {
		return errors.NewError("完了メッセージを送信できません", err)
	}

	return nil
}

// 公式情報の送信内容です
func infoEmbed() *discordgo.MessageEmbed {
	description := `
**🔗｜公式リンク**

**OpenSea** TOKYO ANARCHY
https://opensea.io/collection/tokyoanarchy

**Twitter** しつぎょう✱おとうさん
https://twitter.com/shitsugyou_otou
`
	embed := &discordgo.MessageEmbed{
		Description: description,
		Image: &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/attachments/1103240223376293938/1116523753850028155/1500x500.png",
		},
		Color: internal.ColorYellow,
	}

	return embed
}
