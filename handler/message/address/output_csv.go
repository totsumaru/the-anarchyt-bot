package address

import (
	"bytes"
	"encoding/csv"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
	"github.com/techstart35/the-anarchy-bot/internal/db"
)

// 全てのレコードをcsvにして、Discordに送信します
func OutputCSV(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if m.ChannelID != internal.ChannelID().TEST {
		return nil
	}

	wallets, err := db.FindAll()
	if err != nil {
		return errors.NewError("全ての情報を取得できませんでした", err)
	}

	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)

	// Write CSV header
	writer.Write([]string{"ID", "Address", "Quantity"})

	// Write wallets data to csv
	for _, wallet := range wallets {
		if err = writer.Write([]string{
			wallet.ID,
			wallet.Address,
			strconv.Itoa(wallet.Quantity),
		}); err != nil {
			return errors.NewError("CSVへの書き込みに失敗しました", err)
		}
	}

	// Ensure all data is written to buffer
	writer.Flush()

	if writer.Error() != nil {
		return errors.NewError("CSVのフラッシュ中にエラーが発生しました", writer.Error())
	}

	// Create a new file from our buffer to send
	file := &discordgo.File{
		Name:   "wallets.csv",
		Reader: buf,
	}

	ms := &discordgo.MessageSend{
		Content: "Here are the wallets data",
		Files:   []*discordgo.File{file},
	}

	// Send the message with the CSV attached
	_, err = s.ChannelMessageSendComplex(m.ChannelID, ms)
	if err != nil {
		return errors.NewError("Discordへのメッセージ送信中にエラーが発生しました", err)
	}

	return nil
}
