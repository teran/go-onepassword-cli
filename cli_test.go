package onepassword

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOnePasswordGetByLabel(t *testing.T) {
	r := require.New(t)
	ctx := context.Background()

	svc := New().(*onePassword)
	svc.opBinary = "testdata/opmock"

	v, err := svc.GetByLabel(ctx, KindPassword, "blah")
	r.NoError(err)
	r.Equal("test_password", v)
}
