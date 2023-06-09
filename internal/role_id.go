package internal

import "os"

type Role struct {
	TICKET      string
	VERIFIED    string
	PRIZE1      string
	PRIZE2      string
	PRIZE3      string
	INVITATION1 string
	INVITATION2 string
}

// ロールID
func RoleID() Role {
	if os.Getenv("ENV") == "dev" {
		// テスト環境
		return Role{
			TICKET:      "1115563590825553941", // ガチャチケット
			VERIFIED:    "1112562612044038295", // 入場券
			PRIZE1:      "1115569765101076510", // 当たり1
			PRIZE2:      "1115570175949951067", // 当たり2
			PRIZE3:      "1115570239183269930", // 当たり3
			INVITATION1: "1116535909093998664", // 招待券1
			INVITATION2: "1116536092573839400", // 招待券2
		}
	} else {
		// 本番環境
		return Role{
			TICKET:      "1115563590825553941", // ガチャチケット
			VERIFIED:    "1112562612044038295", // 入場券
			PRIZE1:      "1115569765101076510", // 当たり1
			PRIZE2:      "1115570175949951067", // 当たり2
			PRIZE3:      "1115570239183269930", // 当たり3
			INVITATION1: "1116535909093998664", // 招待券1
			INVITATION2: "1116536092573839400", // 招待券2
		}
	}
}
