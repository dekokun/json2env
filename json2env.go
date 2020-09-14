package json2env

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

const cmdName = "json2env"

// Run the json2env
func Run(ctx context.Context, argv []string, outStream, errStream io.Writer, inStream io.Reader, env []string) error {
	log.SetOutput(errStream)
	fs := flag.NewFlagSet(
		fmt.Sprintf("%s (v%s rev:%s)", cmdName, version, revision), flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", fs.Name())
		fmt.Fprintf(flag.CommandLine.Output(), "  SECRETS='{\"key\":\"value\"}'; %s --input SECRETS --keys \"key\" /path/to/command [...]\n\n", cmdName)
		fs.PrintDefaults()
	}
	fs.SetOutput(errStream)
	ver := fs.Bool("version", false, "display version")
	keys := fs.String("keys", "", "conmma separated environment variable names that you want to export")

	if err := fs.Parse(argv); err != nil {
		return err
	}
	command := fs.Args()
	if len(command) == 0 {
		return errors.New("you must specify one command name")
	}
	if *ver {
		return printVersion(outStream)
	}

	if *keys == "" {
		return errors.New("keys option is required")
	}

	var envJSON map[string]string
	err := json.NewDecoder(inStream).Decode(&envJSON)
	if err != nil {
		switch err := errors.Cause(err).(type) {
		case *json.UnmarshalTypeError:
			return errors.Wrap(err, "failed to decode input, json must be flat key value map, example: {\"key1\":\"value1\", \"key2\":\"value2\"}")
		default:
			return errors.Wrap(err, "failed to decode input")
		}
	}

	newEnv, err := makeNewEnv(env, envJSON, strings.Split(*keys, ","))
	if err != nil {
		return errors.Wrapf(err, "failed to create new env")
	}
	if err := runCommand(command, newEnv, outStream, errStream); err != nil {
		return errors.Wrapf(err, "failed to run command %+s", command)
	}

	return nil
}

func makeNewEnv(origEnv []string, json map[string]string, envNames []string) ([]string, error) {
	newEnv := []string{}
	newJSON := make(map[string]string)
	for _, name := range envNames {
		value, ok := json[name]
		if !ok {
			return newEnv, errors.New(fmt.Sprintf("%s is not exists in json", name))
		}
		newJSON[name] = value
	}
	for key, value := range newJSON {
		newEnv = append(newEnv, fmt.Sprintf("%s=%s", key, value))
	}
	// if same key exists origEnv and newEnv, json parameter overwrite the origEnv
	newEnv = append(origEnv, newEnv...)
	return newEnv, nil
}

func printVersion(out io.Writer) error {
	_, err := fmt.Fprintf(out, "%s v%s (rev:%s)\n", cmdName, version, revision)
	return err
}

func runCommand(command, envVars []string, outStream, errStream io.Writer) error {
	bin, err := exec.LookPath(command[0])
	if err != nil {
		return err
	}
	cmd := exec.Command(bin, command[1:]...)
	cmd.Env = envVars
	cmd.Stderr = errStream
	cmd.Stdout = outStream
	return cmd.Run()
}
