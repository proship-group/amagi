package amagi

import (
	"fmt"
)

// HandleError http error handler TODO
func HandleError(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
