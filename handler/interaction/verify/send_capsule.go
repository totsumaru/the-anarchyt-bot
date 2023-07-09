package verify

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
)

// 認証をします
func Verify(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// 認証済みの場合は取得できません
	if hasRole(i.Member, internal.RoleID().VERIFIED) {
		if err := sendAlreadyVerifiedMessage(s, i); err != nil {
			return errors.NewError("認証済みエラーメッセージを送信できません", err)
		}
		return nil
	}

	// verifyロールとチケットロールを付与します
	if err := addRoles(s, i, []string{internal.RoleID().VERIFIED, internal.RoleID().GATCHA_COIN}); err != nil {
		return errors.NewError("ロールを付与できません", err)
	}

	description := `
<@&%s>を取得しました。
まずはチャットで挨拶をしてみましょう！

<@&%s>を1枚獲得しました。
<#%s>でガチャが回せます！
`

	embed := &discordgo.MessageEmbed{
		Description: fmt.Sprintf(
			description,
			internal.RoleID().VERIFIED,
			internal.RoleID().GATCHA_COIN,
			internal.ChannelID().GATCHA,
		),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://cdn.discordapp.com/attachments/1103240223376293938/1115225765542363166/anarchy.jpg",
		},
		Color: internal.ColorBlue,
	}

	resp := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	}

	if err := s.InteractionRespond(i.Interaction, resp); err != nil {
		return errors.NewError("レスポンスを送信できません", err)
	}

	return nil
}

// ロールを付与します
func addRoles(s *discordgo.Session, i *discordgo.InteractionCreate, roles []string) error {
	for _, role := range roles {
		if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, role); err != nil {
			return errors.NewError("ロールを付与できません", err)
		}
	}

	return nil
}

// ロールを保持していることを確認します
func hasRole(member *discordgo.Member, roleID string) bool {
	for _, role := range member.Roles {
		if role == roleID {
			return true
		}
	}

	return false
}

// 認証済みメッセージを送信します
func sendAlreadyVerifiedMessage(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	embed := &discordgo.MessageEmbed{
		Description: "認証済みです",
		Color:       internal.ColorRed,
	}

	resp := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	}

	if err := s.InteractionRespond(i.Interaction, resp); err != nil {
		return errors.NewError("レスポンスを送信できません", err)
	}

	return nil
}
