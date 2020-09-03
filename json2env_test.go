package json2env

import (
	"bytes"
	"context"
	"reflect"
	"sort"
	"testing"
)

func TestRun(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		stdin := bytes.NewBufferString(`{"test": "test"}`)
		err := Run(context.Background(), []string{"-keys", "test", "ls"}, &bytes.Buffer{}, &bytes.Buffer{}, stdin, []string{})
		if err != nil {
			t.Fatalf("failed test %#v", err)
		}
	})
	t.Run("all environment variable check", func(t *testing.T) {
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
	})
	t.Run("some environment variable check", func(t *testing.T) {
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
	})
	t.Run("override existing environment", func(t *testing.T) {
		stdin := bytes.NewBufferString(`{"a": "after"}`)
		var outStream bytes.Buffer
		err := Run(context.Background(), []string{"-keys", "a", "env"}, &outStream, &bytes.Buffer{}, stdin, []string{"a=before"})
		if err != nil {
			t.Fatalf("failed run %#v", err)
		}
		expect := []byte("a=after\n")
		if !bytes.Equal(outStream.Bytes(), expect) {
			t.Fatalf("failed expect: %#v, got: %#v", string(expect), string(outStream.Bytes()))
		}
	})
	t.Run("error if stdin is not json", func(t *testing.T) {
		stdin := bytes.NewBufferString("{deadbeaf")
		err := Run(context.Background(), []string{"-keys", "test", "ls"}, &bytes.Buffer{}, &bytes.Buffer{}, stdin, []string{})
		if err == nil {
			t.Fatalf("error not occurred, failed test")
		}
	})
	t.Run("error if stdin is nested json", func(t *testing.T) {
		stdin := bytes.NewBufferString(`{"test": {"test": "test"}}`)
		err := Run(context.Background(), []string{"-keys", "test", "ls"}, &bytes.Buffer{}, &bytes.Buffer{}, stdin, []string{})
		if err == nil {
			t.Fatalf("error not occurred, failed test")
		}
	})
	t.Run("error if the key in keys option is not exists in json", func(t *testing.T) {
		stdin := bytes.NewBufferString(`{"test": "test"}`)
		err := Run(context.Background(), []string{"-keys", "notExists", "ls"}, &bytes.Buffer{}, &bytes.Buffer{}, stdin, []string{})
		if err == nil {
			t.Fatalf("error not occurred, failed test")
		}
	})
	t.Run("error if command arg not exists", func(t *testing.T) {
		stdin := bytes.NewBufferString(`{"test": "test"}`)
		err := Run(context.Background(), []string{"-keys", "test"}, &bytes.Buffer{}, &bytes.Buffer{}, stdin, []string{})
		if err == nil {
			t.Fatalf("error not occurred, failed test %#v", err)
		}
	})
	t.Run("error if command not exists in path", func(t *testing.T) {
		stdin := bytes.NewBufferString(`{"test": "test"}`)
		err := Run(context.Background(), []string{"-keys", "test", "notExistsCommand"}, &bytes.Buffer{}, &bytes.Buffer{}, stdin, []string{})
		if err == nil {
			t.Fatalf("error not occurred, failed test %#v", err)
		}
	})
}
func TestNotExistedKey(t *testing.T) {
	_, err := makeNewEnv([]string{"a=b", "b=b"}, map[string]string{"b": "c", "d": "e"}, []string{"x"})
	if err == nil {
		t.Fatal("failed err")
	}
}
