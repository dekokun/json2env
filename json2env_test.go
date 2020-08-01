package json2env

import (
	"bytes"
	"context"
	"testing"
)

func TestNormal(t *testing.T) {
	stdin := bytes.NewBufferString("{}")
	err := Run(context.Background(), []string{}, &bytes.Buffer{}, &bytes.Buffer{}, stdin)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
}
