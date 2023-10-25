package address

import "github.com/techstart35/the-anarchy-bot/internal"

// ロールに対するMint上限です
var RoleMaxMintMap = map[string]int{
	internal.RoleID().AL:            2,
	internal.RoleID().BRONZE:        1,
	internal.RoleID().SILVER:        1,
	internal.RoleID().GOLD:          1,
	internal.RoleID().PLATINUM:      1,
	internal.RoleID().DIAMOND:       1,
	internal.RoleID().CRAZY:         1,
	internal.RoleID().CHAINSAW_CLUB: 1,
	internal.RoleID().PHYSICAL:      1,
}

// Mint上限を取得します
func MaxMintQuantity(roles []string) int {
	res := 0
	for _, role := range roles {
		if quantity, ok := RoleMaxMintMap[role]; ok {
			res += quantity
		}
	}

	if res >= 5 {
		return 5
	}

	return res
}
