// Package test implements test command for PEPCLI.
package test

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/infobloxopen/themis/pdp"
	"github.com/infobloxopen/themis/pep"

	"github.com/infobloxopen/themis/pepcli/requests"
)

const (
	// Name contains title of function implemented by the package.
	Name = "test"
	// Description provides additional information on the package functionality.
	Description = "evaluates given requests on PDP server"
)

// Exec tests requests from input with given pdp server and dumps responses in YAML format
// to given file or standard output if file name is empty.
func Exec(addr string, opts []pep.Option, maxRequestSize, maxResponseObligations uint32, in, out string, n int, v interface{}) error {
	reqs, err := requests.Load(in, maxRequestSize)
	if err != nil {
		return fmt.Errorf("can't load requests from \"%s\": %s", in, err)
	}

	if n < 1 {
		n = len(reqs)
	}

	f := os.Stdout
	if len(out) > 0 {
		f, err = os.Create(out)
		if err != nil {
			return err
		}
		defer f.Close()
	}

	c := pep.NewClient(opts...)
	err = c.Connect(addr)
	if err != nil {
		return fmt.Errorf("can't connect to %s: %s", addr, err)
	}
	defer c.Close()

	obligations := make([]pdp.AttributeAssignment, maxResponseObligations)
	res := pdp.Response{}
	for i := 0; i < n; i++ {
		idx := i % len(reqs)
		req := reqs[idx]

		res.Obligations = obligations
		err := c.Validate(req, &res)
		if err != nil {
			return fmt.Errorf("can't send request %d (%d): %s", idx, i, err)
		}

		err = dump(res, f)
		if err != nil {
			return fmt.Errorf("can't dump response for reqiest %d (%d): %s", idx, i, err)
		}
	}

	return nil
}

func dump(r pdp.Response, f io.Writer) error {
	lines := []string{fmt.Sprintf("- effect: %s", pdp.EffectNameFromEnum(r.Effect))}
	if r.Status != nil {
		lines = append(lines, fmt.Sprintf("  reason: %q", r.Status))
	}

	if len(r.Obligations) > 0 {
		lines = append(lines, "  obligation:")
		for i, o := range r.Obligations {
			id, t, v, err := o.Serialize(nil)
			if err != nil {
				return fmt.Errorf("can't get %d obligation: %s", i+1, err)
			}

			lines = append(lines, fmt.Sprintf("    - id: %q", id))
			lines = append(lines, fmt.Sprintf("      type: %q", t))
			lines = append(lines, fmt.Sprintf("      value: %q", v))
			lines = append(lines, "")
		}
	} else {
		lines = append(lines, "")
	}

	_, err := fmt.Fprintf(f, "%s\n", strings.Join(lines, "\n"))
	return err
}
