package ddbstruct

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type encodeFunc func(reflect.Value) types.AttributeValue
type decodeFunc func(reflect.Value, types.AttributeValue)

type fieldTagData struct {
	pk        bool
	sk        bool
	fieldname string
	structidx int
	defvalue  string // only valid for string typed fields
	optional  bool
	enctype   string
	gotype    reflect.Type
	enc       encodeFunc
	dec       decodeFunc
}

var typeTime = reflect.TypeOf(time.Time{})
var typeDuration = reflect.TypeOf(time.Duration(0))
var typeBytes = reflect.TypeOf([]byte{})

func imp(t reflect.Type, d interface{}) bool {
	return t.Implements(reflect.TypeOf(d).Elem())
}

func isBinEncoder(t reflect.Type) bool {
	return imp(t, new(encoding.BinaryMarshaler)) && imp(t, new(encoding.BinaryUnmarshaler))
}

func isTextEncoder(t reflect.Type) bool {
	return imp(t, new(encoding.TextMarshaler)) && imp(t, new(encoding.TextUnmarshaler))
}

func isJSONEncoder(t reflect.Type) bool {
	return imp(t, new(json.Marshaler)) && imp(t, new(json.Unmarshaler))
}

func (ftd *fieldTagData) D(d interface{}, av types.AttributeValue) (err error) {
	if ftd.dec == nil {
		return fmt.Errorf("no decode function available for field %q of %T", ftd.fieldname, d)
	}
	defer func() {
		if panicVal := recover(); panicVal != nil {
			if panicErr, ok := panicVal.(error); ok {
				err = fmt.Errorf("failed to decode %q of %T from %+v: %w", ftd.fieldname, d, av, panicErr)
			} else {
				err = fmt.Errorf("failed to decode %q of %T from %+v: %v", ftd.fieldname, d, av, panicVal)
			}
		}
	}()
	ref := reflect.ValueOf(d).Elem().Field(ftd.structidx)
	if !ref.CanSet() {
		err = fmt.Errorf("cannot set field %q of %T", ftd.fieldname, d)
		return
	}
	ftd.dec(ref, av)
	return nil
}

func (ftd *fieldTagData) guessScalerCodec() bool {
	switch ftd.gotype.Kind() {
	case reflect.String:
		ftd.enc = func(d reflect.Value) types.AttributeValue {
			return &types.AttributeValueMemberS{Value: d.String()}
		}
		ftd.dec = func(d reflect.Value, av types.AttributeValue) {
			d.SetString(av.(*types.AttributeValueMemberS).Value)
		}
		return true
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
		ftd.enc = func(d reflect.Value) types.AttributeValue {
			return &types.AttributeValueMemberN{Value: strconv.FormatInt(d.Int(), 10)}
		}
		ftd.dec = func(d reflect.Value, av types.AttributeValue) {
			sv := av.(*types.AttributeValueMemberN).Value
			i, err := strconv.ParseInt(sv, 10, 0)
			if err != nil {
				panic(fmt.Errorf("cannot convert %q to %s: %w", sv, d.Kind(), err))
			}
			if d.OverflowInt(i) {
				panic(fmt.Errorf("value %q overflows %s", sv, d.Kind()))
			}
			d.SetInt(i)
		}
		return true
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
		ftd.enc = func(d reflect.Value) types.AttributeValue {
			return &types.AttributeValueMemberN{Value: strconv.FormatUint(d.Uint(), 10)}
		}
		ftd.dec = func(d reflect.Value, av types.AttributeValue) {
			sv := av.(*types.AttributeValueMemberN).Value
			i, err := strconv.ParseUint(sv, 10, 0)
			if err != nil {
				panic(fmt.Errorf("cannot convert %q to %s: %w", sv, d.Kind(), err))
			}
			if d.OverflowUint(i) {
				panic(fmt.Errorf("value %q overflows %s", sv, d.Kind()))
			}
			d.SetUint(i)
		}
		return true
	case reflect.Float32:
		ftd.enc = func(d reflect.Value) types.AttributeValue {
			return &types.AttributeValueMemberN{Value: strconv.FormatFloat(d.Float(), 'G', -1, 32)}
		}
		ftd.dec = func(d reflect.Value, av types.AttributeValue) {
			sv := av.(*types.AttributeValueMemberN).Value
			f, err := strconv.ParseFloat(sv, 32)
			if err != nil {
				panic(fmt.Errorf("cannot convert %q to %s: %w", sv, d.Kind(), err))
			}
			if d.OverflowFloat(f) {
				panic(fmt.Errorf("value %q overflows %s", sv, d.Kind()))
			}
			d.SetFloat(f)
		}
		return true
	case reflect.Float64:
		ftd.enc = func(d reflect.Value) types.AttributeValue {
			return &types.AttributeValueMemberN{Value: strconv.FormatFloat(d.Float(), 'G', -1, 64)}
		}
		ftd.dec = func(d reflect.Value, av types.AttributeValue) {
			sv := av.(*types.AttributeValueMemberN).Value
			f, err := strconv.ParseFloat(sv, 64)
			if err != nil {
				panic(fmt.Errorf("cannot convert %q to %s: %w", sv, d.Kind(), err))
			}
			if d.OverflowFloat(f) {
				panic(fmt.Errorf("value %q overflows %s", sv, d.Kind()))
			}
			d.SetFloat(f)
		}
		return true
	case reflect.Bool:
		ftd.enc = func(d reflect.Value) types.AttributeValue {
			return &types.AttributeValueMemberBOOL{Value: d.Bool()}
		}
		ftd.dec = func(d reflect.Value, av types.AttributeValue) {
			d.SetBool(av.(*types.AttributeValueMemberBOOL).Value)
		}
		return true
	}
	switch ftd.gotype { // not scalar but pretty close
	case typeBytes:
		ftd.enc = func(d reflect.Value) types.AttributeValue {
			cp := make([]byte, d.Len())
			copy(cp, d.Bytes())
			return &types.AttributeValueMemberB{Value: cp}
		}
		ftd.dec = func(d reflect.Value, av types.AttributeValue) {
			bv := av.(*types.AttributeValueMemberB).Value
			cp := make([]byte, len(bv))
			copy(cp, bv)
			d.Set(reflect.ValueOf(cp))
		}
		return true

	}
	return false
}

func (ftd *fieldTagData) guessCommonCodec() bool {
	switch ftd.gotype {
	case typeDuration:
		ftd.enc = func(d reflect.Value) types.AttributeValue {
			return &types.AttributeValueMemberS{Value: d.Interface().(time.Duration).String()}
		}
		ftd.dec = func(d reflect.Value, av types.AttributeValue) {
			t, err := time.ParseDuration(av.(*types.AttributeValueMemberS).Value)
			if err != nil {
				panic(err)
			}
			d.Set(reflect.ValueOf(t))
		}
		return true
	}
	return false
}

func (ftd *fieldTagData) guessImplementsCodec() bool {
	switch {
	case isTextEncoder(ftd.gotype):
		ftd.enc, ftd.dec = encText, decText
	case isJSONEncoder(ftd.gotype):
		ftd.enc, ftd.dec = encJSON, decJSON
	case isBinEncoder(ftd.gotype):
		ftd.enc, ftd.dec = encBinary, decBinary
	default:
		return false // didn't match anything
	}
	return true
}

func (ftd *fieldTagData) typecalc() error {
	if ftd.pk && ftd.sk {
		return fmt.Errorf("pk and sk both set on %q", ftd.fieldname)
	}
	if ftd.enctype == "" { // with no explicit type, let's start by guessing
		if ftd.guessScalerCodec() { // matches a scalar type
			return nil
		}
		if ftd.guessCommonCodec() { // matches a common type
			return nil
		}
		if ftd.guessImplementsCodec() { // matches a known codec
			return nil
		}
		return fmt.Errorf("unable to guess encoding for %q field of type %s", ftd.fieldname, ftd.gotype)
	}
	switch ftd.enctype {
	case "string":
		if ftd.gotype.Kind() == reflect.String {
			ftd.enc = func(d reflect.Value) types.AttributeValue {
				return &types.AttributeValueMemberS{Value: d.String()}
			}
			ftd.dec = func(d reflect.Value, av types.AttributeValue) {
				d.SetString(av.(*types.AttributeValueMemberS).Value)
			}
			return nil
		}
		if isTextEncoder(ftd.gotype) {
			ftd.enc, ftd.dec = encText, decText
			return nil
		}
		return fmt.Errorf("field %q cannot be typed as string automatically", ftd.fieldname)
	case "binary":
		if ftd.gotype == typeBytes {
			ftd.enc = func(d reflect.Value) types.AttributeValue {
				return &types.AttributeValueMemberS{Value: d.String()}
			}
			ftd.dec = func(d reflect.Value, av types.AttributeValue) {
				d.SetString(av.(*types.AttributeValueMemberS).Value)
			}
			return nil
		}
		if isBinEncoder(ftd.gotype) {
			ftd.enc, ftd.dec = encBinary, decBinary
			return nil
		}
		return fmt.Errorf("field %q cannot be typed as binary automatically", ftd.fieldname)
	case "json":
		if isJSONEncoder(ftd.gotype) {
			ftd.enc, ftd.dec = encJSON, decJSON
			return nil
		}
		ftd.enc = func(d reflect.Value) types.AttributeValue {
			sv, err := json.Marshal(d.Interface())
			if err != nil {
				panic(err)
			}
			return &types.AttributeValueMemberS{Value: string(sv)}
		}
		ftd.dec = func(d reflect.Value, av types.AttributeValue) {
			err := json.Unmarshal([]byte(av.(*types.AttributeValueMemberS).Value), d.Interface())
			if err != nil {
				panic(err)
			}
		}
		return nil
	case "nanoseconds":
		switch ftd.gotype {
		case typeDuration:
			ftd.enc = func(d reflect.Value) types.AttributeValue {
				ns := d.Interface().(time.Duration).Nanoseconds()
				return &types.AttributeValueMemberN{Value: strconv.FormatInt(ns, 10)}
			}
			ftd.dec = func(d reflect.Value, av types.AttributeValue) {
				nv, err := strconv.ParseInt(av.(*types.AttributeValueMemberN).Value, 10, 64)
				if err != nil {
					panic(err)
				}
				d.SetInt(int64(time.Nanosecond) * nv)
			}
			return nil
		case typeTime:
			ftd.enc = func(d reflect.Value) types.AttributeValue {
				t := d.Interface().(time.Time)
				return &types.AttributeValueMemberN{Value: strconv.FormatInt(t.UnixNano(), 10)}
			}
			ftd.dec = func(d reflect.Value, av types.AttributeValue) {
				nv, err := strconv.ParseInt(av.(*types.AttributeValueMemberN).Value, 10, 64)
				if err != nil {
					panic(err)
				}
				d.Set(reflect.ValueOf(time.Unix(0, nv)))
			}
			return nil
		}
		return fmt.Errorf("cannot encode field %q (a %s) as %q", ftd.fieldname, ftd.gotype, ftd.enctype)
	case "epoch", "seconds":
		switch ftd.gotype {
		case typeDuration:
			ftd.enc = func(d reflect.Value) types.AttributeValue {
				s := d.Interface().(time.Duration).Seconds()
				return &types.AttributeValueMemberN{Value: strconv.FormatFloat(s, 'G', 64, 0)}
			}
			ftd.dec = func(d reflect.Value, av types.AttributeValue) {
				nv, err := strconv.ParseInt(av.(*types.AttributeValueMemberN).Value, 10, 64)
				if err != nil {
					panic(err)
				}
				d.SetInt(int64(time.Nanosecond) * nv)
			}
			return nil
		case typeTime:
			ftd.enc = func(d reflect.Value) types.AttributeValue {
				t := d.Interface().(time.Time).Unix()
				return &types.AttributeValueMemberN{Value: strconv.FormatInt(t, 10)}
			}
			ftd.dec = func(d reflect.Value, av types.AttributeValue) {
				nv, err := strconv.ParseInt(av.(*types.AttributeValueMemberN).Value, 10, 64)
				if err != nil {
					panic(err)
				}
				d.Set(reflect.ValueOf(time.Unix(0, nv)))
			}
			return nil
		}
		return fmt.Errorf("cannot encode field %q (a %s) as %q", ftd.fieldname, ftd.gotype, ftd.enctype)
	}
	return fmt.Errorf("cannot encode field %q (a %s) with unknown enctype %q", ftd.fieldname, ftd.gotype, ftd.enctype)
}
