package snowflake

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/log"
)

const (
	SERVICE        = "[SNOWFLAKE]"
	ENV_MACHINE_ID = "MACHINE_ID" // specific machine id
	MAX_DELAY      = 100          // max delay millisecond
)

const (
	TS_MASK         = 0x1FFFFFFFFFF // 41bit
	SN_MASK         = 0xFFF         // 12bit
	MACHINE_ID_MASK = 0x3FF         // 10bit
)

var machineId int64
var lastTimestamp int64
var serialNumber int16
var logger log.Logger

// RegisterService is
func RegisterService(e *echo.Echo) {
	e.GET("/uuid", handleGetId)
	logger = e.Logger()
}

type UUID struct {
	Uuid uint64 `json:"uuid"`
}

func handleGetId(context echo.Context) error {
	uniqueID := new(UUID)
	uniqueID.Uuid = generateID()

	return context.JSON(http.StatusOK, uniqueID)
}

func generateID() uint64 {

	env := os.Getenv(ENV_MACHINE_ID)
	if env != "" {
		logger.Debug("Found Machine ID. Using existing one.")
		mId, _ := strconv.Atoi(env)
		machineId = int64(mId)
	} else {

		machineId = (rand.Int63() & MACHINE_ID_MASK) << 12
		os.Setenv("MACHINE_ID", strconv.Itoa(int(machineId)))
		logger.Debug("Did not find Machine ID. Generating a new one.", machineId)
	}

	t := timestamp()

	if t < lastTimestamp {
		randomDelay()
	}

	if t == lastTimestamp {
		serialNumber++
	} else {
		serialNumber = 0
	}

	lastTimestamp = t

	var uuid uint64
	uuid |= (uint64(t) & TS_MASK) << 22
	uuid |= uint64(machineId)
	uuid |= uint64(serialNumber)
	return uuid
}

func timestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func randomDelay() {
	<-time.After(time.Duration(rand.Int63n(MAX_DELAY)) * time.Millisecond)
}
