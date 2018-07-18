package perf

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

type config struct {
	parallel int
	limit    int64
}

var perfFlagSet = flag.NewFlagSet(Name, flag.ExitOnError)

// FlagsParser implements parsing for options specific to the package.
func FlagsParser(args []string) interface{} {
	conf := config{}

	perfFlagSet.Usage = usage
	perfFlagSet.IntVar(&conf.parallel, "p", 0, "make given number of requests in parallel\n\t"+
		"(default and zero - make requests sequentially;\n\t"+
		" negative - make all requess in parallel)")
	perfFlagSet.Int64Var(&conf.limit, "l", 0, "limit request rate by adding 1s/limit pauses\n\t"+
		"(default and less than one - no limit;\n\t"+
		" shouldn't be more than 1,000,000,000)")
	perfFlagSet.Parse(args)

	count := perfFlagSet.NArg()
	if count > 1 {
		tail := strings.Join(perfFlagSet.Args()[1:count], "\", \"")
		fmt.Fprintf(os.Stderr, "trailing arguments after cluster name: \"%s\"\n", tail)
		usage()
		os.Exit(2)
	}

	if conf.limit > 1000000000 {
		fmt.Fprintf(os.Stderr, "request rate limit is too high %d > 1,000,000,000", conf.limit)
		usage()
		os.Exit(2)
	}

	return conf
}

func usage() {
	base := path.Base(os.Args[0])
	fmt.Fprintf(os.Stderr,
		"Usage of %s.%s:\n\n"+
			"  %s [GLOBAL OPTIONS] %s [%s OPTIONS]\n\n"+
			"GLOBAL OPTIONS:\n"+
			"  See %s -h\n\n"+
			"%s OPTIONS:\n", base, Name, base, Name, strings.ToUpper(Name), base, strings.ToUpper(Name))
	perfFlagSet.PrintDefaults()
}
