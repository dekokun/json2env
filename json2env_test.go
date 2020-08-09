package json2env

import (
	"bytes"
	"context"
	"reflect"
	"sort"
	"testing"
)

func TestNormal(t *testing.T) {
	stdin := bytes.NewBufferString("{}")
	err := Run(context.Background(), []string{"ls"}, &bytes.Buffer{}, &bytes.Buffer{}, stdin, []string{})
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
}

func TestErrorIfNotJson(t *testing.T) {
	stdin := bytes.NewBufferString("{deadbeaf")
	err := Run(context.Background(), []string{"ls"}, &bytes.Buffer{}, &bytes.Buffer{}, stdin, []string{})
	if err == nil {
		t.Fatalf("error not occurred, failed test %#v", err)
	}
}

func TestEnv(t *testing.T) {
	stdin := bytes.NewBufferString("{\"a\": \"after\"}")
	var outStream bytes.Buffer
	err := Run(context.Background(), []string{"env"}, &outStream, &bytes.Buffer{}, stdin, []string{"a=before"})
	if err != nil {
		t.Fatalf("failed run %#v", err)
	}
	expect := []byte("a=after\n")
	if !bytes.Equal(outStream.Bytes(), expect) {
		t.Fatalf("failed expect: %#v, got: %#v", string(expect), string(outStream.Bytes()))
	}
}

func TestGetEnv(t *testing.T) {
	env := makeNewEnv([]string{"a=b", "b=b"}, map[string]string{"b": "c", "d": "e"})
	sort.Strings(env)
	expect := []string{"a=b", "b=b", "b=c", "d=e"}
	sort.Strings(expect)
	if !reflect.DeepEqual(env, expect) {
		t.Fatalf("failed getEnv test expect: %#v, got: %#v", expect, env)
	}
}
