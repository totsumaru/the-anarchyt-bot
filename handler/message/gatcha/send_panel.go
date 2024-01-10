package gatcha

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
)

// パネルを送信します
//
// 新規でパネルを送信する場合は`currentPanelURL`を空に、
// パネルを更新する場合は、現在のパネルのURLを入れてください。
func SendPanel(s *discordgo.Session, m *discordgo.MessageCreate, currentPanelURL string) error {
	btn1 := discordgo.Button{
		Label:    "ガチャを回す！",
		Style:    discordgo.PrimaryButton,
		CustomID: internal.Interaction_CustomID_gatcha_Go,
		Emoji:    discordgo.ComponentEmoji{Name: "▶️"},
	}

	btn2 := discordgo.Button{
		Label: "通知 ON/OFF",
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
			URL: getRandImageURL(),
		},
		Title: "ロールガチャ",
		Description: fmt.Sprintf(
			description,
			internal.RoleID().GATCHA_COIN,
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

// ランダムでTOP画像が表示されるようにします
func getRandImageURL() string {
	urls := []string{
		// 1. 初期の赤い画像
		"https://cdn.discordapp.com/attachments/1103240223376293938/1116312631721066517/title_3.png",
		// 2. 青い画像
		"https://cdn.discordapp.com/attachments/1067807967950422096/1141186361571954778/title_02.jpg",
		// 3. ピンクの画像
		"https://cdn.discordapp.com/attachments/1067807967950422096/1141186362599542916/title_04.jpg",
		// 4. モノクロ画像
		"https://cdn.discordapp.com/attachments/1067807967950422096/1141186362209476608/title_03.jpg",
	}
	rand.Seed(time.Now().UnixNano())
	return urls[rand.Intn(len(urls))]
}
