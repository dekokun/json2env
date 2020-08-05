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
	err := Run(context.Background(), []string{}, &bytes.Buffer{}, &bytes.Buffer{}, stdin)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
}

func TestErrorIfNotJson(t *testing.T) {
	stdin := bytes.NewBufferString("{deadbeaf")
	err := Run(context.Background(), []string{}, &bytes.Buffer{}, &bytes.Buffer{}, stdin)
	if err == nil {
		t.Fatalf("error not occurred, failed test")
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
