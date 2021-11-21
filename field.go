package ddbstruct

import (
	"fmt"
	"reflect"
)

type field struct {
	pk       bool
	sk       bool
	name     string
	idx      int
	defvalue string // only valid for string typed fields
	optional bool
	enctype  string
	gotype   reflect.Type
	enc      encodeFunc
	dec      decodeFunc
}

func (f *field) appendAV(m avmap, d interface{}) error {
	if f.enc == nil {
		return fmt.Errorf("no encode function available for field %q of %T", f.name, d)
	}
	if reflect.ValueOf(d).Elem().Field(f.idx).IsZero() {
		if f.optional { // skip zero attribute
			return nil
		}
		if f.defvalue != "" { // apply default value
			m[f.name] = f.enc(&struct{ S string }{S: f.defvalue}, 0)
			return nil
		}
	}
	m[f.name] = f.enc(d, f.idx)
	return nil
}
