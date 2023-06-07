package internal

// メッセージコマンド
const (
	CMD_Send_Rule                   = "!an-rule"         // ルール
	CMD_Send_gatcha_Panel           = "!an-gatcha-panel" // ガチャのパネル
	CMD_Send_gatcha_Add_Ticket_Role = "!an-add-role"     // チケットルール付与
	CMD_Send_verify_Panel           = "!an-verify-panel" // Verifyのパネル
	CMD_Create_Invitation           = "!an-invitation"   // 招待リンク発行
	CMD_Member                      = "member"           // 参加人数を取得
)

// InteractionのカスタムID
const (
	// ガチャ
	Interaction_CustomID_gatcha_Go   = "gatcha-go"
	Interaction_CustomID_gatcha_Open = "gatcha-open"
	// Verify
	Interaction_CustomID_Verify = "verify"
)
