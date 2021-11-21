package ddbstruct

import (
	"fmt"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type encodeFunc func(interface{}, int) types.AttributeValue
type decodeFunc func(interface{}, int, types.AttributeValue)

var typeTime = reflect.TypeOf(time.Time{})
var typeDuration = reflect.TypeOf(time.Duration(0))
var typeBytes = reflect.TypeOf([]byte{})

func (f *field) tryBasicMarshaling() bool {
	switch f.gotype.Kind() {
	case reflect.String:
		f.enc, f.dec = encString, decString
		return true
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
		f.enc, f.dec = encInt, decInt
		return true
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
		f.enc, f.dec = encUint, decUint
		return true
	case reflect.Float32:
		f.enc, f.dec = encFloat(32), decFloat(32)
		return true
	case reflect.Float64:
		f.enc, f.dec = encFloat(64), decFloat(64)
		return true
	case reflect.Bool:
		f.enc, f.dec = encBool, decBool
		return true
	}
	switch f.gotype {
	case typeBytes:
		f.enc, f.dec = encBytes, decBytes
		return true
	case typeDuration:
		f.enc, f.dec = encDurationString, decDurationString
		return true
	}
	return false
}

func (f *field) tryInterfaceMarshaling() bool {
	switch {
	case isTextEncoder(f.gotype):
		f.enc, f.dec = encText, decText
	case isJSONEncoder(f.gotype):
		f.enc, f.dec = encJSON, decJSON
	case isBinEncoder(f.gotype):
		f.enc, f.dec = encBinary, decBinary
	default:
		return false // didn't match anything
	}
	return true
}

func (f *field) typecalc() error {
	if f.enctype == "" { // with no explicit type, let's start by guessing
		if f.tryBasicMarshaling() { // matches basic types
			return nil
		}
		if f.tryInterfaceMarshaling() { // uses standard interfaces
			return nil
		}
		return fmt.Errorf("unable to guess encoding for %q field of type %s", f.name, f.gotype)
	}
	switch f.enctype {
	case "string":
		if f.gotype.Kind() == reflect.String {
			f.enc, f.dec = encString, decString
			return nil
		}
		if isTextEncoder(f.gotype) {
			f.enc, f.dec = encText, decText
			return nil
		}
		return fmt.Errorf("field %q cannot be typed as string automatically", f.name)
	case "binary", "bytes":
		if f.gotype == typeBytes {
			f.enc, f.dec = encBytes, decBytes
			return nil
		}
		if isBinEncoder(f.gotype) {
			f.enc, f.dec = encBinary, decBinary
			return nil
		}
		return fmt.Errorf("field %q cannot be typed as binary automatically", f.name)
	case "json":
		if isJSONEncoder(f.gotype) {
			f.enc, f.dec = encJSON, decJSON
			return nil
		}
		f.enc, f.dec = encJSONRaw, decJSONRaw
		return nil
	case "nano", "nanoseconds":
		switch f.gotype {
		case typeDuration:
			f.enc, f.dec = encDurationNano, decDurationNano
			return nil
		case typeTime:
			f.enc, f.dec = encTimeNano, decTimeNano
			return nil
		}
	case "epoch", "seconds":
		switch f.gotype {
		case typeDuration:
			f.enc, f.dec = encDurationSec, decDurationSec
			return nil
		case typeTime:
			f.enc, f.dec = encTimeEpoch, decTimeEpoch
			return nil
		}
	}
	return fmt.Errorf("cannot encode field %q (a %s) as %q", f.name, f.gotype, f.enctype)
}
