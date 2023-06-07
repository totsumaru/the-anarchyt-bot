package invitation

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"os"
	"strings"
	"time"
)

// 招待を作成します
// - 1回のみ
// - 7日間
func CreateInvitation(s *discordgo.Session, m *discordgo.MessageCreate) error {
	// 管理者のチャンネルのみで実行可能
	adminCmdChannelID := os.Getenv("ADMIN_CHANNEL_ID")
	if m.ChannelID != adminCmdChannelID {
		return nil
	}

	// 20秒待つメッセージを送信する
	_, err := s.ChannelMessageSend(m.ChannelID, "Please wait 20 seconds...")
	if err != nil {
		return errors.NewError("20秒待つメッセージを送信できません", err)
	}

	var maxAge int
	var description string

	maxAge = 604800
	description = `
- 7日
- 1回しか使用できない

%s
`

	urls := make([]string, 0)
	ivChID := os.Getenv("INVITATION_LINK_CHANNEL_ID")

	// 10個のリンクを発行します
	for i := 0; i < 10; i++ {
		iv := discordgo.Invite{
			MaxAge:    maxAge,
			MaxUses:   1, // 1回限り
			Unique:    true,
			Temporary: false,
		}

		invite, err := s.ChannelInviteCreate(ivChID, iv)
		if err != nil {
			return errors.NewError("招待リンクを作成できません", err)
		}

		url := fmt.Sprintf("https://discord.gg/%s", invite.Code)
		urls = append(urls, url)

		time.Sleep(1 * time.Second) // API制限回避のために1秒待つ
	}

	// メッセージを送信
	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(description, strings.Join(urls, "\n")))
	if err != nil {
		return errors.NewError("メッセージを送信できません", err)
	}

	return nil
}
