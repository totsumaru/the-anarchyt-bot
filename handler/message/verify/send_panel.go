package verify

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
)

// パネルを送信します
func SendPanel(s *discordgo.Session, m *discordgo.MessageCreate) error {
	btn1 := discordgo.Button{
		Label:    "入場",
		Style:    discordgo.PrimaryButton,
		CustomID: internal.Interaction_CustomID_Verify,
	}

	actions := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{btn1},
	}

	description := `
ようこそ、**THE ANARCHY**の世界へ！

ボタンを押して、コミュニティに参加しましょう。
`

	embed := &discordgo.MessageEmbed{
		Title:       "Welcome",
		Description: description,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://cdn.discordapp.com/attachments/1103240223376293938/1115225765542363166/anarchy.jpg",
		},
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
