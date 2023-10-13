package address

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal/address"
	"github.com/techstart35/the-anarchy-bot/internal/db"
)

// Modalを開きます
func OpenModal(s *discordgo.Session, i *discordgo.InteractionCreate) error {
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

	maxMint := address.MaxMintQuantity(i.Member.Roles)

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
							Placeholder: "0x64295E222c13BC1817a174e42F35fC0ccbf3c5AD",
							Required:    true,
							MinLength:   42,
							MaxLength:   42,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "quantity",
							Label:       fmt.Sprintf("ミント数(上限: %d)", maxMint),
							Style:       discordgo.TextInputShort,
							Value:       quantityValue,
							Placeholder: strconv.Itoa(maxMint),
							Required:    true,
							MinLength:   1,
							MaxLength:   1,
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
