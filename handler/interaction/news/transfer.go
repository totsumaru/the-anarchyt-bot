package news

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
)

// 送信ボタンが押された時の処理です
func Transfer(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	_, err := s.ChannelMessageSendEmbed(internal.ChannelID().NEWS, i.Message.Embeds[0])
	if err != nil {
		return errors.NewError("メッセージを送信できません", err)
	}

	// 元メッセージを変更します
	edit := &discordgo.MessageEdit{
		ID:         i.Message.ID,
		Channel:    i.ChannelID,
		Components: []discordgo.MessageComponent{},
		Embeds:     i.Message.Embeds,
	}

	_, err = s.ChannelMessageEditComplex(edit)
	if err != nil {
		return errors.NewError("元メッセージを変更できません", err)
	}

	return nil
}
