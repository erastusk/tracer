package utils

import (
	"net"
	"time"

	"github/erastusk/tracer/internal/types"

	"github.com/fatih/color"
)

var (
	Done    chan struct{}
	Red     *color.Color
	Green   *color.Color
	Options []string
)

type kafkaApp struct {
	App types.KafkaApps
}

func NewKafkaApp(a types.KafkaApps) kafkaApp {
	return kafkaApp{
		App: a,
	}
}

func InitColors() {
	Red = color.New(color.FgRed)
	Green = color.New(color.FgGreen)
	Options = []string{"Please enter endpoint", "Please enter username", "Please enter password"}
}

type TestConnectivityImpl struct{}

func NewTestConnectivityImpl() types.TCPDial {
	return TestConnectivityImpl{}
}

func (tc TestConnectivityImpl) TCPDial(server string) error {
	// Dial the server
	conn, err := net.DialTimeout("tcp", server, time.Second*2)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}
