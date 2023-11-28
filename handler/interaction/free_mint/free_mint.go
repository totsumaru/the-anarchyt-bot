package free_mint

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
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

	description := ""

	if hasSubmitted && hasHolder {
		description = fmt.Sprintf("FreeMintのURLはこちらです\nキーワード: Anarchy\n%s", MintURL)
	} else {
		description = "ホルダーかつ第一弾提出済みの人のみが参加可能です"
	}

	embed := &discordgo.MessageEmbed{
		Description: description,
		Color:       internal.ColorBlue,
	}

	resp := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	}

	if err := s.InteractionRespond(i.Interaction, resp); err != nil {
		return errors.NewError("レスポンスを送信できません", err)
	}

	return nil
}
