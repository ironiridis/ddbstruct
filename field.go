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
	val := reflect.ValueOf(d).Elem().Field(f.idx)
	if val.IsZero() {
		if f.optional { // skip zero attribute
			return nil
		}
		if f.defvalue != "" { // apply default value
			val = reflect.ValueOf(f.defvalue)
		}
	}
	m[f.name] = f.enc(val)
	return nil
}
