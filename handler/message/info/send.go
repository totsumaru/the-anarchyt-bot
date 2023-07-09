package info

import (
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
	greetingInfo := Info{
		MessageID: "1120581611860271227",
		Description: `
**🔗｜公式リンク**

**[OpenSea]** TOKYO ANARCHY
https://opensea.io/collection/tokyoanarchy

**[Twitter]** しつぎょう✱おとうさん
https://twitter.com/shitsugyou_otou

━━━━━━━━━━━━━━━━━━

**💬｜あいさつ集**

- 朝のあいさつ「おはーきー！」

━━━━━━━━━━━━━━━━━━

**🤖｜botコマンド**

- /my-roles : 自分のロール確認

<#%s>で実行OK。

━━━━━━━━━━━━━━━━━━
`,
	}

	infos := []Info{greetingInfo}

	for _, info := range infos {
		if info.MessageID == "" {
			if _, err := s.ChannelMessageSendEmbed(
				internal.ChannelID().PUBLIC_INFO,
				&discordgo.MessageEmbed{
					Description: info.Description,
					Color:       internal.ColorYellow,
				},
			); err != nil {
				return errors.NewError("メッセージを送信できません", err)
			}
		} else {
			if _, err := s.ChannelMessageEditEmbed(
				internal.ChannelID().PUBLIC_INFO,
				info.MessageID,
				&discordgo.MessageEmbed{
					Description: info.Description,
					Color:       internal.ColorYellow,
				},
			); err != nil {
				return errors.NewError("メッセージを更新できません", err)
			}
		}
	}

	// 完了メッセージを送信
	if _, err := s.ChannelMessageSend(m.ChannelID, "更新が完了しました"); err != nil {
		return errors.NewError("完了メッセージを送信できません", err)
	}

	return nil
}
