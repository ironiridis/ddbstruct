package ddbstruct

import (
	"context"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func Put(ctx context.Context, svc *dynamodb.Client, table string, data interface{}) (err error) {
	defer func() {
		if panicVal := recover(); panicVal != nil {
			if panicErr, ok := panicVal.(error); ok {
				err = fmt.Errorf("failed to encode for put: %w", panicErr)
			} else {
				err = fmt.Errorf("failed to encode for put: %v", panicVal)
			}
		}
	}()
	src := reflect.ValueOf(data).Elem()
	dmd := cache.get(data)
	putcmd := &dynamodb.PutItemInput{
		Item:      map[string]types.AttributeValue{},
		TableName: &table,
	}
	for idx, f := range dmd.f {
		if f.enc == nil {
			err = fmt.Errorf("no encode function available for field %q of %T", f.fieldname, data)
			return
		}
		val := src.Field(idx)
		if val.IsZero() {
			if f.optional { // skip zero attribute
				continue
			}
			if f.defvalue != "" { // apply default value
				val = reflect.ValueOf(f.defvalue)
			}
		}
		putcmd.Item[f.fieldname] = f.enc(val)
	}
	_, err = svc.PutItem(ctx, putcmd)
	return err
}
