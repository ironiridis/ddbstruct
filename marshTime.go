package ddbstruct

import (
	"reflect"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func encDurationString(s interface{}, f int) types.AttributeValue {
	return &types.AttributeValueMemberS{Value: getF(s, f).Interface().(time.Duration).String()}
}
func decDurationString(s interface{}, f int, av types.AttributeValue) {
	t, err := time.ParseDuration(av.(*types.AttributeValueMemberS).Value)
	if err != nil {
		panic(err)
	}
	getF(s, f).Set(reflect.ValueOf(t))
}
func encDurationNano(s interface{}, f int) types.AttributeValue {
	ns := getF(s, f).Interface().(time.Duration).Nanoseconds()
	return &types.AttributeValueMemberN{Value: strconv.FormatInt(ns, 10)}
}
func decDurationNano(s interface{}, f int, av types.AttributeValue) {
	nv, err := strconv.ParseInt(av.(*types.AttributeValueMemberN).Value, 10, 64)
	if err != nil {
		panic(err)
	}
	getF(s, f).SetInt(int64(time.Nanosecond) * nv)
}
func encDurationSec(s interface{}, f int) types.AttributeValue {
	sv := getF(s, f).Interface().(time.Duration).Seconds()
	return &types.AttributeValueMemberN{Value: strconv.FormatFloat(sv, 'f', -1, 64)}
}
func decDurationSec(s interface{}, f int, av types.AttributeValue) {
	tv, err := time.ParseDuration(av.(*types.AttributeValueMemberN).Value + "s")
	if err != nil {
		panic(err)
	}
	getF(s, f).Set(reflect.ValueOf(tv))
}

func encTimeNano(s interface{}, f int) types.AttributeValue {
	t := getF(s, f).Interface().(time.Time)
	return &types.AttributeValueMemberN{Value: strconv.FormatInt(t.UnixNano(), 10)}
}
func decTimeNano(s interface{}, f int, av types.AttributeValue) {
	nv, err := strconv.ParseInt(av.(*types.AttributeValueMemberN).Value, 10, 64)
	if err != nil {
		panic(err)
	}
	getF(s, f).Set(reflect.ValueOf(time.Unix(0, nv)))
}
func encTimeEpoch(s interface{}, f int) types.AttributeValue {
	t := getF(s, f).Interface().(time.Time).Unix()
	return &types.AttributeValueMemberN{Value: strconv.FormatInt(t, 10)}
}
func decTimeEpoch(s interface{}, f int, av types.AttributeValue) {
	nv, err := strconv.ParseInt(av.(*types.AttributeValueMemberN).Value, 10, 64)
	if err != nil {
		panic(err)
	}
	getF(s, f).Set(reflect.ValueOf(time.Unix(nv, 0)))
}
