package redis_test

import (
	"testing"

	"github/erastusk/tracer/internal/types"
	"github/erastusk/tracer/mocks"

	"go.uber.org/mock/gomock"
)

func TestRedisConnectivity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tc := mocks.NewMockRedisApps(ctrl)
	tc.EXPECT().Connectivity(gomock.Any()).Return(nil)
	var s types.PromptOptions
	var e error
	expected := e
	actual := tc.Connectivity(s)
	if expected != actual {
		t.Errorf("expected : %s, got : %s", expected, actual)
	}
}
