package ddbstruct

import (
	"context"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type avmap map[string]types.AttributeValue

func Get(ctx context.Context, svc *dynamodb.Client, table string, data interface{}) (err error) {
	defer func() {
		if panicVal := recover(); panicVal != nil {
			if panicErr, ok := panicVal.(error); ok {
				err = fmt.Errorf("failed to decode for get: %w", panicErr)
			} else {
				err = fmt.Errorf("failed to decode for get: %v", panicVal)
			}
		}
	}()

	dmd := cache.get(data)
	getcmd := &dynamodb.GetItemInput{
		Key:       avmap{},
		TableName: &table,
	}
	if dmd.pk == nil {
		err = fmt.Errorf("no field is tagged as the partitioning key (pk) for %T", data)
		return
	}
	err = dmd.pk.appendAV(getcmd.Key, data)
	if err != nil {
		return
	}
	if dmd.sk != nil {
		err = dmd.sk.appendAV(getcmd.Key, data)
		if err != nil {
			return
		}
	}
	getres, err := svc.GetItem(ctx, getcmd)
	if err != nil {
		return
	}
	if getres.Item == nil {
		err = &NoItemError{Key: getcmd.Key}
		return
	}
	dst := reflect.ValueOf(data).Elem()
	for _, f := range dmd.f {
		// we don't need to re-decode pk or sk into the struct; it's already there
		if f.pk || f.sk {
			continue
		}
		if av, ok := getres.Item[f.name]; !ok {
			if f.optional {
				if !dst.Field(f.idx).IsZero() {
					// this case deals with the fact that, when reading the response from dynamodb, some attributes
					// may be missing. if they're missing, and they're optional, that's fine. but if they are *also*
					// not already at the zero value in the source struct, there isn't a way to distinguish this
					// situation from the stale value being the value retrieved. this could be relaxed later if it
					// seems like a feature that would be useful, but for now this feels like a footgun.
					err = fmt.Errorf("decoding into %T: field %q is optional, value is not zero, and no attribute returned", data, f.name)
					return
				}
				// optional values may be missing returned attributes, so keep going without decoding anything
				continue
			}
			err = fmt.Errorf("missing attribute in response for field %q not tagged optional", f.name)
			return
		} else {
			if f.dec == nil {
				err = fmt.Errorf("field %q is missing decoder func", f.name)
				return
			}
			if !dst.Field(f.idx).CanSet() {
				err = fmt.Errorf("cannot set field %q", f.name)
				return
			}
			f.dec(dst.Field(f.idx), av)
		}
	}
	return
}

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
	dmd := cache.get(data)
	putcmd := &dynamodb.PutItemInput{
		Item:      avmap{},
		TableName: &table,
	}
	for idx := range dmd.f {
		err = dmd.f[idx].appendAV(putcmd.Item, data)
		if err != nil {
			return
		}
	}
	_, err = svc.PutItem(ctx, putcmd)
	return err
}

func Delete(ctx context.Context, svc *dynamodb.Client, table string, data interface{}) (err error) {
	dmd := cache.get(data)
	getcmd := &dynamodb.DeleteItemInput{
		Key:       avmap{},
		TableName: &table,
	}
	if dmd.pk == nil {
		err = fmt.Errorf("no field is tagged as the partitioning key (pk) for %T", data)
		return
	}
	err = dmd.pk.appendAV(getcmd.Key, data)
	if err != nil {
		return
	}
	if dmd.sk != nil {
		err = dmd.sk.appendAV(getcmd.Key, data)
		if err != nil {
			return
		}
	}
	_, err = svc.DeleteItem(ctx, getcmd)
	return
}
