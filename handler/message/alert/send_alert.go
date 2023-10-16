package alert

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
)

// 注意者リスト
var alertUsers = map[string]bool{
	"1159253963292557474": true, // 平野誠也 sh0577
	"836968004389437455":  true, // Furikake furikake8859
}

// 警告対象者がメッセージを送信したらアラートを送信します
func SendAlert(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if _, ok := alertUsers[m.Author.ID]; !ok {
		return nil
	}

	description := `
✅ 対象者
<@%s>

✅ メッセージ
%s

%s
`

	messageLink := fmt.Sprintf(
		"https://discord.com/channels/%s/%s/%s",
		m.GuildID,
		m.ChannelID,
		m.ID,
	)

	embed := &discordgo.MessageEmbed{
		Title:       "警告対象者がメッセージを送信しました",
		Description: fmt.Sprintf(description, m.Member.User.ID, m.Content, messageLink),
		Color:       internal.ColorYellow,
	}

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		return errors.NewError("パネルメッセージを送信できません", err)
	}

	return nil
}
