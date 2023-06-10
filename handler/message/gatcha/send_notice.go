package gatcha

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
	"strings"
)

// 毎日の通知メッセージを送信します
func SendNotice(s *discordgo.Session) error {
	// 前回のメッセージを削除します
	messages, err := s.ChannelMessages(internal.ChannelID().GATCHA, 1, "", "", "")
	if err != nil {
		return errors.NewError("最新メッセージを取得できません", err)
	}

	latestMessage := messages[0]
	if strings.Contains(latestMessage.Content, "Notice") {
		msg := fmt.Sprintf(
			"<@&%s> おはようございます！今日もガチャで運試しをしてみましょう！！",
			internal.RoleID().GATCHA_NOTICE,
		)

		_, err = s.ChannelMessageSend(internal.ChannelID().GATCHA, msg)
		if err != nil {
			return errors.NewError("通知メッセージを送信できません", err)
		}
	}

	return nil
}
