package free_mint

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
)

// パネルメッセージを送信します
func SendFreeMintPanel(s *discordgo.Session, m *discordgo.MessageCreate) error {
	btn1 := discordgo.Button{
		Label:    "リンクを表示",
		Style:    discordgo.PrimaryButton,
		CustomID: internal.Interaction_CustomID_FreeMint,
	}

	actions := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{btn1},
	}

	description := `
- 第一弾提出済み
- ホルダー

の人のみが参加可能です。
`

	embed := &discordgo.MessageEmbed{
		Title:       "FreeMintリンクの表示",
		Description: description,
		Color:       internal.ColorYellow,
	}

	messageSend := &discordgo.MessageSend{
		Components: []discordgo.MessageComponent{actions},
		Embed:      embed,
	}

	_, err := s.ChannelMessageSendComplex(m.ChannelID, messageSend)
	if err != nil {
		return errors.NewError("パネルメッセージを送信できません", err)
	}

	// コマンドメッセージを削除
	if err = s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
		return errors.NewError("コマンドメッセージを削除できません", err)
	}

	return nil
}
