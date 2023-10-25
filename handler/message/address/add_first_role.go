package address

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
)

// 第一弾で登録した人にロールを付与します
func AddFirstWalletSubmittedRole(s *discordgo.Session, m *discordgo.MessageCreate) error {
	// Totsumaru以外は実行できません
	if m.Author.ID != internal.UserID().TOTSUMARU {
		return nil
	}

	if _, err := s.ChannelMessageSend(m.ChannelID, "開始します"); err != nil {
		return errors.NewError("メッセージを送信できません")
	}

	// 指定されたギルドのメンバーを取得します。
	// 1000は1回のAPIリクエストで取得するメンバーの最大数です。必要に応じて調整してください。
	members, err := s.GuildMembers(m.GuildID, "", 1000)
	if err != nil {
		return err
	}

	// 付与済みカウント
	count := 0

	// 各メンバーを調査して、ロールAを持っているかどうかを確認します。
	for _, member := range members {
		// メンバーがロールAを持っているかどうかを確認します。
		hasSubmittedRole := false
		for _, roleID := range member.Roles {
			if roleID == internal.RoleID().SUBMITTED {
				hasSubmittedRole = true
				break
			}
		}

		// メンバーがロールAを持っている場合、ロールBを付与します。
		if hasSubmittedRole {
			if err = s.GuildMemberRoleAdd(
				m.GuildID,
				member.User.ID,
				internal.RoleID().FIRST_SUBMITTED,
			); err != nil {
				return err
			}
			count++
		}

		if _, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%d人にロールを付与しました。", count)); err != nil {
			return errors.NewError("メッセージを送信できません")
		}
	}

	return nil
}
