package invitation

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
)

// 招待を1つも持っていない人に@招待券ロールを付与します
func AddInvitationRole(s *discordgo.Session, m *discordgo.MessageCreate) error {
	var after string

	for {
		members, err := s.GuildMembers(m.GuildID, after, 1000)
		if err != nil {
			return errors.NewError("メンバーを取得できません", err)
		}

		for _, member := range members {
			hasInvitation := false
			for _, roleID := range member.Roles {
				if roleID == internal.RoleID().INVITATION1 || roleID == internal.RoleID().INVITATION2 {
					hasInvitation = true
				}
			}
			if !hasInvitation {
				if err = s.GuildMemberRoleAdd(m.GuildID, member.User.ID, internal.RoleID().INVITATION1); err != nil {
					return errors.NewError("ロールを付与できません", err)
				}
			}
		}

		if len(members) < 1000 {
			break
		}

		after = members[len(members)-1].User.ID
	}

	return nil
}
