package ddbstruct

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func Put(ctx context.Context, svc *dynamodb.Client, table string, data interface{}) error {
	return fmt.Errorf("unimplemented")
	/*
		var err error
		dmd := cache.get(data)

		putcmd := &dynamodb.PutItemInput{
			Item:      avmap,
			TableName: &table,
		}
		_, err = svc.PutItem(ctx, putcmd)
		return err
	*/
}
