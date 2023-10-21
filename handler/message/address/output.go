package address

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
	"github.com/techstart35/the-anarchy-bot/internal/db"
)

// 現在の集計結果を出力します
func OutputAddresses(s *discordgo.Session, m *discordgo.MessageCreate) error {
	// 出力できるチャンネルのみで実行できます
	if !(m.ChannelID == internal.ChannelID().TEAM ||
		m.ChannelID == internal.ChannelID().TEST) {
		return nil
	}

	// DBから取得します
	allWallets, err := db.FindAll()
	if err != nil {
		return errors.NewError("全てのレコードを取得できません", err)
	}

	var allQuantity int
	for _, wallet := range allWallets {
		allQuantity += wallet.Quantity
	}

	description := `
✅総ウォレット数
%d wallet

✅登録AL数
%d 枚
`

	embed := &discordgo.MessageEmbed{
		Description: fmt.Sprintf(
			description,
			len(allWallets),
			allQuantity,
		),
	}

	_, err = s.ChannelMessageSendEmbed(internal.ChannelID().TEAM, embed)
	if err != nil {
		return errors.NewError("パネルメッセージを送信できません", err)
	}

	return nil
}
