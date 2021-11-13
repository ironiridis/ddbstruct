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

func encBinary(d reflect.Value) types.AttributeValue {
	res := d.MethodByName("MarshalBinary").Call([]reflect.Value{})
	if !res[1].IsNil() {
		panic(res[1].Interface().(error))
	}
	return &types.AttributeValueMemberB{Value: res[0].Bytes()}
}

func encText(d reflect.Value) types.AttributeValue {
	res := d.MethodByName("MarshalText").Call([]reflect.Value{})
	if !res[1].IsNil() {
		panic(res[1].Interface().(error))
	}
	return &types.AttributeValueMemberS{Value: string(res[0].Bytes())}
}

func encJSON(d reflect.Value) types.AttributeValue {
	res := d.MethodByName("MarshalJSON").Call([]reflect.Value{})
	if !res[1].IsNil() {
		panic(res[1].Interface().(error))
	}
	return &types.AttributeValueMemberS{Value: string(res[0].Bytes())}
}

func decBinary(d reflect.Value, av types.AttributeValue) {
	bv := av.(*types.AttributeValueMemberB).Value
	res := d.MethodByName("UnmarshalBinary").Call([]reflect.Value{reflect.ValueOf(bv)})
	if !res[0].IsNil() {
		panic(res[0].Interface().(error))
	}
}

func decText(d reflect.Value, av types.AttributeValue) {
	sv := av.(*types.AttributeValueMemberS).Value
	bv := make([]byte, len(sv))
	copy(bv, sv)
	res := d.MethodByName("UnmarshalText").Call([]reflect.Value{reflect.ValueOf(bv)})
	if !res[0].IsNil() {
		panic(res[0].Interface().(error))
	}
}

func decJSON(d reflect.Value, av types.AttributeValue) {
	sv := av.(*types.AttributeValueMemberS).Value
	bv := make([]byte, len(sv))
	copy(bv, sv)
	res := d.MethodByName("UnmarshalJSON").Call([]reflect.Value{reflect.ValueOf(bv)})
	if !res[0].IsNil() {
		panic(res[0].Interface().(error))
	}
}

func (f *field) guessScalerCodec() bool {
	switch f.gotype.Kind() {
	case reflect.String:
		f.enc = func(d reflect.Value) types.AttributeValue {
			return &types.AttributeValueMemberS{Value: d.String()}
		}
		f.dec = func(d reflect.Value, av types.AttributeValue) {
			d.SetString(av.(*types.AttributeValueMemberS).Value)
		}
		return true
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
		f.enc = func(d reflect.Value) types.AttributeValue {
			return &types.AttributeValueMemberN{Value: strconv.FormatInt(d.Int(), 10)}
		}
		f.dec = func(d reflect.Value, av types.AttributeValue) {
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
		f.enc = func(d reflect.Value) types.AttributeValue {
			return &types.AttributeValueMemberN{Value: strconv.FormatUint(d.Uint(), 10)}
		}
		f.dec = func(d reflect.Value, av types.AttributeValue) {
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
		f.enc = func(d reflect.Value) types.AttributeValue {
			return &types.AttributeValueMemberN{Value: strconv.FormatFloat(d.Float(), 'G', -1, 32)}
		}
		f.dec = func(d reflect.Value, av types.AttributeValue) {
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
		f.enc = func(d reflect.Value) types.AttributeValue {
			return &types.AttributeValueMemberN{Value: strconv.FormatFloat(d.Float(), 'G', -1, 64)}
		}
		f.dec = func(d reflect.Value, av types.AttributeValue) {
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
		f.enc = func(d reflect.Value) types.AttributeValue {
			return &types.AttributeValueMemberBOOL{Value: d.Bool()}
		}
		f.dec = func(d reflect.Value, av types.AttributeValue) {
			d.SetBool(av.(*types.AttributeValueMemberBOOL).Value)
		}
		return true
	}
	switch f.gotype { // not scalar but pretty close
	case typeBytes:
		f.enc = func(d reflect.Value) types.AttributeValue {
			cp := make([]byte, d.Len())
			copy(cp, d.Bytes())
			return &types.AttributeValueMemberB{Value: cp}
		}
		f.dec = func(d reflect.Value, av types.AttributeValue) {
			bv := av.(*types.AttributeValueMemberB).Value
			cp := make([]byte, len(bv))
			copy(cp, bv)
			d.Set(reflect.ValueOf(cp))
		}
		return true

	}
	return false
}

func (f *field) guessCommonCodec() bool {
	switch f.gotype {
	case typeDuration:
		f.enc = func(d reflect.Value) types.AttributeValue {
			return &types.AttributeValueMemberS{Value: d.Interface().(time.Duration).String()}
		}
		f.dec = func(d reflect.Value, av types.AttributeValue) {
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

func (f *field) guessImplementsCodec() bool {
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
	if f.pk && f.sk {
		return fmt.Errorf("pk and sk both set on %q", f.name)
	}
	if f.enctype == "" { // with no explicit type, let's start by guessing
		if f.guessScalerCodec() { // matches a scalar type
			return nil
		}
		if f.guessCommonCodec() { // matches a common type
			return nil
		}
		if f.guessImplementsCodec() { // matches a known codec
			return nil
		}
		return fmt.Errorf("unable to guess encoding for %q field of type %s", f.name, f.gotype)
	}
	switch f.enctype {
	case "string":
		if f.gotype.Kind() == reflect.String {
			f.enc = func(d reflect.Value) types.AttributeValue {
				return &types.AttributeValueMemberS{Value: d.String()}
			}
			f.dec = func(d reflect.Value, av types.AttributeValue) {
				d.SetString(av.(*types.AttributeValueMemberS).Value)
			}
			return nil
		}
		if isTextEncoder(f.gotype) {
			f.enc, f.dec = encText, decText
			return nil
		}
		return fmt.Errorf("field %q cannot be typed as string automatically", f.name)
	case "binary":
		if f.gotype == typeBytes {
			f.enc = func(d reflect.Value) types.AttributeValue {
				return &types.AttributeValueMemberS{Value: d.String()}
			}
			f.dec = func(d reflect.Value, av types.AttributeValue) {
				d.SetString(av.(*types.AttributeValueMemberS).Value)
			}
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
		f.enc = func(d reflect.Value) types.AttributeValue {
			sv, err := json.Marshal(d.Interface())
			if err != nil {
				panic(err)
			}
			return &types.AttributeValueMemberS{Value: string(sv)}
		}
		f.dec = func(d reflect.Value, av types.AttributeValue) {
			err := json.Unmarshal([]byte(av.(*types.AttributeValueMemberS).Value), d.Interface())
			if err != nil {
				panic(err)
			}
		}
		return nil
	case "nanoseconds":
		switch f.gotype {
		case typeDuration:
			f.enc = func(d reflect.Value) types.AttributeValue {
				ns := d.Interface().(time.Duration).Nanoseconds()
				return &types.AttributeValueMemberN{Value: strconv.FormatInt(ns, 10)}
			}
			f.dec = func(d reflect.Value, av types.AttributeValue) {
				nv, err := strconv.ParseInt(av.(*types.AttributeValueMemberN).Value, 10, 64)
				if err != nil {
					panic(err)
				}
				d.SetInt(int64(time.Nanosecond) * nv)
			}
			return nil
		case typeTime:
			f.enc = func(d reflect.Value) types.AttributeValue {
				t := d.Interface().(time.Time)
				return &types.AttributeValueMemberN{Value: strconv.FormatInt(t.UnixNano(), 10)}
			}
			f.dec = func(d reflect.Value, av types.AttributeValue) {
				nv, err := strconv.ParseInt(av.(*types.AttributeValueMemberN).Value, 10, 64)
				if err != nil {
					panic(err)
				}
				d.Set(reflect.ValueOf(time.Unix(0, nv)))
			}
			return nil
		}
		return fmt.Errorf("cannot encode field %q (a %s) as %q", f.name, f.gotype, f.enctype)
	case "epoch", "seconds":
		switch f.gotype {
		case typeDuration:
			f.enc = func(d reflect.Value) types.AttributeValue {
				s := d.Interface().(time.Duration).Seconds()
				return &types.AttributeValueMemberN{Value: strconv.FormatFloat(s, 'G', 64, 0)}
			}
			f.dec = func(d reflect.Value, av types.AttributeValue) {
				nv, err := strconv.ParseInt(av.(*types.AttributeValueMemberN).Value, 10, 64)
				if err != nil {
					panic(err)
				}
				d.SetInt(int64(time.Second) * nv)
			}
			return nil
		case typeTime:
			f.enc = func(d reflect.Value) types.AttributeValue {
				t := d.Interface().(time.Time).Unix()
				return &types.AttributeValueMemberN{Value: strconv.FormatInt(t, 10)}
			}
			f.dec = func(d reflect.Value, av types.AttributeValue) {
				nv, err := strconv.ParseInt(av.(*types.AttributeValueMemberN).Value, 10, 64)
				if err != nil {
					panic(err)
				}
				d.Set(reflect.ValueOf(time.Unix(nv, 0)))
			}
			return nil
		}
		return fmt.Errorf("cannot encode field %q (a %s) as %q", f.name, f.gotype, f.enctype)
	}
	return fmt.Errorf("cannot encode field %q (a %s) with unknown enctype %q", f.name, f.gotype, f.enctype)
}
