package json2env

import (
	"bytes"
	"context"
	"reflect"
	"sort"
	"testing"
)

func TestNormal(t *testing.T) {
	stdin := bytes.NewBufferString("{\"test\": \"test\"}")
	err := Run(context.Background(), []string{"-keys", "test", "ls"}, &bytes.Buffer{}, &bytes.Buffer{}, stdin, []string{})
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
}

func TestErrorIfNotJson(t *testing.T) {
	stdin := bytes.NewBufferString("{deadbeaf")
	err := Run(context.Background(), []string{"-keys", "test", "ls"}, &bytes.Buffer{}, &bytes.Buffer{}, stdin, []string{})
	if err == nil {
		t.Fatalf("error not occurred, failed test")
	}
}

func TestErrorIfNestedJson(t *testing.T) {
	stdin := bytes.NewBufferString("{\"test\": {\"test\": \"test\"}}")
	err := Run(context.Background(), []string{"-keys", "test", "ls"}, &bytes.Buffer{}, &bytes.Buffer{}, stdin, []string{})
	if err == nil {
		t.Fatalf("error not occurred, failed test")
	}
}

func TestErrorIfSpecifyNotExistsKeysJson(t *testing.T) {
	stdin := bytes.NewBufferString("{\"test\": \"test\"}")
	err := Run(context.Background(), []string{"-keys", "notExists", "ls"}, &bytes.Buffer{}, &bytes.Buffer{}, stdin, []string{})
	if err == nil {
		t.Fatalf("error not occurred, failed test")
	}
}
func TestErrorIfCommandNotExists(t *testing.T) {
	stdin := bytes.NewBufferString("{\"test\": \"test\"}")
	err := Run(context.Background(), []string{"-keys", "test"}, &bytes.Buffer{}, &bytes.Buffer{}, stdin, []string{})
	if err == nil {
		t.Fatalf("error not occurred, failed test %#v", err)
	}
}
func TestEnv(t *testing.T) {
	stdin := bytes.NewBufferString("{\"a\": \"after\"}")
	var outStream bytes.Buffer
	err := Run(context.Background(), []string{"-keys", "a", "env"}, &outStream, &bytes.Buffer{}, stdin, []string{"a=before"})
	if err != nil {
		t.Fatalf("failed run %#v", err)
	}
	expect := []byte("a=after\n")
	if !bytes.Equal(outStream.Bytes(), expect) {
		t.Fatalf("failed expect: %#v, got: %#v", string(expect), string(outStream.Bytes()))
	}
}

func TestGetEnv(t *testing.T) {
	env, err := makeNewEnv([]string{"a=b", "b=b"}, map[string]string{"b": "c", "d": "e"}, []string{"b", "d"})
	if err != nil {
		t.Fatalf("failed err:%#v", err)
	}
	sort.Strings(env)
	expect := []string{"a=b", "b=b", "b=c", "d=e"}
	sort.Strings(expect)
	if !reflect.DeepEqual(env, expect) {
		t.Fatalf("failed getEnv test expect: %#v, got: %#v", expect, env)
	}
}

func TestGetSomeEnv(t *testing.T) {
	env, err := makeNewEnv([]string{"a=b", "b=b"}, map[string]string{"b": "c", "d": "e"}, []string{"b"})
	if err != nil {
		t.Fatalf("failed err:%#v", err)
	}
	sort.Strings(env)
	expect := []string{"a=b", "b=b", "b=c"}
	sort.Strings(expect)
	if !reflect.DeepEqual(env, expect) {
		t.Fatalf("failed getEnv test expect: %#v, got: %#v", expect, env)
	}
}

func TestNotExistedKey(t *testing.T) {
	_, err := makeNewEnv([]string{"a=b", "b=b"}, map[string]string{"b": "c", "d": "e"}, []string{"x"})
	if err == nil {
		t.Fatal("failed err")
	}
}
