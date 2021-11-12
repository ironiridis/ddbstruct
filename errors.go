package ddbstruct

import (
	"fmt"
	"strings"

	ddbtype "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type NoItemError struct {
	Key map[string]ddbtype.AttributeValue
}

func (e *NoItemError) Error() string {
	if len(e.Key) == 0 {
		return "no item found (empty key)"
	}

	var buf strings.Builder
	for k := range e.Key {
		if buf.Len() > 0 {
			buf.WriteByte(',')
		}
		fmt.Fprintf(&buf, "%q=%s", k, e.Key[k])
	}

	return "no item found matching key " + buf.String()
}

type EncodeError struct {
	Err     error
	KeyName string
}

func (e *EncodeError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("failed to encode key %q: %s", e.KeyName, e.Err.Error())
	}
	return fmt.Sprintf("failed to encode key %q", e.KeyName)
}

type DecodeError struct {
	Err     error
	KeyName string
}

func (e *DecodeError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("failed to decode key %q: %s", e.KeyName, e.Err.Error())
	}
	return fmt.Sprintf("failed to decode key %q", e.KeyName)
}
