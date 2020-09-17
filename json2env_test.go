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
		err := Run(context.Background(), []string{"-keys", "test", "ls"}, &bytes.Buffer{}, &bytes.Buffer{}, []string{})
		if err != nil {
			t.Fatalf("failed test %#v", err)
		}
	})
	t.Run("override existing environment", func(t *testing.T) {
		stdin := bytes.NewBufferString(`{"a": "after"}`)
		var outStream bytes.Buffer
		err := Run(context.Background(), []string{"-keys", "a", "env"}, &outStream, &bytes.Buffer{}, []string{"a=before"})
		if err != nil {
			t.Fatalf("failed run %#v", err)
		}
		expect := []byte("a=after\n")
		if !bytes.Equal(outStream.Bytes(), expect) {
			t.Fatalf("failed expect: %#v, got: %#v", string(expect), string(outStream.Bytes()))
		}
	})
	t.Run("check error", func(t *testing.T) {
		patterns := []struct {
			Title       string
			InputParams []string
			Stdin       string
		}{
			{
				Title:       "if stdin is not json",
				InputParams: []string{"-keys", "test", "ls"},
				Stdin:       `{deadbeaf`,
			},
			{
				Title:       "if stdin is nested json",
				InputParams: []string{"-keys", "test", "ls"},
				Stdin:       `{"test": {"test": "test"}}`,
			},
			{
				Title:       "if the key in keys option is not exists in json",
				InputParams: []string{"-keys", "notExists", "ls"},
				Stdin:       `{"test": "test"}`,
			},
			{
				Title:       "if command arg not exists",
				InputParams: []string{"-keys", "test"},
				Stdin:       `{"test": "test"}`,
			},
			{
				Title:       "if command not exists in path",
				InputParams: []string{"-keys", "test", "notExistsCommand"},
				Stdin:       `{"test": "test"}`,
			},
			{
				Title:       "if keys option is not provided",
				InputParams: []string{"ls"},
				Stdin:       `{"test": "test"}`,
			},
		}
		for _, pattern := range patterns {
			stdin := bytes.NewBufferString(pattern.Stdin)
			err := Run(context.Background(), pattern.InputParams, &bytes.Buffer{}, &bytes.Buffer{}, []string{})
			if err == nil {
				t.Fatalf("error not occurred, failed test.\n Test title: %s", pattern.Title)
			}
		}
	})
}
func TestMakeNewEnv(t *testing.T) {
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
	t.Run("key does not exist in json", func(t *testing.T) {
		_, err := makeNewEnv([]string{"a=b", "b=b"}, map[string]string{"b": "c", "d": "e"}, []string{"x"})
		if err == nil {
			t.Fatal("failed err")
		}
	})
}
