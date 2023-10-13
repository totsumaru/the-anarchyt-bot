package gatcha

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
)

// 毎日の通知メッセージを送信します
func SendNotice(s *discordgo.Session) error {
	// 通知を送信します
	msg := fmt.Sprintf(
		"<@&%s>\nおはようございます！今日も<#%s>を回してみよう！！",
		internal.RoleID().GATCHA_NOTICE,
		internal.ChannelID().GATCHA,
	)

	_, err := s.ChannelMessageSend(internal.ChannelID().CHAT, msg)
	if err != nil {
		return errors.NewError("通知メッセージを送信できません", err)
	}

	return nil
}
