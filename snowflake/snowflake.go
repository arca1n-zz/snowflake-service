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

// Service name
const Service = "[SNOWFLAKE]"

// EnvMachineID is the environment varialble name for Machine ID
const EnvMachineID string = "MACHINE_ID" // specific machine id
// MaxDelay is the maximum delay for the random delay method
const MaxDelay = 100 // max delay millisecond

const tsMask = 0x1FFFFFFFFFF // 41bit
const snMask = 0xFFF         // 12bit
const machineIDMask = 0x3FF  // 10bit

var machineID int64
var lastTimestamp int64
var serialNumber int16
var logger log.Logger

// RegisterService is
func RegisterService(e *echo.Echo) {
	e.GET("/uuid", handleGetID)
	logger = e.Logger()
}

// UUID struct is the type that contains the generated UUID
type UUID struct {
	ID uint64 `json:"uuid"`
}

func handleGetID(context echo.Context) error {
	uniqueID := new(UUID)
	uniqueID.ID = generateID()

	return context.JSON(http.StatusOK, uniqueID)
}

func generateID() uint64 {

	env := os.Getenv(EnvMachineID)
	if env != "" {
		logger.Debug("Found Machine ID. Using existing one.")
		mID, _ := strconv.Atoi(env)
		machineID = int64(mID)
	} else {

		machineID = (rand.Int63() & machineIDMask) << 12
		os.Setenv(EnvMachineID, strconv.Itoa(int(machineID)))
		logger.Debug("Did not find Machine ID. Generating a new one.", machineID)
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
	uuid |= (uint64(t) & tsMask) << 22
	uuid |= uint64(machineID)
	uuid |= uint64(serialNumber)
	return uuid
}

func timestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func randomDelay() {
	<-time.After(time.Duration(rand.Int63n(MaxDelay)) * time.Millisecond)
}
