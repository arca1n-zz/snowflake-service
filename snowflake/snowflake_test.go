package snowflake

import (
	"testing"
	"time"

	"github.com/labstack/echo"
)

func TestRandomDelay(t *testing.T) {
	randomDelay()
}

func TestTimeStamp(t *testing.T) {
	start := time.Now().UnixNano() / int64(time.Millisecond)
	timestamp := timestamp()
	if timestamp < start {
		t.Fatalf("Start time %ll is greater than timestamp %ll", start, timestamp)
	}
}

func TestSnowflake(t *testing.T) {

	echo := echo.New()
	RegisterService(echo)
	uuid := generateID()
	if uuid == 0 {
		t.Fail()
		t.Logf("UUID generated is invalid", uuid)
	} else {
		t.Logf("UUID generated is %ll", uuid)
	}
}
