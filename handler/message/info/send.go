package info

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
)

type Info struct {
	MessageID   string
	Description string
}

// 公式情報のメッセージを更新します
//
// 新しく追加したもの(MessageIDが空のInfo)は新規送信します。
func UpdatePublicInfos(s *discordgo.Session, m *discordgo.MessageCreate) error {
	info := Info{
		MessageID: "1120581611860271227",
		Description: `
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
**🔗 公式リンク**
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
**[OpenSea]** TOKYO ANARCHY
https://opensea.io/collection/tokyoanarchy

**[Twitter]** しつぎょう✱おとうさん
https://twitter.com/shitsugyou_otou


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
**💬 あいさつ集**
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
- 朝のあいさつ「おはーきー！」


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
**🤖 botコマンド**
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
- /my-roles : 自分のロール確認

<#%s>で実行OK。

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
**🎗️ ロール**
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
<&%s> 当たり21回
<&%s> 当たり18回
<&%s> 当たり15回
<&%s> 当たり12回
<&%s> 当たり9回
<&%s> 当たり6回
<&%s> 当たり3回
<&%s> "TOKYO ANARCHY" ホルダー
<&%s> OG（配布終了）
`,
	}

	if _, err := s.ChannelMessageEditEmbed(
		internal.ChannelID().PUBLIC_INFO,
		info.MessageID,
		&discordgo.MessageEmbed{
			Description: fmt.Sprintf(
				info.Description,
				internal.ChannelID().BOT_COMMAND,
				internal.RoleID().CRAZY,
				internal.RoleID().DIAMOND,
				internal.RoleID().PLATINUM,
				internal.RoleID().GOLD,
				internal.RoleID().SILVER,
				internal.RoleID().BRONZE,
				internal.RoleID().AL,
				internal.RoleID().TOKYO_ANARCHY,
				internal.RoleID().CHAINSAW_CLUB,
			),
			Color: internal.ColorYellow,
		},
	); err != nil {
		return errors.NewError("メッセージを更新できません", err)
	}

	// 完了メッセージを送信
	if _, err := s.ChannelMessageSend(m.ChannelID, "更新が完了しました"); err != nil {
		return errors.NewError("完了メッセージを送信できません", err)
	}

	return nil
}
