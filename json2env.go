package json2env

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
)

const cmdName = "json2env"

// Run the json2env
func Run(ctx context.Context, argv []string, outStream, errStream io.Writer, inStream io.Reader) error {
	log.SetOutput(errStream)
	fs := flag.NewFlagSet(
		fmt.Sprintf("%s (v%s rev:%s)", cmdName, version, revision), flag.ContinueOnError)
	fs.SetOutput(errStream)
	ver := fs.Bool("version", false, "display version")
	// write flag setting

	if err := fs.Parse(argv); err != nil {
		return err
	}
	if *ver {
		return printVersion(outStream)
	}
	var envJSON map[string]string
	err := json.NewDecoder(inStream).Decode(&envJSON)
	if err != nil {
		return err
	}
	fmt.Fprint(outStream, envJSON)

	return nil
}

func makeNewEnv(origEnv []string, json map[string]string) []string {
	newEnv := []string{}
	for key, value := range json {
		newEnv = append(newEnv, fmt.Sprintf("%s=%s", key, value))
	}
	// if same key exists origEnv and newEnv, json parameter overwrite the origEnv
	newEnv = append(origEnv, newEnv...)
	return newEnv
}

func printVersion(out io.Writer) error {
	_, err := fmt.Fprintf(out, "%s v%s (rev:%s)\n", cmdName, version, revision)
	return err
}
