package sneek_peek

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
)

// #チラ見せに投稿されたものを#チャットに転送します
func Transfer(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if m.ChannelID != internal.ChannelID().SNEAK_PEEK {
		return nil
	}

	embed := &discordgo.MessageEmbed{
		Image: &discordgo.MessageEmbedImage{
			URL: m.Attachments[0].URL,
		},
		Description: m.Content,
		Color:       internal.ColorYellow,
	}

	contentTmpl := `
<#%s>に投稿がありました。チェックしましょう！
`

	data := &discordgo.MessageSend{
		Content: fmt.Sprintf(contentTmpl, internal.ChannelID().SNEAK_PEEK),
		Embed:   embed,
	}

	_, err := s.ChannelMessageSendComplex(m.ChannelID, data)
	if err != nil {
		return errors.NewError("パネルメッセージを送信できません", err)
	}

	return nil
}
