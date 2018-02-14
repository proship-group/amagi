package helpers

import (
	"fmt"
	"os"
	"testing"
)

func TestGetEnvIntValue(t *testing.T) {
	os.Setenv("TIMER_HOUR_ENV", "")
	fmt.Println(os.Getenv("TIMER_HOUR_ENV"))
	if GetEnvIntValue("TIMER_HOUR_ENV", 12) != 12 {
		t.Error(fmt.Errorf("got wrong value from GetEnvIntValue"))
	}
}
