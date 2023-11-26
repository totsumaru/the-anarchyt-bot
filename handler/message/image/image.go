package image

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

const BaseImageUrl = "https://arweave.net/OhewBZDW1nsLB0Cmq65YjrKMaIiRizUM3xr3Q7h6xZc"

var idToTrueID = map[int]int{
	86:  89,
	89:  86,
	182: 184,
	184: 182,
}

// 画像を送信します
func SendImage(s *discordgo.Session, m *discordgo.MessageCreate) error {
	// メッセージの内容を数値として解析
	number, err := strconv.Atoi(m.Content)
	if err != nil || number < 1 || number > 571 {
		return nil // 数字でない、または範囲外の場合は何もしない
	}

	// 画像のURLを構築
	if _, ok := idToTrueID[number]; ok {
		number = idToTrueID[number]
	}
	imageUrl := fmt.Sprintf("%s/%d.gif", BaseImageUrl, number)

	// 画像データを取得
	resp, err := http.Get(imageUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 画像を返信メッセージとして送信
	msg := discordgo.MessageSend{
		Reference: m.Reference(),
		Files: []*discordgo.File{
			{
				Name:   "image.gif",
				Reader: resp.Body,
			},
		},
	}
	_, err = s.ChannelMessageSendComplex(m.ChannelID, &msg)
	if err != nil {
		return err
	}

	return nil
}
