package news

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
)

// キャンセルボタンが押された時の処理です
func Cancel(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	if err := s.ChannelMessageDelete(i.ChannelID, i.Message.ID); err != nil {
		return errors.NewError("メッセージを削除できません", err)
	}

	return nil
}
