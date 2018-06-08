// Package jcon implements JSON content (JCON) parser.
package jcon

import (
	"encoding/json"
	"io"

	"github.com/google/uuid"

	"github.com/infobloxopen/themis/pdp"
)

// Unmarshal parses JSON content representation to PDP's internal represntation
// and returns pointer to LocalContent. It sets given tag to the content.
// Content with no tag can't be updated.
func Unmarshal(r io.Reader, tag *uuid.UUID) (*pdp.LocalContent, error) {
	c := &content{
		symbols: pdp.MakeSymbols(),
	}
	err := c.unmarshal(json.NewDecoder(r))
	if err != nil {
		return nil, err
	}

	return pdp.NewLocalContent(c.id, tag, c.symbols, c.items), nil
}

// UnmarshalUpdate parses JSON content update representation to PDP's internal
// represntation. It requires content id and oldTag to match content to update.
// Value of newTag is set to the content when update is applied.
func UnmarshalUpdate(r io.Reader, cID string, oldTag, newTag uuid.UUID, s pdp.Symbols) (*pdp.ContentUpdate, error) {
	u := pdp.NewContentUpdate(cID, oldTag, newTag)
	err := unmarshalCommands(json.NewDecoder(r), s, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}
