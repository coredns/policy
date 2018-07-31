package policy

import (
	"fmt"
	"github.com/infobloxopen/themis/pdp"
	"github.com/infobloxopen/themis/pdp/ast"
	"github.com/infobloxopen/themis/pdp/jcon"
	"log"
	"os"
	"path/filepath"
)

type builtinClient struct {
	policyFile   string
	contentFiles []string

	parser ast.Parser

	p *pdp.PolicyStorage
	c *pdp.LocalContentStorage
}

func newBuiltinClient(policyFile string, contentFiles []string) *builtinClient {
	c := &builtinClient{
		policyFile:   policyFile,
		contentFiles: contentFiles,
	}
	return c
}

func (c *builtinClient) setPolicyParser() {
	if c.policyFile != "" {
		ext := filepath.Ext(c.policyFile)
		switch ext {
		case ".json":
			c.parser = ast.NewJSONParser()
		case ".yaml":
			c.parser = ast.NewYAMLParser()
		}
	}
}

func (c *builtinClient) loadPolicies() error {
	log.Printf("[INFO] Loading policy '%s'", c.policyFile)
	c.setPolicyParser()

	pf, err := os.Open(c.policyFile)
	if err != nil {
		log.Printf("[ERROR] Failed to open policy file: %s", err)
		return err
	}

	log.Printf("[INFO] Parsing policy '%s'", c.policyFile)
	p, err := c.parser.Unmarshal(pf, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to parse policy: %s", err)
		return err
	}

	c.p = p
	return nil
}

func (c *builtinClient) loadContents() error {
	log.Print("[INFO] Loading content")

	var items []*pdp.LocalContent
	for _, path := range c.contentFiles {
		err := func() error {
			log.Printf("[INFO] Opening content '%s'", path)
			f, err := os.Open(path)
			if err != nil {
				log.Printf("[ERROR] Failed to open content: %s", err)
				return err
			}

			defer f.Close()

			log.Printf("[INFO] Parsing content '%s'", path)
			item, err := jcon.Unmarshal(f, nil)
			if err != nil {
				log.Printf("[ERROR] Failed to parse content: %s", err)
				return err
			}

			items = append(items, item)
			return nil
		}()
		if err != nil {
			return err
		}
	}

	c.c = pdp.NewLocalContentStorage(items)
	return nil
}

func (c *builtinClient) Connect(addr string) error {
	if err := c.loadPolicies(); err != nil {
		return err
	}

	if err := c.loadContents(); err != nil {
		return err
	}

	return nil
}

func (c *builtinClient) Close() {
	c.p = nil
	c.c = nil
}

func (c *builtinClient) Validate(in, out interface{}) error {
	req, ok := in.([]pdp.AttributeAssignment)
	if !ok {
		return fmt.Errorf("unknown request type passed to Validate()")
	}
	res, ok := out.(*pdp.Response)
	if !ok {
		return fmt.Errorf("unknown response type passed to Validate()")
	}

	ctx, err := pdp.NewContext(c.c, len(req), func(i int) (string, pdp.AttributeValue, error) {
		id := req[i].GetID()
		v, err := req[i].GetValue()
		return id, v, err
	})
	if err != nil {
		return fmt.Errorf("error creating pdp context '%s'", err)
	}

	r := c.p.Root().Calculate(ctx)

	res.Effect = r.Effect
	res.Status = r.Status
	if res.Obligations == nil {
		res.Obligations = r.Obligations
	} else {
		obLen := len(r.Obligations)
		if obLen > len(res.Obligations) {
			return fmt.Errorf("result obligations too small %d < %d",
				len(res.Obligations), obLen)
		}

		copy(res.Obligations, r.Obligations)
		res.Obligations = res.Obligations[:obLen]
	}

	return nil
}
