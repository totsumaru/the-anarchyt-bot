package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
	"os"
)

// スラッシュコマンドを登録します
func RegisterSlashCommand(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if m.ChannelID != internal.ChannelID().TEST {
		return nil
	}

	if err := registerCommand(s, m.GuildID); err != nil {
		return errors.NewError("コマンドを登録できません", err)
	}

	_, err := s.ChannelMessageSend(m.ChannelID, "Slashコマンドを追加しました")
	if err != nil {
		return errors.NewError("完了メッセージを送信できません", err)
	}

	return nil
}

// コマンドを登録します
func registerCommand(session *discordgo.Session, guildID string) error {
	commands := []discordgo.ApplicationCommand{
		{
			Name:        internal.Slash_CMD_MyRoles,
			Description: "自分のロールを確認できます",
			Options:     []*discordgo.ApplicationCommandOption{},
		},
	}

	for _, command := range commands {
		_, err := session.ApplicationCommandCreate(os.Getenv("BOT_APPLICATION_ID"), guildID, &command)
		if err != nil {
			return errors.NewError("コマンドを登録できません Name: ", command.Name)
		}
	}

	return nil
}
