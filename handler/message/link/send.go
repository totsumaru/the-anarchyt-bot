package link

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
)

type Link struct {
	Name     string
	URL      string
	ImageURL string
}

// 公式リンクのメッセージを送信します
func SendPublicURL(s *discordgo.Session, m *discordgo.MessageCreate) error {
	links := []Link{
		{
			Name:     "OpenSea",
			URL:      "https://opensea.io/collection/tokyoanarchy",
			ImageURL: "https://i.seadn.io/gae/TobG0yEHvcmk2u3UxfHaUFrWtAFbOLhFgbtT3Kir_2Fk8Y27km8MHGwXheQ4XTPnwMHFxDdZ_XGYACzBdLysFpaRlwR0sFPgc3xZ?auto=format&dpr=1&w=3840",
		},
		{
			Name:     "Twitter",
			URL:      "https://twitter.com/shitsugyou_otou",
			ImageURL: "https://cdn.discordapp.com/attachments/1103240223376293938/1116523753850028155/1500x500.png",
		},
	}

	for _, link := range links {
		embed := &discordgo.MessageEmbed{
			Title:       link.Name,
			Description: link.URL,
			Image: &discordgo.MessageEmbedImage{
				URL: link.ImageURL,
			},
			Color: internal.ColorYellow,
		}

		_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
		if err != nil {
			return errors.NewError("メッセージを送信できません", err)
		}
	}

	// コマンドメッセージを削除
	if err := s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
		return errors.NewError("コマンドメッセージを削除できません", err)
	}

	return nil
}
