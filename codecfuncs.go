package ddbstruct

import (
	"reflect"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

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
