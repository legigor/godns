package main

import (
	"context"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_server(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := run(ctx, os.Stdout)
		require.NoError(t, err)
	}()
	defer cancel()

	// TODO:
}
