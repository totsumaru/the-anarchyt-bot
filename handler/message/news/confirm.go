package news

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
	"strings"
	"time"
)

// 送信内容を確認します
func Confirm(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if m.ChannelID != internal.ChannelID().TEST {
		return nil
	}

	components := strings.Split(m.Content, "$")

	title := components[0]
	content := components[1]
	url := ""

	if len(components) > 2 {
		url = components[2]
	}

	btn1 := discordgo.Button{
		Label:    "この内容で送信します",
		Style:    discordgo.SecondaryButton,
		CustomID: internal.Interaction_CustomID_News_Send,
	}

	btn2 := discordgo.Button{
		Label:    "キャンセル",
		Style:    discordgo.PrimaryButton,
		CustomID: internal.Interaction_CustomID_News_Cancel,
	}

	actions := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{btn1, btn2},
	}

	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: content,
		Image: &discordgo.MessageEmbedImage{
			URL: url,
		},
		Color:     internal.ColorBlue,
		Timestamp: time.Now().Format("2006-01-02T15:04:05+09:00"),
	}

	data := &discordgo.MessageSend{
		Embed:      embed,
		Components: []discordgo.MessageComponent{actions},
	}

	_, err := s.ChannelMessageSendComplex(m.ChannelID, data)
	if err != nil {
		return errors.NewError("メッセージを送信できません", err)
	}

	return nil
}
