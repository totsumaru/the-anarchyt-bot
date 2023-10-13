package address

import "github.com/techstart35/the-anarchy-bot/internal"

// ロールに対するMint上限です
var RoleMaxMintMap = map[string]int{
	internal.RoleID().AL: 2,
}

// Mint上限を取得します
func MaxMintQuantity(roles []string) int {
	res := 0
	for _, role := range roles {
		if quantity, ok := RoleMaxMintMap[role]; ok {
			res += quantity
		}
	}

	return res
}
