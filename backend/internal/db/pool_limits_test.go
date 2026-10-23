package db

import (
	"context"
	"math"
	"strings"
	"testing"
)

func TestNewPool_RejectsConnectionCountOverflow(t *testing.T) {
	for _, config := range []PoolConfig{{MaxConns: math.MaxInt32 + 1}, {MinConns: math.MaxInt32 + 1}} {
		pool, err := NewPool(context.Background(), "postgres://unused@127.0.0.1:1/unused", config)
		if pool != nil || err == nil || !strings.Contains(err.Error(), "int32 range") {
			t.Fatalf("config=%+v, pool=%v, err=%v", config, pool, err)
		}
	}
}
