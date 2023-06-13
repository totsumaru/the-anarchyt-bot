package invitation

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
	"strings"
	"time"
)

// 管理者用の招待を作成します
// - 1回のみ
// - 7日間
func CreateInvitationForAdmin(s *discordgo.Session, m *discordgo.MessageCreate) error {
	// 管理者のチャンネルのみで実行可能
	if m.ChannelID != internal.ChannelID().TEST {
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

	// 10個のリンクを発行します
	for i := 0; i < 10; i++ {
		iv := discordgo.Invite{
			MaxAge:    maxAge,
			MaxUses:   1, // 1回限り
			Unique:    true,
			Temporary: false,
		}

		invite, err := s.ChannelInviteCreate(internal.ChannelID().START, iv)
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
