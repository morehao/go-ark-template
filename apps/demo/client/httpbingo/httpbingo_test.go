package httpbingo

import (
	"testing"

	"github.com/morehao/goark/pkg/testsetup"
	"github.com/morehao/golib/glog"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	testsetup.Initialize(testsetup.AppNameDemo)

	ctx := testsetup.NewContext()
	res, err := Get(ctx, &GetRequest{
		ID: 1,
	})
	assert.Nil(t, err)
	t.Log(glog.ToJsonString(res))
}
