package address

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
	"github.com/techstart35/the-anarchy-bot/internal/address"
	"github.com/techstart35/the-anarchy-bot/internal/db"
)

// 登録したアドレスを確認します
func Check(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	maxMint := address.MaxMintQuantity(i.Member.Roles)

	// maxMintが0の人にはエラーメッセージを送信します
	if maxMint == 0 {
		embed := &discordgo.MessageEmbed{
			Description: "AL対象のロールを保持していません。",
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

	// IDでウォレットを取得します
	wallet, err := db.FindByID(i.Member.User.ID)
	if err != nil {
		return errors.NewError("IDでウォレットを取得できません", err)
	}

	btn1 := discordgo.Button{
		Label:    "変更/登録",
		Style:    discordgo.PrimaryButton,
		CustomID: internal.Interaction_CustomID_Address_Modal_Open,
	}

	actions := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{btn1},
	}

	description := "ウォレットアドレス\n" +
		"```%s```\n" +
		"ミント数\n" +
		"```%d mint```\n" +
		"※上限は %d mint"
	addr := wallet.Address
	if addr == "" {
		addr = "登録なし"
	}
	quantity := wallet.Quantity

	embed := &discordgo.MessageEmbed{
		Description: fmt.Sprintf(description, addr, quantity, maxMint),
		Color:       internal.ColorYellow,
	}

	resp := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{actions},
			Embeds:     []*discordgo.MessageEmbed{embed},
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	}

	if err = s.InteractionRespond(i.Interaction, resp); err != nil {
		return errors.NewError("レスポンスを送信できません", err)
	}

	return nil
}
