package member

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
)

// 参加人数を取得します
func Member(s *discordgo.Session, m *discordgo.MessageCreate) error {
	guild, err := s.Guild(m.GuildID)
	if err != nil {
		return errors.NewError("ギルドを取得できません", err)
	}

	num := guild.MemberCount

	embed := &discordgo.MessageEmbed{
		Description: fmt.Sprintf("参加人数: %d人", num),
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		return errors.NewError("埋め込みメッセージを送信できません", err)
	}

	return nil
}
