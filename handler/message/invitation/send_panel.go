package invitation

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
)

// パネルを送信します
func SendPanel(s *discordgo.Session, m *discordgo.MessageCreate) error {
	btn1 := discordgo.Button{
		Label:    "招待リンクを発行",
		Style:    discordgo.PrimaryButton,
		CustomID: internal.Interaction_CustomID_Invitation,
	}

	actions := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{btn1},
	}

	description := `
----------

✅️-条件
- <@&%s>を保持

✅️-招待リンクについて
- **1人のみ** 招待可
- 発行から **7日間** 有効

✅️-招待券の取得方法
<#%s>で当たると招待券が2枚をもらえます。
その他イベントでも、貰える可能性があります。
`

	embed := &discordgo.MessageEmbed{
		Title: "招待リンクの発行",
		Description: fmt.Sprintf(
			description,
			internal.RoleID().INVITATION1,
			internal.ChannelID().GATCHA,
		),
		Color: internal.ColorDarkGray,
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
