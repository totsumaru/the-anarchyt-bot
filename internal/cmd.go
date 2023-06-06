package internal

// メッセージコマンド
const (
	// ルール
	CMD_Send_Rule = "!an-rule"

	// ガチャ
	CMD_Send_gatcha_Panel           = "!an-gatcha-panel"
	CMD_Send_gatcha_Add_Ticket_Role = "!an-add-role"
)

// InteractionのカスタムID
const (
	// ガチャを回す
	Interaction_CustomID_gatcha_Go   = "gatcha-go"
	Interaction_CustomID_gatcha_Open = "gatcha-open"
)
