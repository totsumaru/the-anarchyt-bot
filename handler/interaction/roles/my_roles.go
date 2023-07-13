package roles

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
	"sort"
	"strings"
)

// `/my-roles`コマンドが実行された時の処理です。
//
// 自分のロールを出力します
func GetRoles(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	formatRoles := make([]string, 0)

	roleIDs, err := RoleListInOrder(s, i.GuildID, i.Member.User.ID)
	if err != nil {
		return errors.NewError("ロール一覧を取得できません", err)
	}

	var thumbnailURL string

	for _, roleID := range roleIDs {
		formatRoles = append(formatRoles, fmt.Sprintf("<@&%s>", roleID))

		switch roleID {
		case internal.RoleID().BRONZE:
			thumbnailURL = "https://media.discordapp.net/attachments/1103240223376293938/1128924752396963900/bronze.png?width=256&height=256"
		}
	}

	description := `
**YOUR ROLES**

%s
`

	embed := &discordgo.MessageEmbed{
		Description: fmt.Sprintf(
			description,
			strings.Join(formatRoles, "\n"),
		),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: thumbnailURL,
		},
		Color: internal.ColorBlue,
	}

	resp := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	}

	if err = s.InteractionRespond(i.Interaction, resp); err != nil {
		return errors.NewError("レスポンスを送信できません", err)
	}

	return nil
}

// RoleListInOrder prints the roles of a user in the order they appear on the server
func RoleListInOrder(s *discordgo.Session, guildID, userID string) ([]string, error) {
	guild, err := s.State.Guild(guildID)
	if err != nil {
		return nil, err
	}

	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		return nil, err
	}

	// Create a map of role IDs to roles
	roleMap := make(map[string]*discordgo.Role)
	for _, role := range guild.Roles {
		roleMap[role.ID] = role
	}

	// Get the member's roles
	memberRoles := make([]*discordgo.Role, len(member.Roles))
	for i, roleID := range member.Roles {
		memberRoles[i] = roleMap[roleID]
	}

	// Sort the member's roles according to their position on the server
	sort.SliceStable(memberRoles, func(i, j int) bool {
		// Compare the positions of the roles on the server
		// Higher position should come first
		return memberRoles[i].Position > memberRoles[j].Position
	})

	// Convert the member's roles from IDs to names
	roleList := make([]string, len(member.Roles))
	for i, role := range memberRoles {
		roleList[i] = role.ID
	}

	return roleList, nil
}
