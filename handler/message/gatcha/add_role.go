package gatcha

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
)

// @Verifiedを持っている全員に@ガチャ券を付与します
//
// #logsでのみ起動します。
//
// @はずれロールを持っている場合は、そのロールを削除します。
func AddRole(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if m.ChannelID != internal.ChannelID().LOGS {
		return nil
	}

	// 進捗メッセージを送信
	if _, err := s.ChannelMessageSend(m.ChannelID, "Process is running..."); err != nil {
		return errors.NewError("進捗メッセージを送信できません", err)
	}

	users, err := getUsersWithRole(s, m.GuildID, internal.RoleID().VERIFIED)
	if err != nil {
		return errors.NewError("特定のロールを持つユーザーを取得できません", err)
	}

	for _, user := range users {
		// ユーザーにロールを付与します
		if internal.RoleID().GATCHA_TICKET != "" {
			if err = s.GuildMemberRoleAdd(m.GuildID, user.User.ID, internal.RoleID().GATCHA_TICKET); err != nil {
				return errors.NewError("ロールを付与できません", err)
			}
		}

		// ハズレロールを持っている人は、削除します
		for _, role := range user.Roles {
			if role == internal.RoleID().HAZURE {
				if err = s.GuildMemberRoleRemove(m.GuildID, user.User.ID, internal.RoleID().HAZURE); err != nil {
					return errors.NewError("ハズレロールを削除できません", err)
				}
			}
		}
	}

	// 完了メッセージを送信
	if _, err = s.ChannelMessageSend(m.ChannelID, "Finished!"); err != nil {
		return errors.NewError("完了メッセージを送信できません", err)
	}

	return nil
}

// 特定のロールを持つユーザーを全て取得します
func getUsersWithRole(s *discordgo.Session, guildID, roleID string) ([]*discordgo.Member, error) {
	members := make([]*discordgo.Member, 0)

	var lastID string

	for {
		guildMembers, err := s.GuildMembers(guildID, lastID, 1000)
		if err != nil {
			return nil, err
		}

		for _, member := range guildMembers {
			for _, memberRoleID := range member.Roles {
				if memberRoleID == roleID {
					members = append(members, member)
					break
				}
			}
		}

		if len(guildMembers) < 1000 {
			break
		}

		lastID = guildMembers[len(guildMembers)-1].User.ID
	}

	return members, nil
}
