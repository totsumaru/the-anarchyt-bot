package free_mint

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
)

const MintURL = "https://0xspread.com/0x12d72cC8/fa40a5ea-b9bc-46bc-bf4e-d0529fc88fcb"

// 指定のロールを持っている人にはURLを返し、そうでない人にはエラーメッセージを返す
func FreeMint(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	hasSubmitted := false
	hasHolder := false

	for _, v := range i.Member.Roles {
		switch v {
		case "1166540813384306758": // 第一弾提出済み
			hasSubmitted = true
		case "1175956583977594911": // ホルダー
			hasHolder = true
		}
	}

	if hasSubmitted && hasHolder {
		_, err := s.ChannelMessageSend(i.ChannelID, MintURL)
		if err != nil {
			return errors.NewError("メッセージの送信に失敗しました", err)
		}
	} else {
		_, err := s.ChannelMessageSend(i.ChannelID, "ホルダーかつ第一弾提出済みの人のみが参加可能です")
		if err != nil {
			return errors.NewError("メッセージの送信に失敗しました", err)
		}
	}

	return nil
}
