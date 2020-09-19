package json2env

import (
	"bytes"
	"context"
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestRun(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		os.Setenv("TEST", `{"test": "test"}`)
		err := Run(context.Background(), []string{"-keys", "test", "-envname", "TEST", "ls"}, &bytes.Buffer{}, &bytes.Buffer{}, []string{})
		if err != nil {
			t.Errorf("failed test %#v", err)
		}
	})
	t.Run("override existing environment", func(t *testing.T) {
		os.Setenv("TEST", `{"a": "after"}`)
		var outStream bytes.Buffer
		err := Run(context.Background(), []string{"-keys", "a", "-envname", "TEST", "env"}, &outStream, &bytes.Buffer{}, []string{"a=before"})
		if err != nil {
			t.Errorf("failed run %#v", err)
		}
		expect := []byte("a=after\n")
		if !bytes.Equal(outStream.Bytes(), expect) {
			t.Errorf("failed expect: %#v, got: %#v", string(expect), string(outStream.Bytes()))
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
				InputParams: []string{"-keys", "test", "-envname", "TEST", "ls"},
				Stdin:       `{deadbeaf`,
			},
			{
				Title:       "if stdin is nested json",
				InputParams: []string{"-keys", "test", "-envname", "TEST", "ls"},
				Stdin:       `{"test": {"test": "test"}}`,
			},
			{
				Title:       "if the key in keys option is not exists in json",
				InputParams: []string{"-keys", "notExists", "-envname", "TEST", "ls"},
				Stdin:       `{"test": "test"}`,
			},
			{
				Title:       "if command arg not exists",
				InputParams: []string{"-keys", "test", "-envname", "TEST"},
				Stdin:       `{"test": "test"}`,
			},
			{
				Title:       "if command not exists in path",
				InputParams: []string{"-keys", "test", "-envname", "TEST", "notExistsCommand"},
				Stdin:       `{"test": "test"}`,
			},
			{
				Title:       "if keys option is not provided",
				InputParams: []string{"-envname", "TEST", "ls"},
				Stdin:       `{"test": "test"}`,
			},
			{
				Title:       "if envname option is not exists",
				InputParams: []string{"-keys", "test", "ls"},
				Stdin:       `{"test": "test"}`,
			},
		}
		for _, pattern := range patterns {
			os.Setenv("TEST", pattern.Stdin)
			err := Run(context.Background(), pattern.InputParams, &bytes.Buffer{}, &bytes.Buffer{}, []string{})
			if err == nil {
				t.Errorf("error not occurred, failed test.\n Test title: %s", pattern.Title)
			}
		}
	})
}
func TestMakeNewEnv(t *testing.T) {
	t.Run("all environment variable check", func(t *testing.T) {
		env, err := makeNewEnv([]string{"a=b", "b=b"}, map[string]string{"b": "c", "d": "e"}, []string{"b", "d"})
		if err != nil {
			t.Errorf("failed err:%#v", err)
		}
		sort.Strings(env)
		expect := []string{"a=b", "b=b", "b=c", "d=e"}
		sort.Strings(expect)
		if !reflect.DeepEqual(env, expect) {
			t.Errorf("failed getEnv test expect: %#v, got: %#v", expect, env)
		}
	})
	t.Run("some environment variable check", func(t *testing.T) {
		env, err := makeNewEnv([]string{"a=b", "b=b"}, map[string]string{"b": "c", "d": "e"}, []string{"b"})
		if err != nil {
			t.Errorf("failed err:%#v", err)
		}
		sort.Strings(env)
		expect := []string{"a=b", "b=b", "b=c"}
		sort.Strings(expect)
		if !reflect.DeepEqual(env, expect) {
			t.Errorf("failed getEnv test expect: %#v, got: %#v", expect, env)
		}
	})
	t.Run("key does not exist in json", func(t *testing.T) {
		_, err := makeNewEnv([]string{"a=b", "b=b"}, map[string]string{"b": "c", "d": "e"}, []string{"x"})
		if err == nil {
			t.Fatal("failed err")
		}
	})
}
