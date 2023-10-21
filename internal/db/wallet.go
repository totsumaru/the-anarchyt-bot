package db

import (
	defaultErrors "errors"
	"os"

	"github.com/techstart35/the-anarchy-bot/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ウォレットアドレス提出のテーブルです
type Wallet struct {
	ID       string `json:"id"`
	Address  string `json:"address"`
	Quantity int    `json:"quantity"`
}

// DBに接続します
func ConnectDB() {
	dialector := postgres.Open(os.Getenv("DB_URL"))
	db, err := gorm.Open(dialector, &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(errors.NewError("DBに接続できません", err))
	}

	// テーブルが存在していない場合のみテーブルを作成します
	// 存在している場合はスキーマを同期します
	if err = db.AutoMigrate(&Wallet{}); err != nil {
		panic(errors.NewError("テーブルのスキーマが一致しません", err))
	}

	DB = db
}

// Upsertします
func Upsert(tx *gorm.DB, id, address string, quantity int) error {
	var wallet Wallet

	// Addressをキーとして検索
	if err := tx.Where("id = ?", id).First(&wallet).Error; err != nil {
		// レコードが見つからない場合は新しいレコードを作成
		if defaultErrors.Is(err, gorm.ErrRecordNotFound) {
			wallet = Wallet{
				ID:       id,
				Address:  address,
				Quantity: quantity,
			}
			if err = tx.Create(&wallet).Error; err != nil {
				return errors.NewError("レコードを作成できません", err)
			}
		} else {
			return errors.NewError("Addressからレコードを取得できません", err)
		}
	} else {
		// レコードが見つかった場合はQuantityを更新
		wallet.Address = address
		wallet.Quantity = quantity
		if err = tx.Save(&wallet).Error; err != nil {
			return errors.NewError("レコードを更新できません", err)
		}
	}

	return nil
}

// 削除します
func Remove(tx *gorm.DB, id string) error {
	var wallet Wallet

	// IDをキーとして検索
	if err := tx.Where("id = ?", id).First(&wallet).Error; err != nil {
		// レコードが見つからない場合
		if defaultErrors.Is(err, gorm.ErrRecordNotFound) {
			return errors.NewError("削除するレコードが見つかりません", err)
		}
		// その他のエラーの場合
		return errors.NewError("レコードの検索に失敗しました", err)
	}

	// レコードを削除
	if err := tx.Delete(&wallet).Error; err != nil {
		return errors.NewError("レコードの削除に失敗しました", err)
	}

	return nil
}

// IDで取得します
func FindByID(id string) (Wallet, error) {
	var wallet Wallet

	// IDをキーとして検索
	if err := DB.Where("id = ?", id).First(&wallet).Error; err != nil {
		// レコードが見つからない場合
		if defaultErrors.Is(err, gorm.ErrRecordNotFound) {
			return Wallet{}, nil
		}
		// その他のエラーの場合
		return Wallet{}, errors.NewError("レコードの取得に失敗しました", err)
	}

	return wallet, nil
}

// 全ての情報を取得します
func FindAll() ([]Wallet, error) {
	var wallets []Wallet

	// 全てのレコードを取得
	if err := DB.Find(&wallets).Error; err != nil {
		// その他のエラーの場合
		return nil, errors.NewError("全てのレコードの取得に失敗しました", err)
	}

	return wallets, nil
}
