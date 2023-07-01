package internal

import "os"

type User struct {
	TOTSUMARU string
	MUG       string
}

// ユーザーID
func UserID() User {
	if os.Getenv("ENV") == "dev" {
		// テスト環境
		return User{
			TOTSUMARU: "960104306151948328",
			MUG:       "954798318742044672",
		}
	} else {
		// 本番環境
		return User{
			TOTSUMARU: "960104306151948328",
			MUG:       "954798318742044672",
		}
	}
}
