package ddbstruct

import (
	"encoding"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var intfTextMarshaler = reflect.TypeOf(new(encoding.TextMarshaler)).Elem()
var intfTextUnmarshaler = reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()

func encText(s interface{}, f int) types.AttributeValue {
	fp := getF(s, f)
	enc, ok := fp.Interface().(encoding.TextMarshaler)
	if !ok {
		enc, ok = fp.Addr().Interface().(encoding.TextMarshaler)
		if !ok {
			panic("neither " + fp.Type().String() + " nor *" + fp.Type().String() + " implements MarshalText")
		}
	}

	buf, err := enc.MarshalText()
	if err != nil {
		panic(err)
	}
	return &types.AttributeValueMemberS{Value: string(buf)}
}

func decText(s interface{}, f int, av types.AttributeValue) {
	fp := getF(s, f)
	dec, ok := fp.Interface().(encoding.TextUnmarshaler)
	if !ok {
		dec, ok = fp.Addr().Interface().(encoding.TextUnmarshaler)
		if !ok {
			panic("neither " + fp.Type().String() + " nor *" + fp.Type().String() + " implements UnmarshalText")
		}
	}

	err := dec.UnmarshalText([]byte(av.(*types.AttributeValueMemberS).Value))
	if err != nil {
		panic(err)
	}
}

func isTextEncoder(t reflect.Type) bool {
	var e, d bool
	e = t.Implements(intfTextMarshaler)
	d = t.Implements(intfTextUnmarshaler)

	// if the type isn't already a pointer, check to see if a pointer to
	// the type satisfies the interface instead
	if t.Kind() != reflect.Pointer {
		e = e || reflect.PointerTo(t).Implements(intfTextMarshaler)
		d = d || reflect.PointerTo(t).Implements(intfTextUnmarshaler)
	}
	return e && d
}
