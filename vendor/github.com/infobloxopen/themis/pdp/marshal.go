package pdp

import (
	"encoding/json"
	"io"
)

// StorageMarshal interface defines functions
// to capturing storage state information
type StorageMarshal interface {
	GetID() (id string, hidden bool)
	MarshalWithDepth(out io.Writer, depth int) error
}

// PolicySet/Policy representation for marshaling
type storageEvalFmt struct {
	Ord         int                   `json:"ord"`
	ID          string                `json:"id"`
	Target      Target                `json:"target"`
	Obligations []AttributeAssignment `json:"obligations"`
	Algorithm   json.Marshaler        `json:"algorithm"`
}

// Rule representation for marshaling
type storageRuleFmt struct {
	Ord         int                   `json:"ord"`
	ID          string                `json:"id"`
	Target      Target                `json:"target"`
	Obligations []AttributeAssignment `json:"obligations"`
	Effect      string                `json:"effect"`
}

func marshalHeader(v interface{}, out io.Writer) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	n := len(b)
	if n < 1 || b[n-1] != '}' {
		return newInvalidHeaderError(v)
	}
	_, err = out.Write(b[:n-1])
	return err
}

// Algorithm representation for marshalling
type algFmt struct {
	Type string `json:"type"`
}

type mapperAlgFmt struct {
	Type      string         `json:"type"`
	Default   string         `json:"def"`
	Error     string         `json:"err"`
	Algorithm json.Marshaler `json:"alg"`
}
