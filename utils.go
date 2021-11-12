package ddbstruct

import (
	"fmt"
	"strings"

	ddbtype "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func attrValString(av ddbtype.AttributeValue) string {
	switch v := av.(type) {
	case *ddbtype.AttributeValueMemberBOOL:
		if v.Value {
			return "true"
		}
		return "false"
	case *ddbtype.AttributeValueMemberNULL:
		return "(null)"
	case *ddbtype.AttributeValueMemberB:
	case *ddbtype.AttributeValueMemberS:
		return fmt.Sprintf("%q", v.Value)
	case *ddbtype.AttributeValueMemberBS:
	case *ddbtype.AttributeValueMemberSS:
		var buf strings.Builder
		buf.WriteByte('[')
		for idx := range v.Value {
			if idx > 0 {
				buf.WriteByte(',')
			}
			fmt.Fprintf(&buf, "%q", v.Value[idx])
		}
		buf.WriteByte(']')
		return buf.String()
	case *ddbtype.AttributeValueMemberN:
		return v.Value
	case *ddbtype.AttributeValueMemberNS:
		var buf strings.Builder
		buf.WriteByte('[')
		for idx := range v.Value {
			if idx > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(v.Value[idx])
		}
		buf.WriteByte(']')
		return buf.String()
	case *ddbtype.AttributeValueMemberL:
		var buf strings.Builder
		buf.WriteByte('[')
		for idx := range v.Value {
			if idx > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(attrValString(v.Value[idx]))
		}
		buf.WriteByte(']')
		return buf.String()
	case *ddbtype.AttributeValueMemberM:
		var buf strings.Builder
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
		return buf.String()
	}

	return fmt.Sprintf("?unk %+v", av)
}
