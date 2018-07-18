package jcon

import (
	"encoding/json"
	"strings"

	"github.com/infobloxopen/themis/jparser"
	"github.com/infobloxopen/themis/pdp"
)

func unmarshalCommand(d *json.Decoder, s pdp.Symbols, u *pdp.ContentUpdate) error {
	var op int
	opOk := false

	var path []string
	pathOk := false

	var entity *pdp.ContentItem
	entityOk := false

	err := jparser.UnmarshalObject(d, func(k string, d *json.Decoder) error {
		switch strings.ToLower(k) {
		case "op":
			if opOk {
				return newDuplicateCommandFieldError(k)
			}

			s, err := jparser.GetString(d, "operation")
			if err != nil {
				return err
			}

			op, opOk = pdp.UpdateOpIDs[strings.ToLower(s)]
			if !opOk {
				return newUnknownContentUpdateOperationError(s)
			}

			return nil

		case "path":
			if pathOk {
				return newDuplicateCommandFieldError(k)
			}
			path = []string{}
			err := jparser.GetStringSequence(d, func(idx int, s string) error {
				path = append(path, s)
				return nil
			}, "path")
			if err != nil {
				return err
			}

			pathOk = true
			return nil

		case "entity":
			if entityOk {
				return newDuplicateCommandFieldError(k)
			}

			var err error
			entity, err = unmarshalContentItem("", s, d)
			if err != nil {
				return err
			}

			entityOk = true
			return nil
		}

		return newUnknownCommadFieldError(k)
	}, "command")
	if err != nil {
		return err
	}

	if !opOk {
		return newMissingCommandOpError()
	}

	if !pathOk {
		return newMissingCommandPathError()
	}

	if op == pdp.UOAdd && !entityOk {
		return newMissingCommandEntityError()
	}

	u.Append(op, path, entity)
	return nil
}

func unmarshalCommands(d *json.Decoder, s pdp.Symbols, u *pdp.ContentUpdate) error {
	ok, err := jparser.CheckRootArrayStart(d)
	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	err = jparser.UnmarshalObjectArray(d, func(idx int, d *json.Decoder) error {
		if err := unmarshalCommand(d, s, u); err != nil {
			return bindErrorf(err, "%d", idx)
		}
		return nil
	}, "update")

	if err != nil {
		return err
	}

	return jparser.CheckEOF(d)
}
