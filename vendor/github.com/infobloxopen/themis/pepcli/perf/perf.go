// Package perf implements perf command for PEPCLI.
package perf

import (
	"fmt"

	"github.com/infobloxopen/themis/pep"
	"github.com/infobloxopen/themis/pepcli/requests"
)

const (
	// Name contains title of function implemented by the package.
	Name = "perf"
	// Description provides additional information on the package functionality.
	Description = "measures performance of evaluation given requests on PDP server"
)

// Exec runs performance test for given server with requests from input and
// dumps timings in JSON format to given file or standard output if file name is empty.
func Exec(addr string, opts []pep.Option, maxRequestSize, maxResponseObligations uint32, in, out string, n int, v interface{}) error {
	reqs, err := requests.Load(in, maxRequestSize)
	if err != nil {
		return fmt.Errorf("can't load requests from \"%s\": %s", in, err)
	}

	if n < 1 {
		n = len(reqs)
	}

	c := pep.NewClient(opts...)
	err = c.Connect(addr)
	if err != nil {
		return fmt.Errorf("can't connect to %s: %s", addr, err)
	}
	defer c.Close()

	recs, err := measurement(c, n, v.(config).parallel, v.(config).limit, reqs, maxResponseObligations)
	if err != nil {
		return err
	}

	return dump(recs, out)
}
