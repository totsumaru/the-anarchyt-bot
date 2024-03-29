package gatcha

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
)

// @Verifiedを持っている全員に@ガチャコインを付与します
//
// #testでのみ起動します。
//
// @はずれ,@コイン再付与済 ロールを持っている場合は、そのロールを削除します。
func AddCoinRole(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if m.ChannelID != internal.ChannelID().TEST {
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

	// ユーザーが不明の場合は処理が止まってしまうので、エラーは返さず通知のみ実施します
	for i, user := range users {
		// ユーザーにロールを付与します
		if internal.RoleID().GATCHA_COIN != "" {
			if err = s.GuildMemberRoleAdd(m.GuildID, user.User.ID, internal.RoleID().GATCHA_COIN); err != nil {
				errors.SendErrMsg(s, errors.NewError("ロールを付与できないので無視します"), user.User)
			}
		}

		// @ハズレ,@コイン2枚目付与済み を持っている人は、そのロールを削除します
		for _, role := range user.Roles {
			switch role {
			case internal.RoleID().HAZURE, internal.RoleID().COIN_2_ADDED:
				if err = s.GuildMemberRoleRemove(m.GuildID, user.User.ID, role); err != nil {
					errors.SendErrMsg(s, errors.NewError("ロールを付与できないので無視します"), user.User)
				}
			}
		}

		// 進捗メッセージを送信
		if i%50 == 0 {
			if _, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%d件の処理が完了", i)); err != nil {
				return errors.NewError("進捗メッセージを送信できません", err)
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
