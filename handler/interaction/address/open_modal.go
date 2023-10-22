package address

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
	"github.com/techstart35/the-anarchy-bot/internal/address"
	"github.com/techstart35/the-anarchy-bot/internal/db"
)

// Modalを開きます
func OpenModal(s *discordgo.Session, i *discordgo.InteractionCreate) error {
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

	// 既存の登録情報を取得します
	// IDでウォレットを取得します
	wallet, err := db.FindByID(i.Member.User.ID)
	if err != nil {
		return errors.NewError("IDでウォレットを取得できません", err)
	}

	quantityValue := ""
	if wallet.Quantity != 0 {
		quantityValue = strconv.Itoa(wallet.Quantity)
	}

	// Modalを表示します
	if err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "modal",
			Title:    "Wallet Address",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "address",
							Label:       "ウォレットアドレス",
							Style:       discordgo.TextInputShort,
							Value:       wallet.Address,
							Placeholder: "0x0000....000",
							Required:    true,
							MinLength:   42,
							MaxLength:   42,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  "quantity",
							Label:     fmt.Sprintf("ミント数(上限: %d)", maxMint),
							Style:     discordgo.TextInputShort,
							Value:     quantityValue,
							Required:  true,
							MinLength: 1,
							MaxLength: 1,
						},
					},
				},
			},
		},
	}); err != nil {
		return errors.NewError("ModalをOpenできません", err)
	}

	return nil
}
