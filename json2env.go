package json2env2

import (
	"bufio"
	"context"
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
	// write implementation
	stdin := bufio.NewScanner(inStream)
	input := ""
	for stdin.Scan() {
		input += stdin.Text()
	}
	fmt.Fprint(outStream, input)

	return nil
}

func printVersion(out io.Writer) error {
	_, err := fmt.Fprintf(out, "%s v%s (rev:%s)\n", cmdName, version, revision)
	return err
}
