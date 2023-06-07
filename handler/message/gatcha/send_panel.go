package gatcha

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
)

// パネルを送信します
func SendPanel(s *discordgo.Session, m *discordgo.MessageCreate) error {
	btn1 := discordgo.Button{
		Label:    "ガチャを回す（1日1回）",
		Style:    discordgo.PrimaryButton,
		CustomID: internal.Interaction_CustomID_gatcha_Go,
	}

	actions := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{btn1},
	}

	description := `
1日に1回、ガチャを回せます🎲
毎日チャレンジしてみてね！！

- <@&%s>で参加（毎日1枚プレゼント）
- 当選すると<@&%s>ロールがもらえるよ
- 3回当選で、AL確定GET！
- 確率は状況に応じて変動するよ
`

	embed := &discordgo.MessageEmbed{
		Image: &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/attachments/1103240223376293938/1115571819362136064/AdobeStock_483707441.jpeg",
		},
		Title:       "ロールガチャ",
		Description: fmt.Sprintf(description, internal.RoleID().TICKET, internal.RoleID().PRIZE1),
		Color:       internal.ColorYellow,
	}

	data := &discordgo.MessageSend{
		Components: []discordgo.MessageComponent{actions},
		Embed:      embed,
	}

	_, err := s.ChannelMessageSendComplex(m.ChannelID, data)
	if err != nil {
		return errors.NewError("パネルメッセージを送信できません", err)
	}

	// コマンドメッセージを削除
	if err = s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
		return errors.NewError("コマンドメッセージを削除できません", err)
	}

	return nil
}
