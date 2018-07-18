package test

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

type config struct {
}

var testFlagSet = flag.NewFlagSet(Name, flag.ExitOnError)

// FlagsParser implements parsing for options specific to the package.
func FlagsParser(args []string) interface{} {
	conf := config{}

	testFlagSet.Usage = usage
	testFlagSet.Parse(args)

	count := testFlagSet.NArg()
	if count > 1 {
		tail := strings.Join(testFlagSet.Args()[1:count], "\", \"")
		fmt.Fprintf(os.Stderr, "trailing arguments after cluster name: \"%s\"\n", tail)
		usage()
		os.Exit(2)
	}

	return conf
}

func usage() {
	base := path.Base(os.Args[0])
	fmt.Fprintf(os.Stderr,
		"Usage of %s.%s:\n\n"+
			"  %s [GLOBAL OPTIONS] %s\n\n"+
			"GLOBAL OPTIONS:\n"+
			"  See %s -h\n\n", base, Name, base, Name, base)
	testFlagSet.PrintDefaults()
}
