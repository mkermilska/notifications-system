package nats

import (
	"fmt"
	"testing"

	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const (
	testPort       = 8369
	testServiceID  = "notifications-api"
	testStreamName = "notifications-stream"
	streamSubjects = "EMAIL_MESSAGE"
)

func RunServerOnPort(port int) *server.Server {
	opts := natsserver.DefaultTestOptions
	opts.Port = port
	return RunServerWithOptions(&opts)
}

func RunServerWithOptions(opts *server.Options) *server.Server {
	return natsserver.RunServer(opts)
}

func TestCreateConnectionAndStreamContext(t *testing.T) {
	s := RunServerOnPort(testPort)
	defer s.Shutdown()

	natsSvc := NewService(testStreamName, testServiceID, zap.NewNop())
	err := natsSvc.Connect(fmt.Sprintf("%s:%d", "127.0.0.1", testPort))

	assert.Nil(t, err)
	assert.NotNil(t, natsSvc.js)
	defer natsSvc.conn.Close()
}
