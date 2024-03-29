package address

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
	"github.com/techstart35/the-anarchy-bot/internal/db"
)

// 提出した人全員に提出済みロールを付与します
func AddSubmittedRole(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if m.ChannelID != internal.ChannelID().TEST {
		return nil
	}

	wallets, err := db.FindAll()
	if err != nil {
		return errors.NewError("全ての情報を取得できません", err)
	}

	for i, wallet := range wallets {
		if i%50 == 0 {
			if _, err = s.ChannelMessageSend(
				m.ChannelID,
				fmt.Sprintf("%d件が完了", i),
			); err != nil {
				return errors.NewError("途中経過メッセージを送信できません", err)
			}
		}
		if err = s.GuildMemberRoleAdd(m.GuildID, wallet.ID, internal.RoleID().SUBMITTED); err != nil {
			return errors.NewError("提出した人全員にロールを付与できません", err)
		}
	}

	return nil
}
