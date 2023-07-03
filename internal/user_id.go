package internal

import "os"

type User struct {
	TOTSUMARU string
	MUG       string
	OTOUSAN   string
}

// ユーザーID
func UserID() User {
	if os.Getenv("ENV") == "dev" {
		// テスト環境
		return User{
			TOTSUMARU: "960104306151948328",
			MUG:       "954798318742044672",
			OTOUSAN:   "795588909576749067",
		}
	} else {
		// 本番環境
		return User{
			TOTSUMARU: "960104306151948328",
			MUG:       "954798318742044672",
			OTOUSAN:   "795588909576749067",
		}
	}
}
