package sneek_peek

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
)

// #チラ見せに投稿されたら#チャットに通知します
func Notice(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if m.ChannelID != internal.ChannelID().SNEAK_PEEK {
		return nil
	}

	description := `
<#%s>に投稿がありました。チェックしましょう！
`

	embed := &discordgo.MessageEmbed{
		Description: fmt.Sprintf(description, internal.ChannelID().SNEAK_PEEK),
		Color:       internal.ColorYellow,
	}

	_, err := s.ChannelMessageSendEmbed(internal.ChannelID().CHAT, embed)
	if err != nil {
		return errors.NewError("パネルメッセージを送信できません", err)
	}

	return nil
}
