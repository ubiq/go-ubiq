package openrpc_test

import (
	"testing"

	"github.com/ubiq/go-ubiq/v3/internal/openrpc"
	"github.com/ubiq/go-ubiq/v3/rpc"
)

func TestDefaultSchema(t *testing.T) {
	if err := rpc.SetDefaultOpenRPCSchemaRaw(openrpc.OpenRPCSchema); err != nil {
		t.Fatal(err)
	}
}
