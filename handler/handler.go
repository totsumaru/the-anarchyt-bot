package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/handler/interaction"
	"github.com/techstart35/the-anarchy-bot/handler/message"
	"github.com/techstart35/the-anarchy-bot/handler/reaction"
)

// ハンドラを追加します
func AddHandler(s *discordgo.Session) {
	s.AddHandler(message.MessageCreateHandler)
	s.AddHandler(interaction.InteractionCreateHandler)
	s.AddHandler(reaction.ReactionCreateHandler)
}
