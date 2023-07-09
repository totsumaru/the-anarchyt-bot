package gatcha

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
)

// @Verifiedを持っている全員に@ボーナスコインを付与します
//
// #testでのみ起動します。
func AddBonusCoinRole(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if m.ChannelID != internal.ChannelID().TEST {
		return nil
	}

	// 進捗メッセージを送信
	if _, err := s.ChannelMessageSend(m.ChannelID, "@Verifiedを持っている全員にボーナスコインを付与します"); err != nil {
		return errors.NewError("進捗メッセージを送信できません", err)
	}

	users, err := getUsersWithRole(s, m.GuildID, internal.RoleID().VERIFIED)
	if err != nil {
		return errors.NewError("特定のロールを持つユーザーを取得できません", err)
	}

	for _, user := range users {
		// ユーザーにロールを付与します
		if internal.RoleID().BONUS_COIN != "" {
			if err = s.GuildMemberRoleAdd(m.GuildID, user.User.ID, internal.RoleID().BONUS_COIN); err != nil {
				return errors.NewError("ロールを付与できません", err)
			}
		}
	}

	// 完了メッセージを送信
	if _, err = s.ChannelMessageSend(m.ChannelID, "ボーナスコインの配布が完了しました"); err != nil {
		return errors.NewError("完了メッセージを送信できません", err)
	}

	return nil
}
