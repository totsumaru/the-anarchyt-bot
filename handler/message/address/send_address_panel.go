package address

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
)

// アドレス集計のパネルを送信します
func SendAddressPanel(s *discordgo.Session, m *discordgo.MessageCreate) error {
	btn1 := discordgo.Button{
		Label:    "登録/変更",
		Style:    discordgo.PrimaryButton,
		CustomID: internal.Interaction_CustomID_Address_Modal_Open,
	}

	btn2 := discordgo.Button{
		Label:    "チェック",
		Style:    discordgo.SecondaryButton,
		CustomID: internal.Interaction_CustomID_Address_Check,
	}

	actions := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{btn1, btn2},
	}

	description := `
ウォレットアドレスの提出フォームです。

- ウォレットアドレス
- Mint数(ロールに応じて)

を指定できます。
`

	embed := &discordgo.MessageEmbed{
		Description: description,
		Color:       internal.ColorYellow,
	}

	// 新規のパネルを作成します
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
