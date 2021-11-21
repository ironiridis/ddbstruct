package ddbstruct

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func encInt(s interface{}, f int) types.AttributeValue {
	return &types.AttributeValueMemberN{Value: strconv.FormatInt(getF(s, f).Int(), 10)}
}
func decInt(s interface{}, f int, av types.AttributeValue) {
	d := getF(s, f)
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

func encUint(s interface{}, f int) types.AttributeValue {
	return &types.AttributeValueMemberN{Value: strconv.FormatUint(getF(s, f).Uint(), 10)}
}
func decUint(s interface{}, f int, av types.AttributeValue) {
	d := getF(s, f)
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

func encFloat(bits int) encodeFunc {
	return func(s interface{}, f int) types.AttributeValue {
		d := getF(s, f)
		return &types.AttributeValueMemberN{Value: strconv.FormatFloat(d.Float(), 'G', -1, bits)}
	}
}
func decFloat(bits int) decodeFunc {
	return func(s interface{}, f int, av types.AttributeValue) {
		d := getF(s, f)
		sv := av.(*types.AttributeValueMemberN).Value
		fv, err := strconv.ParseFloat(sv, bits)
		if err != nil {
			panic(fmt.Errorf("cannot convert %q to %s: %w", sv, d.Kind(), err))
		}
		// would check d.OverflowFloat(fv) here, but actually ParseFloat checks for us
		d.SetFloat(fv)
	}
}
