package sql

import (
	"github.com/jinzhu/gorm"
)

// KeepAliveExec keep alive query and exec
func KeepAliveExec(db *gorm.DB) error {
	if err := db.Exec(`SELECT 1`).Error; err != nil {
		return err
	}
	return nil
}
