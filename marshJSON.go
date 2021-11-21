package ddbstruct

import (
	"encoding/json"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var intfJSONMarshaler = reflect.TypeOf(new(json.Marshaler)).Elem()
var intfJSONUnmarshaler = reflect.TypeOf(new(json.Unmarshaler)).Elem()

func encJSON(s interface{}, f int) types.AttributeValue {
	fp := getF(s, f)
	enc, ok := fp.Interface().(json.Marshaler)
	if !ok {
		enc, ok = fp.Addr().Interface().(json.Marshaler)
		if !ok {
			panic("neither " + fp.Type().String() + " nor *" + fp.Type().String() + " implements MarshalJSON")
		}
	}

	buf, err := enc.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return &types.AttributeValueMemberS{Value: string(buf)}
}

func decJSON(s interface{}, f int, av types.AttributeValue) {
	fp := getF(s, f)
	dec, ok := fp.Interface().(json.Unmarshaler)
	if !ok {
		dec, ok = fp.Addr().Interface().(json.Unmarshaler)
		if !ok {
			panic("neither " + fp.Type().String() + " nor *" + fp.Type().String() + " implements UnmarshalJSON")
		}
	}

	err := dec.UnmarshalJSON([]byte(av.(*types.AttributeValueMemberS).Value))
	if err != nil {
		panic(err)
	}
}

func isJSONEncoder(t reflect.Type) bool {
	var e, d bool
	e = t.Implements(intfJSONMarshaler)
	d = t.Implements(intfJSONUnmarshaler)

	// if the type isn't already a pointer, check to see if a pointer to
	// the type satisfies the interface instead
	if t.Kind() != reflect.Pointer {
		e = e || reflect.PointerTo(t).Implements(intfJSONMarshaler)
		d = d || reflect.PointerTo(t).Implements(intfJSONUnmarshaler)
	}
	return e && d
}

func encJSONRaw(s interface{}, f int) types.AttributeValue {
	buf, err := json.Marshal(getF(s, f).Interface())
	if err != nil {
		panic(err)
	}
	return &types.AttributeValueMemberS{Value: string(buf)}
}

func decJSONRaw(s interface{}, f int, av types.AttributeValue) {
	fp := getF(s, f)
	if fp.Type().Kind() != reflect.Pointer {
		fp = fp.Addr()
	}
	err := json.Unmarshal([]byte(av.(*types.AttributeValueMemberS).Value), fp.Interface())
	if err != nil {
		panic(err)
	}
}
