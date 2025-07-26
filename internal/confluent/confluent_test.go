package confluent

import (
	"testing"

	"github/erastusk/tracer/internal/types"
	"github/erastusk/tracer/mocks"

	"go.uber.org/mock/gomock"
)

func TestConfluentProduce(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tc := mocks.NewMockConfluentOauth(ctrl)
	tc.EXPECT().Produce(gomock.Any(), gomock.Any()).Return(nil)
	var s types.PromptOptions
	var e error
	expected := e
	actual := tc.Produce(s, "test")
	if expected != actual {
		t.Errorf("expected : %s, got : %s", expected, actual)
	}
}

func TestConfluentConsume(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tc := mocks.NewMockConfluentOauth(ctrl)
	tc.EXPECT().Consume(gomock.Any(), gomock.Any()).Return(nil)
	var s types.PromptOptions
	var e error
	expected := e
	actual := tc.Consume(s, "test")
	if expected != actual {
		t.Errorf("expected : %s, got : %s", expected, actual)
	}
}

func TestConfluentListTopic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tc := mocks.NewMockConfluentOauth(ctrl)
	tc.EXPECT().ListTopics(gomock.Any()).Return(nil)
	var e error
	expected := e
	var s types.PromptOptions
	actual := tc.ListTopics(s)
	if expected != actual {
		t.Errorf("expected : %s, got : %s", expected, actual)
	}
}
