package address

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
	"github.com/techstart35/the-anarchy-bot/internal/address"
	"github.com/techstart35/the-anarchy-bot/internal/db"
	"gorm.io/gorm"
)

// ModalからSubmitされた内容を処理します
func UpsertFromModal(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// 入力値のバリデーションを行います
	addr := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	quantityStr := i.ModalSubmitData().Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	if !IsValidETHAddress(addr) {
		description := "【ERROR】アドレスの形式が一致していません。"
		if err := reply(s, i, description, internal.ColorRed); err != nil {
			return errors.NewError("エラーメッセージを送信できません")
		}
		return nil
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		description := "【ERROR】ミント数が半角数字では無いため、登録できませんでした。"
		if err = reply(s, i, description, internal.ColorRed); err != nil {
			return errors.NewError("エラーメッセージを送信できません")
		}
		return nil
	}

	// quantityが上限を超えていないか確認します
	if address.MaxMintQuantity(i.Member.Roles) < quantity {
		description := "【ERROR】Mint上限数を超えているため、登録できませんでした。"
		if err = reply(s, i, description, internal.ColorRed); err != nil {
			return errors.NewError("エラーメッセージを送信できません")
		}
		return nil
	}

	// DBに保存します
	err = db.DB.Transaction(func(tx *gorm.DB) error {
		// quantityが0の場合は削除します
		if quantity == 0 {
			// 削除します
			if err = db.Remove(tx, i.Member.User.ID); err != nil {
				return errors.NewError("削除に失敗しました", err)
			}
		} else {
			// 更新します
			if err = db.Upsert(tx, i.Member.User.ID, addr, quantity); err != nil {
				return errors.NewError("Upsertに失敗しました", err)
			}
		}

		return nil
	})
	if err != nil {
		return errors.NewError("DBの保存に失敗しました", err)
	}

	// 返信を送信します
	description := `
✅登録完了

アドレス
- %s
ミント数
- %d
`

	desc := fmt.Sprintf(description, addr, quantity)
	if err = reply(s, i, desc, internal.ColorBlue); err != nil {
		return errors.NewError("返信を送信します", err)
	}

	return nil
}

// 回答者に対する送信です
func reply(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	description string,
	color int,
) error {
	embed := &discordgo.MessageEmbed{
		Description: description,
		Color:       color,
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

// ウォレットアドレスを検証します
func IsValidETHAddress(address string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(address)
}
