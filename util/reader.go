package util

import (
	"errors"
	"io"
	"testing"

	"github.com/cloudwego/eino/schema"
	"github.com/stretchr/testify/require"
)

func PrintContentByReader(tb testing.TB, reader *schema.StreamReader[*schema.Message]) {
	print("============ reader starting ============\n\n")
	defer func() {
		print("============ reader finished ============\n\n")
	}()
	for {
		chunk, err := reader.Recv()
		if errors.Is(err, io.EOF) {
			print("\n\n")
			return
		}
		require.NoError(tb, err)
		print(chunk.Content)
	}
}
