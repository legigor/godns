package dns

import (
	"context"
	"github.com/stretchr/testify/require"
	"log/slog"
	"net"
	"os"
	"testing"
	"time"
)

func Test_server(t *testing.T) {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()

	srv := NewServer(logger, ctx)
	addr, err := srv.Start()
	require.NoError(t, err)

	logger.Info("server: " + addr.String())

	client, err := net.DialUDP("udp", nil, addr)
	require.NoError(t, err)

	err = client.SetDeadline(time.Now().Add(5 * time.Second))
	require.NoError(t, err)

	message := "Hello, UDP server!"

	_, err = client.Write([]byte(message))
	require.NoError(t, err)

	buffer := make([]byte, 1024)
	n, _, err := client.ReadFromUDP(buffer)
	require.NoError(t, err)

	response := string(buffer[:n])
	logger.Info("response: " + response)
}
