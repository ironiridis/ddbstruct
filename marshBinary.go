package ddbstruct

import (
	"encoding"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var intfBinaryMarshaler = reflect.TypeOf(new(encoding.BinaryMarshaler)).Elem()
var intfBinaryUnmarshaler = reflect.TypeOf(new(encoding.BinaryUnmarshaler)).Elem()

func encBinary(s interface{}, f int) types.AttributeValue {
	fp := getF(s, f)
	enc, ok := fp.Interface().(encoding.BinaryMarshaler)
	if !ok {
		enc, ok = fp.Addr().Interface().(encoding.BinaryMarshaler)
		if !ok {
			panic("neither " + fp.Type().String() + " nor *" + fp.Type().String() + " implements MarshalBinary")
		}
	}

	buf, err := enc.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return &types.AttributeValueMemberB{Value: buf}
}

func decBinary(s interface{}, f int, av types.AttributeValue) {
	fp := getF(s, f)
	dec, ok := fp.Interface().(encoding.BinaryUnmarshaler)
	if !ok {
		dec, ok = fp.Addr().Interface().(encoding.BinaryUnmarshaler)
		if !ok {
			panic("neither " + fp.Type().String() + " nor *" + fp.Type().String() + " implements UnmarshalBinary")
		}
	}

	err := dec.UnmarshalBinary(av.(*types.AttributeValueMemberB).Value)
	if err != nil {
		panic(err)
	}
}

func isBinEncoder(t reflect.Type) bool {
	var e, d bool
	e = t.Implements(intfBinaryMarshaler)
	d = t.Implements(intfBinaryUnmarshaler)

	// if the type isn't already a pointer, check to see if a pointer to
	// the type satisfies the interface instead
	if t.Kind() != reflect.Pointer {
		e = e || reflect.PointerTo(t).Implements(intfBinaryMarshaler)
		d = d || reflect.PointerTo(t).Implements(intfBinaryUnmarshaler)
	}
	return e && d
}
