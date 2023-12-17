package roles

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
)

// `/my-roles`コマンドが実行された時の処理です。
//
// 自分のロールを出力します
func GetRoles(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	roleNames := make([]string, 0)

	roleIDs, err := RoleListInOrder(s, i.GuildID, i.Member.User.ID)
	if err != nil {
		return errors.NewError("ロール一覧を取得できません", err)
	}

	guildRoles, err := s.GuildRoles(i.GuildID)
	if err != nil {
		return errors.NewError("全てのロールを取得できません", err)
	}

	var thumbnailURL string
	var color = internal.ColorBlack

	for _, roleID := range roleIDs {
		for _, role := range guildRoles {
			if role.ID == roleID {
				roleNames = append(roleNames, role.Name)
				break
			}
		}

		switch roleID {
		case internal.RoleID().BRONZE:
			thumbnailURL = "https://media.discordapp.net/attachments/1103240223376293938/1128924752396963900/bronze.png?width=256&height=256"
			color = 0xd87b5f
		case internal.RoleID().SILVER:
			thumbnailURL = "https://media.discordapp.net/attachments/1103240223376293938/1133570074867929188/silver.png?width=1068&height=1068"
			color = 0xacb8b8
		case internal.RoleID().GOLD:
			thumbnailURL = "https://cdn.discordapp.com/attachments/1103240223376293938/1133570074461077504/gold.png"
			color = 0xfdbd52
		case internal.RoleID().PLATINUM:
			thumbnailURL = "https://cdn.discordapp.com/attachments/1103240223376293938/1140805359162884217/platinum.png"
			color = 0x7e76f0
		case internal.RoleID().DIAMOND:
			thumbnailURL = "https://cdn.discordapp.com/attachments/1103240223376293938/1140805358839926824/diamond.png"
			color = 0x429df5
		case internal.RoleID().CRAZY:
			thumbnailURL = "https://cdn.discordapp.com/attachments/1103240223376293938/1140805358294675589/crazy.png"
			color = 0xe41b67
		}
	}

	description := `
**YOUR ROLES**

%s
`

	userName := i.Member.User.Username
	if i.Member.User.GlobalName != "" {
		userName = i.Member.User.GlobalName
	}
	if i.Member.Nick != "" {
		userName = i.Member.Nick
	}

	point := 0
	for _, role := range i.Member.Roles {
		switch role {
		case internal.RoleID().PRIZE1,
			internal.RoleID().PRIZE2,
			internal.RoleID().PRIZE3,
			internal.RoleID().PRIZE4,
			internal.RoleID().PRIZE5,
			internal.RoleID().PRIZE6,
			internal.RoleID().PRIZE7,
			internal.RoleID().PRIZE8,
			internal.RoleID().PRIZE9,
			internal.RoleID().PRIZE10,
			internal.RoleID().PRIZE11,
			internal.RoleID().PRIZE12,
			internal.RoleID().PRIZE13,
			internal.RoleID().PRIZE14:
			point += 1
		case internal.RoleID().BRONZE:
			point += 6
		case internal.RoleID().SILVER:
			point += 9
		case internal.RoleID().GOLD:
			point += 12
		case internal.RoleID().PLATINUM:
			point += 15
		case internal.RoleID().DIAMOND:
			point += 18
		case internal.RoleID().CRAZY:
			point += 21
		case internal.RoleID().FUCKIN:
			point += 36
		}
	}

	text := "カード画像はコピーして、Xでポストしてみよう！"
	if point < 6 {
		text = "IDカードはブロンズ以上で表示されるよ"
	}

	embed := &discordgo.MessageEmbed{
		Description: fmt.Sprintf(
			description,
			"- "+strings.Join(roleNames, "\n- "),
		),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: thumbnailURL,
		},
		Color: color,
		Footer: &discordgo.MessageEmbedFooter{
			Text: text,
		},
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

	// 画像を送信
	if point >= 6 {
		imageURL := fmt.Sprintf(
			"https://the-anarchy-gatcha-image.vercel.app/api/card?username=%s&avatar=%s&point=%d",
			url.QueryEscape(userName),
			url.QueryEscape(i.Member.User.AvatarURL("")),
			point,
		)

		// 画像URLからデータを取得
		httpRes, err := http.Get(imageURL)
		if err != nil {
			return errors.NewError("画像を取得できません", err)
		}
		defer httpRes.Body.Close()

		imageData, err := io.ReadAll(httpRes.Body)
		if err != nil {
			return errors.NewError("画像データの読み取り中にエラーが発生しました", err)
		}

		imageReader := bytes.NewReader(imageData)

		// 画像データをDiscordに送信
		if _, err = s.ChannelFileSend(i.Interaction.ChannelID, "image.jpg", imageReader); err != nil {
			return errors.NewError("画像を送信できません", err)
		}
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
