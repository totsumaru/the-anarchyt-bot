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
**当たり**
- 21回 <@&%s>
- 18回 <@&%s>
- 15回 <@&%s> 
- 12回 <@&%s> 
- 9回  <@&%s> 
- 6回  <@&%s> 
- 3回  <@&%s> 

**ホルダー**
- <@&%s>

**その他**
- <@&%s> OG（配布終了）
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
