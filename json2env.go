package json2env2

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
)

const cmdName = "json2env2"

// Run the json2env2
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
	var envJSON interface{}
	err := json.NewDecoder(inStream).Decode(&envJSON)
	if err != nil {
		return err
	}
	fmt.Fprint(outStream, envJSON)

	return nil
}

func printVersion(out io.Writer) error {
	_, err := fmt.Fprintf(out, "%s v%s (rev:%s)\n", cmdName, version, revision)
	return err
}
