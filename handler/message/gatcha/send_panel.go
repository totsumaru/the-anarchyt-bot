package gatcha

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
	"strings"
)

// パネルを送信します
//
// 新規でパネルを送信する場合は`currentPanelURL`を空に、
// パネルを更新する場合は、現在のパネルのURLを入れてください。
func SendPanel(s *discordgo.Session, m *discordgo.MessageCreate, currentPanelURL string) error {
	btn1 := discordgo.Button{
		Label:    "ガチャを回す",
		Style:    discordgo.PrimaryButton,
		CustomID: internal.Interaction_CustomID_gatcha_Go,
	}

	btn2 := discordgo.Button{
		Label: "通知",
		Emoji: discordgo.ComponentEmoji{
			Name: "🔔",
		},
		Style:    discordgo.SecondaryButton,
		CustomID: internal.Interaction_CustomID_gatcha_Notice,
	}

	actions := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{btn1, btn2},
	}

	description := `
1日に1回、ガチャを回せます🎲
毎日チャレンジしてみてね！！

- <@&%s>で参加（毎日プレゼント）
- 当選すると<@&%s>,<@&%s>x2 GET！
- 3回当選で、AL確定GET！
- 確率は10％前後で変動するよ

---
忘れたくない人は、通知をONに！
毎朝8:00に通知されます。
`

	embed := &discordgo.MessageEmbed{
		Image: &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/attachments/1103240223376293938/1116312631721066517/title_3.png",
		},
		Title: "ロールガチャ",
		Description: fmt.Sprintf(
			description,
			internal.RoleID().GATCHA_TICKET,
			internal.RoleID().PRIZE1,
			internal.RoleID().INVITATION1,
		),
		Color: internal.ColorYellow,
	}

	if currentPanelURL == "" {
		// 新規のパネルを作成します
		messageSend := &discordgo.MessageSend{
			Components: []discordgo.MessageComponent{actions},
			Embed:      embed,
		}

		_, err := s.ChannelMessageSendComplex(m.ChannelID, messageSend)
		if err != nil {
			return errors.NewError("パネルメッセージを送信できません", err)
		}
	} else {
		// パネルを更新します

		// URL例: https://discord.com/channels/1067806759034572870/1067807967950422096/1116242093434732595
		replaced := strings.Replace(currentPanelURL, "https://discord.com/channels/", "", -1)
		ids := strings.Split(replaced, "/")

		currentPanelChannelID := ids[1]
		currentPanelMessageID := ids[2]

		edit := &discordgo.MessageEdit{
			ID:         currentPanelMessageID,
			Channel:    currentPanelChannelID,
			Components: []discordgo.MessageComponent{actions},
			Embed:      embed,
		}

		_, err := s.ChannelMessageEditComplex(edit)
		if err != nil {
			return errors.NewError("パネルを更新できません", err)
		}
	}

	// コマンドメッセージを削除
	if err := s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
		return errors.NewError("コマンドメッセージを削除できません", err)
	}

	return nil
}
