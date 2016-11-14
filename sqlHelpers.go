package amagi

import (
	"github.com/b-eee/amagi/api/sql"
	"github.com/jinzhu/gorm"
)

// KeepAlive keep alive query and exec
func KeepAlive(db *gorm.DB) error {
	return sql.KeepAliveExec(db)
}
