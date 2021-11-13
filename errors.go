package ddbstruct

import (
	"fmt"
	"strings"

	ddbtype "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func attrValString(av ddbtype.AttributeValue) string {
	// simple types that don't need built-up return strings
	switch v := av.(type) {
	case *ddbtype.AttributeValueMemberBOOL:
		if v.Value {
			return "true"
		}
		return "false"
	case *ddbtype.AttributeValueMemberNULL:
		return "(null)"
	case *ddbtype.AttributeValueMemberB:
		return fmt.Sprintf("%q", v.Value)
	case *ddbtype.AttributeValueMemberS:
		return fmt.Sprintf("%q", v.Value)
	case *ddbtype.AttributeValueMemberN:
		return v.Value
	}

	// more complicated types that benefit from a string builder
	var buf strings.Builder
	switch v := av.(type) {
	case *ddbtype.AttributeValueMemberBS:
		buf.WriteByte('[')
		for idx := range v.Value {
			if idx > 0 {
				buf.WriteByte(',')
			}
			fmt.Fprintf(&buf, "%q", v.Value[idx])
		}
		buf.WriteByte(']')
	case *ddbtype.AttributeValueMemberSS:
		buf.WriteByte('[')
		for idx := range v.Value {
			if idx > 0 {
				buf.WriteByte(',')
			}
			fmt.Fprintf(&buf, "%q", v.Value[idx])
		}
		buf.WriteByte(']')
	case *ddbtype.AttributeValueMemberNS:
		buf.WriteByte('[')
		for idx := range v.Value {
			if idx > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(v.Value[idx])
		}
		buf.WriteByte(']')
	case *ddbtype.AttributeValueMemberL:
		buf.WriteByte('[')
		for idx := range v.Value {
			if idx > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(attrValString(v.Value[idx]))
		}
		buf.WriteByte(']')
	case *ddbtype.AttributeValueMemberM:
		buf.WriteByte('{')
		var multi bool
		for k := range v.Value {
			if multi {
				buf.WriteByte(',')
			}
			buf.WriteString(fmt.Sprintf("%q:", k))
			buf.WriteString(attrValString(v.Value[k]))
			multi = true
		}
		buf.WriteByte('}')
	default:
		return fmt.Sprintf("unknown attribute %#v", av)
	}
	return buf.String()
}

type NoItemError struct {
	Key avmap
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
		fmt.Fprintf(&buf, "%q=%s", k, attrValString(e.Key[k]))
	}

	return "no item found matching key " + buf.String()
}
