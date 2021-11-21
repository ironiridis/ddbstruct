package ddbstruct

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func encString(s interface{}, f int) types.AttributeValue {
	return &types.AttributeValueMemberS{Value: getF(s, f).String()}
}
func decString(s interface{}, f int, av types.AttributeValue) {
	getF(s, f).SetString(av.(*types.AttributeValueMemberS).Value)
}

func encBytes(s interface{}, f int) types.AttributeValue {
	return &types.AttributeValueMemberB{Value: getF(s, f).Bytes()}
}
func decBytes(s interface{}, f int, av types.AttributeValue) {
	getF(s, f).SetBytes(av.(*types.AttributeValueMemberB).Value)
}

func encBool(s interface{}, f int) types.AttributeValue {
	return &types.AttributeValueMemberBOOL{Value: getF(s, f).Bool()}
}
func decBool(s interface{}, f int, av types.AttributeValue) {
	getF(s, f).SetBool(av.(*types.AttributeValueMemberBOOL).Value)
}
