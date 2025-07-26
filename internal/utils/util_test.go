package utils_test

import (
	"testing"

	"github/erastusk/tracer/mocks"

	"go.uber.org/mock/gomock"
)

func TestConnectivity(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	tc := mocks.NewMockTCPDial(ctrl)
	tc.EXPECT().TCPDial(gomock.Any()).Return(nil)
	// act
	server := "endpoint"
	// assert
	var e error
	expected := e
	actual := tc.TCPDial(server)
	if expected != actual {
		t.Errorf("expected : %s, got : %s", expected, actual)
	}
}
