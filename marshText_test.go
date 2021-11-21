package ddbstruct

import (
	"net"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestTimeTextRoundTrip(t *testing.T) {
	type z struct{ X time.Time }
	in := &z{X: time.Now()}
	av := encText(in, 0)
	expectT(t, new(types.AttributeValueMemberS), av)
	out := &z{}
	decText(out, 0, av)
	if in.X.UnixNano() != out.X.UnixNano() {
		t.Fatalf("expected %v, got %v", in.X, out.X)
	}
}

func TestTimePtrTextRoundTrip(t *testing.T) {
	type z struct{ X *time.Time }
	orig := time.Now()
	in := &z{X: &orig}
	av := encText(in, 0)
	expectT(t, new(types.AttributeValueMemberS), av)
	out := &z{}
	decText(out, 0, av)
	if out.X == nil {
		t.Fatal("out is nil")
	}
	if in.X == out.X {
		t.Fatal("out is not a new pointer")
	}
	if in.X.UnixNano() != out.X.UnixNano() {
		t.Fatalf("expected %v, got %v", *in.X, *out.X)
	}
}

func TestIPTextRoundTrip(t *testing.T) {
	type z struct{ X net.IP }
	in := &z{X: net.IP{192, 168, 0, 1}}
	av := encText(in, 0)
	expectT(t, new(types.AttributeValueMemberS), av)
	out := &z{}
	decText(out, 0, av)
	if in.X.String() != out.X.String() {
		t.Fatalf("expected %s, got %s", in.X, out.X)
	}
}

func TestIPPtrTextRoundTrip(t *testing.T) {
	type z struct{ X *net.IP }
	orig := net.IP{192, 168, 0, 2}
	in := &z{X: &orig}
	av := encText(in, 0)
	expectT(t, new(types.AttributeValueMemberS), av)
	out := &z{}
	decText(out, 0, av)
	if out.X == nil {
		t.Fatal("out is nil")
	}
	if in.X == out.X {
		t.Fatal("out is not a new pointer")
	}
	if in.X.String() != out.X.String() {
		t.Fatalf("expected %s, got %s", *in.X, *out.X)
	}
}
