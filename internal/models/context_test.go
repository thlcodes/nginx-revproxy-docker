package models_test

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"

	"github.com/stretchr/testify/assert"

	"dev.beta.audi/gorepo/gopher_skeleton/internal/models"
	"github.com/stretchr/testify/require"
)

// Test_Context_RequetID
func Test_Context_RequetID(t *testing.T) {
	reqID := "abc123"
	type args struct {
		ctx context.Context
	}
	type wants struct {
		reqid string
		ok    bool
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{"nil ctx", args{}, wants{}},
		{"empty ctx", args{context.TODO()}, wants{}},
		{"non-incoming ctx", args{metadata.NewOutgoingContext(context.TODO(), metadata.Pairs())}, wants{}},
		{"incoming ctx w/o reqid", args{metadata.NewIncomingContext(context.TODO(), metadata.Pairs("whatever", reqID))}, wants{}},
		{"incoming ctx w/ empy reqid", args{metadata.NewIncomingContext(context.TODO(), metadata.Pairs(models.RequestIDKey, ""))}, wants{}},
		{"incoming ctx w/ valid reqid", args{metadata.NewIncomingContext(context.TODO(), metadata.Pairs(models.RequestIDKey, reqID))}, wants{reqID, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := models.NewContext(tt.args.ctx)
			require.NotNil(t, ctx)
			gotID, gotOk := ctx.RequestID()
			assert.Equal(t, tt.wants.reqid, gotID)
			assert.Equal(t, tt.wants.ok, gotOk)
		})
	}
}
