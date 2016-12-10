package amagi

import (
	"github.com/rs/xid"
)

// GenerateUUID generate uuid
func GenerateUUID() string {
	guid := xid.New()

	return guid.String()
}
